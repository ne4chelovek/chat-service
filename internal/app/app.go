package app

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ne4chelovek/chat_common/pkg/closer"
	"github.com/ne4chelovek/chat_service/internal/interceptor"
	desc "github.com/ne4chelovek/chat_service/pkg/chat_v1"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	_ "github.com/ne4chelovek/chat_service/statik"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort       = 9070
	httpAddress    = "localhost:8080"
	swaggerAddress = "localhost:8090"
	grpcAddress    = "localhost:9090"
	servicePort    = "87.228.39.226:9000"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
	swaggerServer   *http.Server
	httpClient      *http.Client
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
		if err := a.serviceProvider.KafkaConsumer().Stop(); err != nil {
			log.Printf("failed to stop kafka consumer: %v", err)
		}
	}()
	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			log.Fatalf("failed to run GRPC server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			log.Fatalf("failed to run Swagger server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		a.serviceProvider.KafkaConsumer().Start()
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initServiceProvider(ctx context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	//	creds, err := credentials.NewServerTLSFromFile("certs/service.pem", "certs/service.key")
	//	if err != nil {
	//		log.Fatalf("failed to create credentials: %v", err)
	//		return err
	//	}

	authInterceptor := &interceptor.Client{Client: a.serviceProvider.AuthClient()}

	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			interceptor.ValidateInterceptor,
			authInterceptor.AuthInterceptor,
		),
	)

	reflection.Register(a.grpcServer)
	desc.RegisterChatServer(a.grpcServer, a.serviceProvider.ChatImpl(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	//	creds, err := credentials.NewClientTLSFromFile("certs/service.pem", "")
	//	if err != nil {
	//		fmt.Printf("failed to load TLS keys: %v", err)
	//		return err
	//	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := desc.RegisterChatHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:    httpAddress,
		Handler: corsMiddleware.Handler(mux),
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:    swaggerAddress,
		Handler: mux,
	}

	return nil
}

func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running ob %d", grpcPort)

	list, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runHTTPServer() error {
	log.Printf("HTTP server is running on %s", httpAddress)

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	log.Printf("Swagger server is running on %s", swaggerAddress)

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Served swagger file: %s", path)
	}
}

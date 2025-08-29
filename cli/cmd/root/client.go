package root

import (
	descAuth "github.com/ne4chelovek/auth-service/pkg/auth_v1"
	descRegister "github.com/ne4chelovek/auth-service/pkg/users_v1"
	"github.com/ne4chelovek/chat_common/pkg/closer"
	descChat "github.com/ne4chelovek/chat_service/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func authRegisterClient(address string) (descRegister.UsersV1Client, error) {
	//creds, err := credentials.NewClientTLSFromFile(certPath, "")
	//if err != nil {
	//	log.Fatalf("failed to get credentials of registration service: %v", err)
	//	return nil, err
	//}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	closer.Add(conn.Close)

	return descRegister.NewUsersV1Client(conn), nil
}

func authClient(address string) (descAuth.AuthV1Client, error) {
	//	creds, err := credentials.NewClientTLSFromFile(certPath, "")
	//	if err != nil {
	//		log.Fatalf("failed to get credentials of authentication service: %v", err)
	//		return nil, err
	//	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	closer.Add(conn.Close)

	return descAuth.NewAuthV1Client(conn), nil
}

func chatClient(address string) (descChat.ChatClient, error) {
	//creds, err := credentials.NewClientTLSFromFile(certPath, "")
	//if err != nil {
	//	log.Fatalf("failed to get credentials of chat service: %v", err)
	//	return nil, err
	//}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	closer.Add(conn.Close)

	return descChat.NewChatClient(conn), nil
}

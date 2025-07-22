package root

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"log"
)

const filename = ".access_token"

var rootCmd = &cobra.Command{
	Use:   "chat-service",
	Short: "CLI для работы с чат-сервисом",
}

var (
	addressAuth string
	addressChat string
	certPath    string
	chatID      int
	users       []string
	username    string
	email       string
	password    string
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Регистрация нового пользователя",
	Run: func(cmd *cobra.Command, _ []string) {
		err := register(context.Background(), addressAuth, username, email, password)
		if err != nil {
			log.Fatalf("Ошибка регистрации: %v", err)
		}
		fmt.Print(color.GreenString("Успешная регистрация!\n"))
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Вход в чат",
	Run: func(cmd *cobra.Command, _ []string) {
		err := login(context.Background(), addressAuth)
		if err != nil {
			log.Fatalf("failed to login: %v", err)
		}
		fmt.Print(color.GreenString("Successfully logged in\n"))
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Выход из чата",
	Run: func(_ *cobra.Command, _ []string) {
		err := logout(context.Background(), addressAuth)
		if err != nil {
			log.Fatalf("failed to logout: %v", err)
		}
		fmt.Print(color.GreenString("Successfully logged out\n"))
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Создаёт чат",
	Run: func(_ *cobra.Command, _ []string) {
		if len(users) < 2 {
			log.Fatal("Нужно минимум 2 пользователя")
		}
		id, err := createChat(context.Background(), addressChat, users)
		if err != nil {
			log.Fatalf("failed to create chat: %v", err)
		}
		fmt.Printf(color.GreenString("Successfully created chat with id: %v\n"), id)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Удаляет чат",
	Run: func(_ *cobra.Command, _ []string) {
		err := deleteChat(context.Background(), addressChat, int64(chatID))
		if err != nil {
			log.Fatalf("failed to delete chat: %v", err)
		}
		fmt.Print(color.GreenString("Chat deleted successfully\n"))
	},
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Подключение к чату",
	Run: func(_ *cobra.Command, _ []string) {
		err := connect(context.Background(), addressChat, int64(chatID))
		if err != nil {
			log.Fatalf("failed to connect: %v", err)
		}
		fmt.Print(color.GreenString("Successfully connected\n"))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	//	if err := godotenv.Load("../.env"); err != nil {
	//		log.Fatal("Error loading .env file")
	//	}
	//	certPath = os.Getenv("CERT_PATH_CLI")
	//	if certPath == "" {
	//		log.Fatal("CERT_PATH not set in .env")
	//	}

	rootCmd.PersistentFlags().StringVar(&addressAuth, "auth-address", "87.228.39.226:9000", "Адрес сервера аутентификации")
	rootCmd.PersistentFlags().StringVar(&addressChat, "chat-address", "87.228.39.227:9070", "Адрес чат-сервера")
	//rootCmd.PersistentFlags().StringVar(&certPath, "cert", certPath, "Путь к TLS сертификату")

	registerCmd.Flags().StringVarP(&username, "username", "u", "", "Имя пользователя")
	registerCmd.Flags().StringVarP(&email, "email", "e", "", "Email пользователя")
	registerCmd.Flags().StringVarP(&password, "password", "p", "", "Пароль пользователя")
	_ = registerCmd.MarkFlagRequired("username")
	_ = registerCmd.MarkFlagRequired("email")
	_ = registerCmd.MarkFlagRequired("password")

	createCmd.Flags().StringSliceVarP(&users, "users", "u", []string{}, "Список пользователей")
	_ = createCmd.MarkFlagRequired("users")

	deleteCmd.Flags().IntVarP(&chatID, "chat-id", "c", 0, "ID чата")
	_ = deleteCmd.MarkFlagRequired("chat-id")

	connectCmd.Flags().IntVarP(&chatID, "chat-id", "c", 0, "ID чата")
	_ = connectCmd.MarkFlagRequired("chat-id")

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(connectCmd)
}

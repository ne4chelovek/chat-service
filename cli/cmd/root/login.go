package root

import (
	"bufio"
	"context"
	"fmt"
	"github.com/ne4chelovek/auth-service/pkg/auth_v1"
	"github.com/ne4chelovek/chat_common/pkg/closer"
	"golang.org/x/term"
	"os"
)

func login(ctx context.Context, address string) error {
	client, err := authClient(address)
	if err != nil {
		return err
	}

	var username string
	fmt.Println("Login: ")
	_, err = fmt.Scanln(&username)
	if err != nil {
		return err
	}

	fmt.Println("Password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("ошибка ввода пароля: %w", err)
	}
	fmt.Println()

	password := string(passwordBytes)

	response, err := client.Login(ctx, &auth_v1.LoginRequest{Password: password, Usernames: username})
	if err != nil {
		return err
	}

	resAccessToken, err := client.GetAccessToken(ctx, &auth_v1.GetAccessTokenRequest{
		RefreshToken: response.GetRefreshToken(),
	})
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	closer.Add(file.Close)

	wr := bufio.NewWriter(file)
	_, err = wr.WriteString(resAccessToken.GetAccessToken())
	if err != nil {
		return err
	}
	err = wr.Flush()
	if err != nil {
		return err
	}
	return nil
}

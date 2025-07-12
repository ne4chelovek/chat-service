package root

import (
	"context"
	"github.com/ne4chelovek/auth-service/pkg/auth_v1"
	"os"
)

func logout(ctx context.Context, address string) error {
	client, err := authClient(address)
	if err != nil {
		return err
	}

	accessToken, err := readToken()

	_, err = client.Logout(ctx, &auth_v1.LogoutRequest{AccessToken: accessToken})
	if err != nil {
		return err
	}
	return os.Remove(filename)
}

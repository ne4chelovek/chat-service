package root

import (
	"context"
	"github.com/ne4chelovek/auth-service/pkg/users_v1"
)

func register(ctx context.Context, address, userName, email, password string) error {
	client, err := authRegisterClient(address)
	if err != nil {
		return err
	}

	newUser := &users_v1.CreateUser{
		Name:            userName,
		Email:           email,
		Password:        password,
		PasswordConfirm: password,
		Role:            users_v1.Role_user,
	}

	_, err = client.Create(ctx, &users_v1.CreateRequest{User: newUser})
	if err != nil {
		return err
	}

	return nil
}

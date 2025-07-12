package auth

import (
	"context"
	"github.com/ne4chelovek/auth-service/pkg/access_v1"
	"github.com/ne4chelovek/chat_service/internal/client/rpc"
)

type authClient struct {
	client access_v1.AccessV1Client
}

func NewAuthClient(client access_v1.AccessV1Client) rpc.AuthClient {
	return &authClient{client: client}
}

func (a *authClient) Check(ctx context.Context, endpoint string) error {
	_, err := a.client.Check(ctx, &access_v1.CheckRequest{EndpointAddress: endpoint})
	return err
}

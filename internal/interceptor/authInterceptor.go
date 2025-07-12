package interceptor

import (
	"context"
	"github.com/ne4chelovek/chat_service/internal/client/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Client struct {
	Client rpc.AuthClient
}

func (c *Client) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	err := c.Client.Check(metadata.NewOutgoingContext(ctx, md), info.FullMethod)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "access denied: %v", err)
	}

	return handler(ctx, req)
}

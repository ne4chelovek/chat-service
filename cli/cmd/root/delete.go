package root

import (
	"context"
	"github.com/ne4chelovek/chat_service/pkg/chat_v1"
	"google.golang.org/grpc/metadata"
)

func deleteChat(ctx context.Context, address string, id int64) error {
	accessToken, err := readToken()

	client, err := chatClient(address)
	if err != nil {
		return err
	}

	claims, err := getTokenClaims(accessToken)
	if err != nil {
		return err
	}

	err = isTokenExpired(claims)
	if err != nil {
		return err
	}

	md := metadata.New(map[string]string{"Authorization": "Bearer " + accessToken})
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err = client.Delete(ctx, &chat_v1.DeleteRequest{Id: id})
	if err != nil {
		return err
	}

	return nil
}

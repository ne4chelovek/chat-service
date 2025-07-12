package root

import (
	"context"
	"fmt"
	"github.com/ne4chelovek/chat_service/pkg/chat_v1"
	"google.golang.org/grpc/metadata"
)

func createChat(ctx context.Context, address string, usernames []string) (int64, error) {
	accessToken, err := readToken()

	client, err := chatClient(address)
	if err != nil {
		return 0, err
	}

	claims, err := getTokenClaims(accessToken)
	if err != nil {
		fmt.Println("Ошибка в анмаршале")
		return 0, err
	}

	err = isTokenExpired(claims)
	if err != nil {
		return 0, err
	}

	md := metadata.New(map[string]string{"Authorization": "Bearer " + accessToken})
	ctx = metadata.NewOutgoingContext(ctx, md)

	res, err := client.Create(ctx, &chat_v1.CreateRequest{Usernames: usernames})
	if err != nil {
		return 0, err
	}

	return res.GetId(), nil
}

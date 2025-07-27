package root

import (
	"bufio"
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/ne4chelovek/chat_service/pkg/chat_v1"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"os"
	"strconv"
)

func connect(ctx context.Context, address string, chatId int64) error {
	accessToken, err := readToken()
	if err != nil {
		return err
	}

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

	id := strconv.Itoa(int(chatId))

	stream, err := client.ConnectChat(ctx, &chat_v1.ConnectChatRequest{ChatId: id, Username: claims.Username})
	if err != nil {
		return err
	}
	fmt.Println(color.GreenString("Successfully connect"))

	go func() {
		for {
			//Recv() будет "ждать" сообщения, не мешая при этом вводу с клавиатуры
			message, errRec := stream.Recv()
			if errRec == io.EOF {
				return
			}
			if errRec != nil {
				log.Printf("Error receiving message: %s", errRec)
				return
			}
			fmt.Printf("%v - from: %s: %s",
				message.CreatedAt,
				color.GreenString(message.GetFrom()),
				message.Text,
			)
		}
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if err != nil {
			return err
		}
		fmt.Println("\033[1A")
		_, err = client.SendMessage(ctx, &chat_v1.SendMessageRequest{ChatId: chatId, Message: &chat_v1.Message{Text: text, From: claims.Username}})
		if err != nil {
			log.Printf("failed to send message: %v", err)
		}
	}

}

package tests

import (
	"context"
	"fmt"
	"github.com/ne4chelovek/chat_service/internal/api/chat"
	"github.com/ne4chelovek/chat_service/internal/service"
	"github.com/ne4chelovek/chat_service/internal/service/mocks"
	"testing"

	desc "github.com/ne4chelovek/chat_service/pkg/chat_v1"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	type chatSeviceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		ctx context.Context
		req *desc.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id         = gofakeit.Int64()
		users      = []string{gofakeit.Username(), gofakeit.Username(), gofakeit.Username()}
		serviceErr = fmt.Errorf("service error")

		req = &desc.GetRequest{
			Id: id,
		}

		res = &desc.GetResponse{
			Usernames: users,
		}
	)

	test := []struct {
		name           string
		args           args
		want           *desc.GetResponse
		err            error
		chatSeviceMock chatSeviceMockFunc
	}{
		{
			name: "saccess case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			chatSeviceMock: func(mc *minimock.Controller) service.ChatService {
				mock := mocks.NewChatServiceMock(mc)
				mock.GetChatInfoMock.Expect(ctx, id).Return(users, nil)
				return mock
			},
		},

		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			chatSeviceMock: func(mc *minimock.Controller) service.ChatService {
				mock := mocks.NewChatServiceMock(mc)
				mock.GetChatInfoMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range test {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			chatServiceMock := tt.chatSeviceMock(mc)
			api := chat.NewImplementation(chatServiceMock)

			resHand, err := api.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resHand)
		})
	}

}

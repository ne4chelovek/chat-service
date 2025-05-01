package tests

import (
	"context"
	"fmt"
	"github.com/ne4chelovek/chat_service/internal/api/chat"
	"github.com/ne4chelovek/chat_service/internal/service"
	serviceMocks "github.com/ne4chelovek/chat_service/internal/service/mocks"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	desc "github.com/ne4chelovek/chat_service/pkg/chat_v1"

	"github.com/brianvoe/gofakeit"
)

func TestCreat(t *testing.T) {
	type chatSeviceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		users = []string{gofakeit.Username(), gofakeit.Username(), gofakeit.Username()}

		serviceErr = fmt.Errorf("service error")

		req = &desc.CreateRequest{
			Usernames: users,
		}

		res = &desc.CreateResponse{
			Id: id,
		}
	)

	test := []struct {
		name           string
		args           args
		want           *desc.CreateResponse
		err            error
		chatSeviceMock chatSeviceMockFunc
	}{
		{
			name: "chat api saccess case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			chatSeviceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.CreateMock.Expect(ctx, users).Return(id, nil)
				return mock
			},
		},

		{
			name: "chat api error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			chatSeviceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.CreateMock.Expect(ctx, users).Return(0, serviceErr)
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

			newID, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}

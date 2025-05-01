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
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestDelet(t *testing.T) {
	type chatSeviceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		ctx context.Context
		req *desc.DeleteRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id         = gofakeit.Int64()
		serviceErr = fmt.Errorf("service error")

		req = &desc.DeleteRequest{
			Id: id,
		}
	)

	test := []struct {
		name           string
		args           args
		want           *emptypb.Empty
		err            error
		chatSeviceMock chatSeviceMockFunc
	}{
		{
			name: "saccess case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  nil,
			chatSeviceMock: func(mc *minimock.Controller) service.ChatService {
				mock := mocks.NewChatServiceMock(mc)
				mock.DeleteChatMock.Expect(ctx, id).Return(nil, nil)
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
				mock.DeleteChatMock.Expect(ctx, id).Return(nil, serviceErr)
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

			resHand, err := api.Delete(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resHand)
		})
	}
}

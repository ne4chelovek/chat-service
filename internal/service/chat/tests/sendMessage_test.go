package tests

import (
	"context"
	"fmt"
	"github.com/ne4chelovek/chat_common/pkg/db"
	dbMocks "github.com/ne4chelovek/chat_common/pkg/db/mocks"
	"github.com/ne4chelovek/chat_common/pkg/db/transaction"
	"github.com/ne4chelovek/chat_service/internal/model"
	"github.com/ne4chelovek/chat_service/internal/repository"
	repoMocks "github.com/ne4chelovek/chat_service/internal/repository/mocks"
	"github.com/ne4chelovek/chat_service/internal/service/chat"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"github.com/brianvoe/gofakeit"
)

func TestSendMessage(t *testing.T) {
	t.Parallel()
	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository
	type logRepositoryMockFunc func(mc *minimock.Controller) repository.LogRepository
	type transactorMockFunc func(mc *minimock.Controller) db.Transactor

	type args struct {
		ctx context.Context
		req *model.Message
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		users = gofakeit.Username()
		text  = "someone text"

		opts = pgx.TxOptions{IsoLevel: pgx.ReadCommitted}

		repositoryErr = fmt.Errorf("failed to send message")

		req = &model.Message{
			From: users,
			Text: text,
		}
		log    = fmt.Sprintf("SendMessage in chat with id: %d, message: %v", id, req)
		status = "SENT"
	)

	test := []struct {
		name              string
		args              args
		want              string
		err               error
		chatSeviceMock    chatRepositoryMockFunc
		logRepositoryMock logRepositoryMockFunc
		transactorMock    transactorMockFunc
	}{
		{
			name: "success create case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: status,
			err:  nil,
			chatSeviceMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.SendMessageMock.Expect(minimock.AnyContext, id, req).Return(status, nil)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				mock.LogMock.Expect(minimock.AnyContext, log).Return(nil)
				return mock
			},
			transactorMock: func(mc *minimock.Controller) db.Transactor {
				mock := dbMocks.NewTransactorMock(mc)
				txMock := dbMocks.NewTxMock(mc)
				mock.BeginTxMock.Expect(minimock.AnyContext, opts).Return(txMock, nil)
				txMock.CommitMock.Expect(minimock.AnyContext).Return(nil)
				return mock
			},
		},
		{
			name: "error create case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: "",
			err:  repositoryErr,
			chatSeviceMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.SendMessageMock.Expect(minimock.AnyContext, id, req).Return("", repositoryErr)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				return mock
			},
			transactorMock: func(mc *minimock.Controller) db.Transactor {
				mock := dbMocks.NewTransactorMock(mc)
				txMock := dbMocks.NewTxMock(mc)
				mock.BeginTxMock.Expect(minimock.AnyContext, opts).Return(txMock, nil)
				txMock.RollbackMock.Expect(minimock.AnyContext).Return(nil)
				return mock
			},
		},
		{
			name: "log create error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: "",
			err:  repositoryErr,
			chatSeviceMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.SendMessageMock.Expect(minimock.AnyContext, id, req).Return(status, nil)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				mock.LogMock.Expect(minimock.AnyContext, log).Return(repositoryErr)
				return mock
			},
			transactorMock: func(mc *minimock.Controller) db.Transactor {
				mock := dbMocks.NewTransactorMock(mc)
				txMock := dbMocks.NewTxMock(mc)
				mock.BeginTxMock.Expect(minimock.AnyContext, opts).Return(txMock, nil)
				txMock.RollbackMock.Expect(minimock.AnyContext).Return(nil)
				return mock
			},
		},
	}

	for _, tt := range test {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			chatServiceMock := tt.chatSeviceMock(mc)
			logRepositoryMock := tt.logRepositoryMock(mc)
			transactorMock := transaction.NewTransactionManager(tt.transactorMock(mc))
			service := chat.NewService(chatServiceMock, logRepositoryMock, transactorMock)

			res, err := service.SendMessage(tt.args.ctx, id, tt.args.req)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

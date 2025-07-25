package tests

import (
	"context"
	"fmt"
	"github.com/ne4chelovek/chat_common/pkg/db"
	dbMocks "github.com/ne4chelovek/chat_common/pkg/db/mocks"
	"github.com/ne4chelovek/chat_common/pkg/db/transaction"
	clientApi "github.com/ne4chelovek/chat_service/internal/openApi"
	"github.com/ne4chelovek/chat_service/internal/repository"
	repoMocks "github.com/ne4chelovek/chat_service/internal/repository/mocks"
	"github.com/ne4chelovek/chat_service/internal/service/chat"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"github.com/brianvoe/gofakeit"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository
	type logRepositoryMockFunc func(mc *minimock.Controller) repository.LogRepository
	type apiInternalMock func(mc *minimock.Controller) clientApi.ApiCat
	type transactorMockFunc func(mc *minimock.Controller) db.Transactor

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		users = []string{gofakeit.Username(), gofakeit.Username(), gofakeit.Username()}

		opts = pgx.TxOptions{IsoLevel: pgx.ReadCommitted}

		repositoryErr = fmt.Errorf("failed to get chat")

		req = id

		log = fmt.Sprintf("get chat with id: %v", id)
	)

	test := []struct {
		name string
		args args
		want []string
		err  error

		chatSeviceMock    chatRepositoryMockFunc
		logRepositoryMock logRepositoryMockFunc
		apiInternalMock   apiInternalMock
		transactorMock    transactorMockFunc
	}{
		{
			name: "chat repository success get case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: users,
			err:  nil,
			chatSeviceMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.GetChatInfoMock.Expect(minimock.AnyContext, id).Return(users, nil)
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
			name: "chat repository get error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  repositoryErr,
			chatSeviceMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.GetChatInfoMock.Expect(minimock.AnyContext, id).Return(nil, repositoryErr)
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
			name: "chat repository log get error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  repositoryErr,
			chatSeviceMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.GetChatInfoMock.Expect(minimock.AnyContext, id).Return(users, nil)
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
			clientApiMock := tt.apiInternalMock(mc)
			transactorMock := transaction.NewTransactionManager(tt.transactorMock(mc))
			service := chat.NewService(chatServiceMock, logRepositoryMock, clientApiMock, transactorMock)

			res, err := service.GetChatInfo(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

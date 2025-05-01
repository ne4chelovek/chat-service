package client

import "context"

type AuthClient interface {
	Check(ctx context.Context, endpoint string) (bool, error)
}

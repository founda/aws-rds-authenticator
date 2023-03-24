package authtoken

import (
	"context"
)

type Builder interface {
	BuildToken(ctx context.Context, endpoint, region, dbUser string) (string, error)
}

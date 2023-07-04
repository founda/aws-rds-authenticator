package authtoken

import (
	"context"
)

type TokenBuilder struct {
	password string
}

func NewTokenBuilder(passwd string) *TokenBuilder {
	return &TokenBuilder{
		password: passwd,
	}
}

func (b *TokenBuilder) BuildToken(ctx context.Context, endpoint, region, dbUser string) (string, error) {
	return b.password, nil
}

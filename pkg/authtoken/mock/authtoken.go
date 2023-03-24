package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type TokenBuilder struct {
	mock.Mock
}

func NewMockTokenBuilder() *TokenBuilder {
	m := new(TokenBuilder)
	m.On("BuildToken", context.TODO(), "rds.amazon.com:5432", "eu-west-1", "postgres").Return("t0k3n", nil)

	return m
}

func (m *TokenBuilder) BuildToken(ctx context.Context, endpoint, region, dbUser string) (string, error) {
	args := m.Called(ctx, endpoint, region, dbUser)
	return args.String(0), args.Error(1)
}

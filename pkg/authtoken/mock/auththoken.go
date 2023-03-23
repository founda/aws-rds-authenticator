package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/founda/aws-rds-authenticator/pkg/authtoken"
)

type MockTokenBuilder struct {
	mock.Mock
}

func NewMockTokenBuilder() *MockTokenBuilder {
	m := new(MockTokenBuilder)
	m.On("BuildAuthToken", context.TODO(), "rds.amazon.com:5432", "eu-west-1", "postgres", nil).Return("t0k3n", nil)

	return m
}

func (m *MockTokenBuilder) BuildAuthToken(ctx context.Context, endpoint, region, dbUser string, creds authtoken.CredentialsProvider) (string, error) {
	args := m.Called(ctx, endpoint, region, dbUser, creds)
	return args.String(0), args.Error(1)
}

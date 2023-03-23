package authtoken

import (
	"context"
)

type Credentials interface {
	Expired() bool
	HasKeys() bool
}

type CredentialsProvider interface {
	Retrieve(ctx context.Context) (Credentials, error)
}

type Builder interface {
	BuildAuthToken(ctx context.Context, endpoint, region, dbUser string, creds CredentialsProvider) (string, error)
}

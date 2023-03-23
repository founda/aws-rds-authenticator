package authtoken

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
)

type RDSTokenBuilder struct{}

func (b RDSTokenBuilder) BuildAuthToken(ctx context.Context, endpoint, region, dbUser string, creds aws.CredentialsProvider) (string, error) {
	return auth.BuildAuthToken(ctx, endpoint, region, dbUser, creds)
}

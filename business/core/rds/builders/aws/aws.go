package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
)

type TokenBuilder struct {
	cfg aws.Config
}

func NewTokenBuilder(ctx context.Context) (*TokenBuilder, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return &TokenBuilder{}, err
	}

	return &TokenBuilder{
		cfg: cfg,
	}, nil
}

func (b *TokenBuilder) BuildToken(ctx context.Context, endpoint, region, dbUser string) (string, error) {
	return auth.BuildAuthToken(ctx, endpoint, region, dbUser, b.cfg.Credentials)
}

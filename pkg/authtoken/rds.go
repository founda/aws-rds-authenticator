package authtoken

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
)

type RDSTokenBuilder struct {
	cfg aws.Config
}

func NewRDSTokenBuilder(ctx context.Context) (*RDSTokenBuilder, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return &RDSTokenBuilder{}, err
	}

	return &RDSTokenBuilder{
		cfg: cfg,
	}, nil
}

func (b *RDSTokenBuilder) BuildToken(ctx context.Context, endpoint, region, dbUser string) (string, error) {
	return auth.BuildAuthToken(ctx, endpoint, region, dbUser, b.cfg.Credentials)
}

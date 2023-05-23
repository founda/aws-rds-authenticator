package rds

import (
	"context"
	"fmt"
)

type Builder interface {
	BuildToken(ctx context.Context, endpoint, region, dbUser string) (string, error)
}

type Core struct {
	builder Builder
}

func NewCore(builder Builder) *Core {
	return &Core{
		builder: builder,
	}
}

func (c *Core) Create(ctx context.Context, nc NewConfig) (Config, error) {
	endpoint := fmt.Sprintf("%s:%d", nc.Host, nc.Port)

	token, err := c.builder.BuildToken(ctx, endpoint, nc.Region, nc.Username)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Username: nc.Username,
		Password: token,
		Host:     nc.Host,
		Port:     nc.Port,
		Region:   nc.Region,
	}

	return cfg, nil
}

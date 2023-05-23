package commands

import (
	"errors"

	"github.com/founda/aws-rds-authenticator/business/core/rds"
)

// ErrHelp provides context that help was given.
var ErrHelp = errors.New("provided help")

type AppNewConfig struct {
	Username   string
	Password   string
	Host       string
	Port       int
	Region     string
	DisableTLS bool
	Name       string
	CAPath     string
}

func toCoreNewConfig(config AppNewConfig) rds.NewConfig {
	return rds.NewConfig{
		Username: config.Username,
		Password: config.Password,
		Host:     config.Host,
		Port:     config.Port,
		Region:   config.Region,
	}
}

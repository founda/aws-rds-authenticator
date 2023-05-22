package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ardanlabs/conf/v3"

	"github.com/founda/aws-rds-authenticator/app/aws-rds-authenticator/commands"
	"github.com/founda/aws-rds-authenticator/business/core/rds"
	"github.com/founda/aws-rds-authenticator/business/core/rds/builders/aws"
	"github.com/founda/aws-rds-authenticator/business/core/rds/builders/local"
)

var build = "develop"

type config struct {
	conf.Version
	conf.Args
	DB struct {
		Username   string `conf:"required"`
		Password   string `conf:"mask"`
		Host       string `conf:"required"`
		Port       int    `conf:"default:5432"`
		Region     string `conf:"required"`
		DisableTLS bool   `conf:"default:false"`
		Name       string
		CAPath     string
	}
}

func main() {
	if err := run(); err != nil {
		if !errors.Is(err, commands.ErrHelp) {
			fmt.Println("ERROR", err)
		}
		os.Exit(1)
	}
}

func run() error {
	cfg := config{
		Version: conf.Version{
			Build: build,
			Desc:  "copyright information here",
		},
	}

	const prefix = "RDS_AUTH"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}

		return err
	}

	return processCommands(cfg.Args, cfg)
}

// processCommands handles the execution of the commands specified on
// the command line.
func processCommands(args conf.Args, cfg config) error {
	dbConf := commands.AppNewConfig{
		Username:   cfg.DB.Username,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Port:       cfg.DB.Port,
		Region:     cfg.DB.Region,
		DisableTLS: cfg.DB.DisableTLS,
		Name:       cfg.DB.Name,
		CAPath:     cfg.DB.CAPath,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tb rds.Builder
	var err error
	if dbConf.Password != "" {
		tb = local.NewTokenBuilder(dbConf.Password)
	} else {
		tb, err = aws.NewTokenBuilder(ctx)
		if err != nil {
			return err
		}
	}

	switch args.Num(0) {
	case "create-db":
		out, err := conf.String(&cfg)
		if err != nil {
			return fmt.Errorf("generating config for output: %w", err)
		}
		fmt.Println(out)

		dbName := args.Num(1)
		if dbName == "" {
			return errors.New("missing database name")
		}
		if err := commands.CreateDB(ctx, dbConf, tb, dbName); err != nil {
			return fmt.Errorf("creating database failed: %w", err)
		}
	case "print-dsn":
		if err := commands.PrintDSN(ctx, dbConf, tb, os.Stdout); err != nil {
			return fmt.Errorf("printing DSN failed: %w", err)
		}
	default:
		fmt.Println("create-db NAME:  create the database if it does not already exist")
		fmt.Println("print-dsn:       print the dsn")
		fmt.Println("provide a command to get more help.")
		return commands.ErrHelp
	}

	return nil
}

package commands

import (
	"context"
	"fmt"

	"github.com/founda/aws-rds-authenticator/business/core/rds"
	"github.com/founda/aws-rds-authenticator/business/data/dbcreate"
	"github.com/founda/aws-rds-authenticator/business/sys/database"
)

func CreateDB(ctx context.Context, cfg AppNewConfig, tb rds.Builder, dbName string) error {
	rdsCore := rds.NewCore(tb)
	rdsConfig, err := rdsCore.Create(ctx, toCoreNewConfig(cfg))
	if err != nil {
		return err
	}

	db, err := database.Open(database.Config{
		User:       rdsConfig.Username,
		Password:   rdsConfig.Password,
		Host:       rdsConfig.Host,
		Port:       rdsConfig.Port,
		Region:     rdsConfig.Region,
		DisableTLS: cfg.DisableTLS,
		Name:       cfg.Name,
		CAPath:     cfg.CAPath,
	})
	if err != nil {
		return fmt.Errorf("connection to database failed: %w", err)
	}
	defer db.Close()

	if err = dbcreate.Create(ctx, db, dbName); err != nil {
		return fmt.Errorf("create database %q failed: %w", dbName, err)
	}

	return nil
}

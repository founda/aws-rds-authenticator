package commands

import (
	"context"
	"fmt"
	"io"

	"github.com/founda/aws-rds-authenticator/business/core/rds"
	"github.com/founda/aws-rds-authenticator/business/sys/database"
)

// Migrate creates the schema in the database.
func PrintDSN(ctx context.Context, cfg AppNewConfig, tb rds.Builder, w io.Writer) error {
	rdsCore := rds.NewCore(tb)
	rdsConfig, err := rdsCore.Create(ctx, toCoreNewConfig(cfg))
	if err != nil {
		return err
	}

	dbCfg := database.Config{
		User:       rdsConfig.Username,
		Password:   rdsConfig.Password,
		Host:       rdsConfig.Host,
		Port:       rdsConfig.Port,
		Region:     rdsConfig.Region,
		DisableTLS: cfg.DisableTLS,
		Name:       cfg.Name,
		CAPath:     cfg.CAPath,
	}
	db, err := database.Open(dbCfg)
	if err != nil {
		return fmt.Errorf("connection to database failed: %w", err)
	}
	defer db.Close()

	err = database.StatusCheck(ctx, db)
	if err != nil {
		return fmt.Errorf("database not ready: %w", err)
	}

	connStr, err := database.CreateDSN(dbCfg)
	if err != nil {
		return fmt.Errorf("error creating connection string: %w", err)
	}

	fmt.Fprint(w, connStr)

	return nil
}

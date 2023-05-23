package dbcreate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/founda/aws-rds-authenticator/business/sys/database"
)

func Create(ctx context.Context, db *sql.DB, dbName string) error {
	err := database.StatusCheck(ctx, db)
	if err != nil {
		return fmt.Errorf("database not ready: %w", err)
	}

	q := `SELECT FROM pg_database WHERE datname = $1`

	if err := db.QueryRowContext(ctx, q, dbName).Scan(); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			fmt.Printf("DB %q not found, creating it...\n", dbName)

			// TODO write this is unsafe by design can't do it otherwise
			q = fmt.Sprintf("CREATE DATABASE %s", dbName)

			_, err = db.ExecContext(ctx, q)
			if err != nil {
				return fmt.Errorf("unable to execute query: %w", err)
			}

			fmt.Printf("DB %q successfully created\n", dbName)

			return nil
		default:
			return fmt.Errorf("unable to query databases: %w", err)
		}
	}

	fmt.Printf("DB %q already present\n", dbName)

	return nil
}

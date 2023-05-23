package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	User       string
	Password   string
	Host       string
	Port       int
	Region     string
	DisableTLS bool
	Name       string
	CAPath     string
}

func CreateDSN(cfg Config) (string, error) {
	sslMode := "verify-full"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	q := make(url.Values)
	q.Set("sslmode", sslMode)

	if cfg.CAPath != "" {
		if err := pathExist(cfg.CAPath); err != nil {
			return "", err
		}

		q.Set("sslrootcert", cfg.CAPath)
	}

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	return u.String(), nil
}

// Open knows how to open a database connection based on the configuration.
func Open(cfg Config) (*sql.DB, error) {
	dsn, err := CreateDSN(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating connection string: %w", err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot open database connection: %w", err)
	}

	return db, nil
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *sql.DB) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second)
		defer cancel()
	}

	var pingError error
	for attempts := 1; ; attempts++ {
		// TODO: Ping() or PingContext(ctx) do not seems to respect the context
		// and if there is a network issue just hangs for about 60s.
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		fmt.Println(pingError)
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Run a simple query to determine connectivity.
	// Running this query forces a round trip through the database.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}

func pathExist(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("CA path does not exist: %s", path)
		}
		return fmt.Errorf("error checking for CA path: %w", err)
	}

	return nil
}

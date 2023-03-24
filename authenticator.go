package authenticator

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/founda/aws-rds-authenticator/pkg/authtoken"
)

type option func(*authenticator) error

type authenticator struct {
	output           io.Writer
	engine           string
	host             string
	port             int
	region           string
	user             string
	database         string
	authTokenBuilder authtoken.Builder
}

func NewAuthenticator(opts ...option) (authenticator, error) {
	auth := authenticator{
		output: os.Stdout,
	}

	for _, opt := range opts {
		err := opt(&auth)
		if err != nil {
			return authenticator{}, err
		}
	}

	return auth, nil
}

func WithOutput(w io.Writer) option {
	return func(a *authenticator) error {
		if w == nil {
			return errors.New("nil output writer")
		}
		a.output = w
		return nil
	}
}

func FromArgs(args []string) option {
	return func(a *authenticator) error {
		fset := flag.NewFlagSet("aws-rds-authenticator", flag.ExitOnError)

		hostPtr := fset.String("host", "", "Endpoint of the database instance")
		portPtr := fset.Int("port", 0, "Port number used for connecting to your DB instance (default postgres: 5432, default mysql: 3306)")
		regionPtr := fset.String("region", "", "AWS Region where the database instance is running")
		userPtr := fset.String("user", "", "Database account that you want to access")
		databasePtr := fset.String("database", "", "Database that you want to access")
		enginePtr := fset.String("engine", "postgres", "Database engine that you want to access: postgres|mysql")

		err := fset.Parse(args)
		if err != nil {
			return err
		}

		//TODO: add more validation
		if *hostPtr == "" {
			return errors.New("missing required host")
		}
		if *regionPtr == "" {
			return errors.New("missing required region")
		}
		if *userPtr == "" {
			return errors.New("missing required user")
		}
		if *databasePtr == "" {
			return errors.New("missing required database")
		}
		if *enginePtr == "" {
			return errors.New("missing required engine")
		}
		if *enginePtr != "postgres" && *enginePtr != "mysql" {
			return errors.New("invalid engine: must be postgres or mysql")
		}
		if *portPtr == 0 && *enginePtr == "postgres" {
			*portPtr = 5432
		} else if *portPtr == 0 && *enginePtr == "mysql" {
			*portPtr = 3306
		}

		a.host = *hostPtr
		a.port = *portPtr
		a.region = *regionPtr
		a.user = *userPtr
		a.database = *databasePtr
		a.engine = *enginePtr

		return nil
	}
}

func WithAuthTokenBuilder(authTokenBuilder authtoken.Builder) option {
	return func(a *authenticator) error {
		a.authTokenBuilder = authTokenBuilder
		return nil
	}
}

func (a authenticator) PrintConnectionString() error {
	endpoint := fmt.Sprintf("%s:%d", a.host, a.port)

	token, err := a.authTokenBuilder.BuildToken(context.TODO(), endpoint, a.region, a.user)
	if err != nil {
		return err
	}

	switch a.engine {
	case "postgres":
		fmt.Fprintf(a.output, "postgres://%s:%s@%s/%s", a.user, token, endpoint, a.database)
	case "mysql":
		fmt.Fprintf(a.output, "%s:%s@tcp(%s)/%s?tls=true&allowCleartextPasswords=true", a.user, token, endpoint, a.database)
	}

	return nil
}

func PrintConnectionString() {
	tokenBuilder, err := authtoken.NewRDSTokenBuilder(context.TODO())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	auth, err := NewAuthenticator(
		FromArgs(os.Args[1:]),
		WithAuthTokenBuilder(tokenBuilder),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = auth.PrintConnectionString()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

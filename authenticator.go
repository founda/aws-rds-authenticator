package authenticator

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/founda/aws-rds-authenticator/v2/pkg/authtoken"
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
	sslMode          string
	rootCertFilePath string
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
		portPtr := fset.Int("port", 0, `Port number used for connecting to your DB instance
			default postgres: 5432
			default mysql: 3306`)
		regionPtr := fset.String("region", "us-east-1", "AWS Region where the database instance is running")
		userPtr := fset.String("user", "", "Database account that you want to access")
		databasePtr := fset.String("database", "", "Database that you want to access (optional)")
		enginePtr := fset.String("engine", "postgres", "Database engine that you want to access: postgres|mysql")
		sslModePtr := fset.String("ssl-mode", "", `SSL mode to connect to the database instance.
			postgres: disable|require|verify-ca|verify-full (default: verify-ca)
			mysql: DISABLED|PREFERRED|REQUIRED|VERIFY_CA (default: VERIFY_CA)`)
		rootCertFilePathPtr := fset.String("root-cert-file", "", "Path to the root certificate file")

		err := fset.Parse(args)
		if err != nil {
			return err
		}

		//TODO: add more validation
		if *hostPtr == "" {
			return errors.New("missing required host")
		}
		if *userPtr == "" {
			return errors.New("missing required user")
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
		if *sslModePtr == "" && *enginePtr == "postgres" {
			*sslModePtr = "verify-ca"
		} else if *sslModePtr == "" && *enginePtr == "mysql" {
			*sslModePtr = "VERIFY_CA"
		}

		var validSSLMode bool

		switch *enginePtr {
		case "postgres":
			switch *sslModePtr {
			case "disable", "require", "verify-ca", "verify-full":
				validSSLMode = true
			}
		case "mysql":
			switch *sslModePtr {
			case "DISABLED", "PREFERRED", "REQUIRED", "VERIFY_CA":
				validSSLMode = true
			}
		}

		if !validSSLMode {
			return fmt.Errorf("invalid ssl-mode: must be one of %v", getValidSSLMode(*enginePtr))
		}

		if (*sslModePtr == "verify-ca" || *sslModePtr == "verify-full" || *sslModePtr == "VERIFY_CA") && *rootCertFilePathPtr == "" {
			return fmt.Errorf("root certificate file path is required for ssl-mode %q", *sslModePtr)
		}

		if *rootCertFilePathPtr != "" {
			if _, err := os.Stat(*rootCertFilePathPtr); err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("root certificate file path does not exist: %s", *rootCertFilePathPtr)
				}
				return fmt.Errorf("error checking for root certificate file: %v", err)
			}
		}

		a.host = *hostPtr
		a.port = *portPtr
		a.region = *regionPtr
		a.user = *userPtr
		a.database = *databasePtr
		a.engine = *enginePtr
		a.sslMode = *sslModePtr
		a.rootCertFilePath = *rootCertFilePathPtr

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
		q := make(url.Values)
		q.Set("sslmode", a.sslMode)

		if a.sslMode == "verify-ca" || a.sslMode == "verify-full" {
			q.Set("sslrootcert", a.rootCertFilePath)
		}

		u := url.URL{
			Scheme:   "postgres",
			User:     url.UserPassword(a.user, token),
			Host:     fmt.Sprintf("%s:%d", a.host, a.port),
			Path:     a.database,
			RawQuery: q.Encode(),
		}

		fmt.Fprintf(a.output, "%s", u.String())

	case "mysql":
		params := []string{
			"allowCleartextPasswords=true",
			fmt.Sprintf("ssl-mode=%s", a.sslMode),
		}

		if a.sslMode == "VERIFY_CA" {
			params = append(params, fmt.Sprintf("ssl-ca=%s", a.rootCertFilePath))
		}

		fmt.Fprintf(a.output, "%s:%s@tcp(%s)/%s?%s", a.user, token, endpoint, a.database, strings.Join(params, "&"))
	}

	return nil
}

func PrintConnectionString() {
	var tokenBuilder authtoken.Builder
	var err error

	if passwd := os.Getenv("PG_PASSWORD"); passwd != "" {
		tokenBuilder = authtoken.NewTokenBuilder(passwd)
	} else {
		tokenBuilder, err = authtoken.NewRDSTokenBuilder(context.TODO())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
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

func getValidSSLMode(engine string) []string {
	if engine == "postgres" {
		return []string{"disable", "require", "verify-ca", "verify-full"}
	}
	return []string{"DISABLED", "PREFERRED", "REQUIRED", "VERIFY_CA"}
}

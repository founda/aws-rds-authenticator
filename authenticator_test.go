package authenticator_test

import (
	"bytes"
	"net/url"
	"os"
	"testing"

	authenticator "github.com/founda/aws-rds-authenticator/v2"
	"github.com/founda/aws-rds-authenticator/v2/pkg/authtoken/mock"
)

func TestPrintsConnectionStringToWriter(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "aws-rds-authenticator-test")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile.Close()

	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "postgres with database name",
			args: []string{"-host", "rds.amazon.com", "-user", "postgres", "-region", "eu-west-1", "-database", "prod-test", "-ssl-mode", "verify-full", "-root-cert-file", tmpFile.Name()},
			want: "postgres://postgres:t0k3n@rds.amazon.com:5432/prod-test?sslmode=verify-full&sslrootcert=" + url.QueryEscape(tmpFile.Name()),
		},
		{
			name: "postgres without database name",
			args: []string{"-host", "rds.amazon.com", "-user", "postgres", "-region", "eu-west-1", "-ssl-mode", "verify-full", "-root-cert-file", tmpFile.Name()},
			want: "postgres://postgres:t0k3n@rds.amazon.com:5432?sslmode=verify-full&sslrootcert=" + url.QueryEscape(tmpFile.Name()),
		},
		{
			name: "postgres without database name, disable ssl",
			args: []string{"-host", "rds.amazon.com", "-user", "postgres", "-region", "eu-west-1", "-ssl-mode", "disable"},
			want: "postgres://postgres:t0k3n@rds.amazon.com:5432?sslmode=disable",
		},
		{
			name: "postgres without database name, require ssl",
			args: []string{"-host", "rds.amazon.com", "-user", "postgres", "-region", "eu-west-1", "-ssl-mode", "require"},
			want: "postgres://postgres:t0k3n@rds.amazon.com:5432?sslmode=require",
		},
		{
			name: "postgres without database name, verify-ca ssl",
			args: []string{"-host", "rds.amazon.com", "-user", "postgres", "-region", "eu-west-1", "-ssl-mode", "verify-ca", "-root-cert-file", tmpFile.Name()},
			want: "postgres://postgres:t0k3n@rds.amazon.com:5432?sslmode=verify-ca&sslrootcert=" + url.QueryEscape(tmpFile.Name()),
		},
		{
			name: "mysql with database name",
			args: []string{"-engine", "mysql", "-host", "rds.amazon.com", "-user", "maria", "-region", "eu-west-1", "-database", "prod-test", "-root-cert-file", tmpFile.Name()},
			want: "maria:t0k3n@tcp(rds.amazon.com:3306)/prod-test?allowCleartextPasswords=true&ssl-mode=VERIFY_CA&ssl-ca=" + tmpFile.Name(),
		},
		{
			name: "mysql without database name",
			args: []string{"-engine", "mysql", "-host", "rds.amazon.com", "-user", "maria", "-region", "eu-west-1", "-root-cert-file", tmpFile.Name()},
			want: "maria:t0k3n@tcp(rds.amazon.com:3306)/?allowCleartextPasswords=true&ssl-mode=VERIFY_CA&ssl-ca=" + tmpFile.Name(),
		},
		{
			name: "mysql without database name, explicit enable ssl",
			args: []string{"-engine", "mysql", "-host", "rds.amazon.com", "-user", "maria", "-region", "eu-west-1", "-ssl-mode", "VERIFY_CA", "-root-cert-file", tmpFile.Name()},
			want: "maria:t0k3n@tcp(rds.amazon.com:3306)/?allowCleartextPasswords=true&ssl-mode=VERIFY_CA&ssl-ca=" + tmpFile.Name(),
		},
		{
			name: "mysql without database name, disable ssl",
			args: []string{"-engine", "mysql", "-host", "rds.amazon.com", "-user", "maria", "-region", "eu-west-1", "-ssl-mode", "DISABLED"},
			want: "maria:t0k3n@tcp(rds.amazon.com:3306)/?allowCleartextPasswords=true&ssl-mode=DISABLED",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fakeTerminal := &bytes.Buffer{}

			mockTokenBuilder := mock.NewMockTokenBuilder()
			auth, err := authenticator.NewAuthenticator(
				authenticator.WithOutput(fakeTerminal),
				authenticator.FromArgs(tt.args),
				authenticator.WithAuthTokenBuilder(mockTokenBuilder),
			)
			if err != nil {
				t.Fatal(err)
			}
			err = auth.PrintConnectionString()
			if err != nil {
				t.Fatal(err)
			}

			got := fakeTerminal.String()
			if tt.want != got {
				t.Errorf("want %s, got %s", tt.want, got)
			}
		})
	}
}

func TestMissingRequiredArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "missing host",
			args: []string{"-user", "user", "-region", "eu-west-1", "-database", "db"},
			want: "missing required host",
		},
		{
			name: "missing user",
			args: []string{"-host", "host", "-region", "eu-west-1", "-database", "db"},
			want: "missing required user",
		},
		{
			name: "missing region",
			args: []string{"-host", "host", "-user", "user", "-database", "db"},
			want: "missing required region",
		},
		{
			name: "incorrect engine",
			args: []string{"-host", "host", "-user", "user", "-region", "eu-west-1", "-database", "db", "-engine", "oracle"},
			want: "invalid engine: must be postgres or mysql",
		},
		{
			name: "incorrect ssl mode",
			args: []string{"-host", "host", "-user", "user", "-region", "eu-west-1", "-database", "db", "-ssl-mode", "invalid"},
			want: "invalid ssl-mode: must be one of [disable require verify-ca verify-full]",
		},
		{
			name: "incorrect ssl mode for mysql",
			args: []string{"-host", "host", "-user", "user", "-region", "eu-west-1", "-database", "db", "-ssl-mode", "invalid", "-engine", "mysql"},
			want: "invalid ssl-mode: must be one of [DISABLED PREFERRED REQUIRED VERIFY_CA]",
		},
		{
			name: "missing root cert file",
			args: []string{"-host", "host", "-user", "user", "-region", "eu-west-1", "-database", "db", "-ssl-mode", "verify-ca"},
			want: "root certificate file path is required for ssl-mode \"verify-ca\"",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, gotError := authenticator.NewAuthenticator(
				authenticator.FromArgs(tt.args),
			)
			got := gotError.Error()
			if tt.want != got {
				t.Errorf("want %q\ngot %q", tt.want, got)
			}
		})
	}
}

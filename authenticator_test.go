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

func TestPrintsClientCertificateConnectionString(t *testing.T) {
	t.Parallel()

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
			name: "postgres with client certificate",
			args: []string{"-host", "cnpg-rw.cnpg.svc", "-user", "dicom-job-manager", "-database", "dicom_job_manager", "-ssl-mode", "verify-full", "-root-cert-file", tmpFile.Name(), "-ssl-cert", tmpFile.Name(), "-ssl-key", tmpFile.Name()},
			want: "postgres://dicom-job-manager@cnpg-rw.cnpg.svc:5432/dicom_job_manager?sslcert=" + url.QueryEscape(tmpFile.Name()) + "&sslkey=" + url.QueryEscape(tmpFile.Name()) + "&sslmode=verify-full&sslrootcert=" + url.QueryEscape(tmpFile.Name()),
		},
		{
			name: "postgres with client certificate, no database name",
			args: []string{"-host", "cnpg-rw.cnpg.svc", "-user", "dicom-job-manager", "-ssl-mode", "verify-full", "-root-cert-file", tmpFile.Name(), "-ssl-cert", tmpFile.Name(), "-ssl-key", tmpFile.Name()},
			want: "postgres://dicom-job-manager@cnpg-rw.cnpg.svc:5432?sslcert=" + url.QueryEscape(tmpFile.Name()) + "&sslkey=" + url.QueryEscape(tmpFile.Name()) + "&sslmode=verify-full&sslrootcert=" + url.QueryEscape(tmpFile.Name()),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fakeTerminal := &bytes.Buffer{}

			auth, err := authenticator.NewAuthenticator(
				authenticator.WithOutput(fakeTerminal),
				authenticator.FromArgs(tt.args),
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

func TestClientCertificateValidation(t *testing.T) {
	t.Parallel()

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
			name: "ssl-cert without ssl-key",
			args: []string{"-host", "host", "-user", "user", "-ssl-mode", "require", "-ssl-cert", tmpFile.Name()},
			want: "ssl-cert and ssl-key must be provided together",
		},
		{
			name: "client certificate with mysql",
			args: []string{"-engine", "mysql", "-host", "host", "-user", "user", "-ssl-mode", "REQUIRED", "-ssl-cert", tmpFile.Name(), "-ssl-key", tmpFile.Name()},
			want: "client certificate authentication is only supported for engine postgres",
		},
		{
			name: "client certificate with ssl disabled",
			args: []string{"-host", "host", "-user", "user", "-ssl-mode", "disable", "-ssl-cert", tmpFile.Name(), "-ssl-key", tmpFile.Name()},
			want: `client certificate authentication requires an ssl-mode other than "disable"`,
		},
		{
			name: "missing client certificate file",
			args: []string{"-host", "host", "-user", "user", "-ssl-mode", "require", "-ssl-cert", "/does/not/exist.crt", "-ssl-key", tmpFile.Name()},
			want: "client certificate file path does not exist: /does/not/exist.crt",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, gotError := authenticator.NewAuthenticator(
				authenticator.FromArgs(tt.args),
			)
			if gotError == nil {
				t.Fatalf("want error %q, got none", tt.want)
			}
			if got := gotError.Error(); tt.want != got {
				t.Errorf("want %q\ngot %q", tt.want, got)
			}
		})
	}
}

func TestMissingAuthTokenBuilder(t *testing.T) {
	t.Parallel()

	auth, err := authenticator.NewAuthenticator(
		authenticator.WithOutput(&bytes.Buffer{}),
		authenticator.FromArgs([]string{"-host", "host", "-user", "user", "-ssl-mode", "require"}),
	)
	if err != nil {
		t.Fatal(err)
	}

	gotError := auth.PrintConnectionString()
	if gotError == nil {
		t.Fatal("want an error, got none")
	}
	if want := "no auth token builder configured"; gotError.Error() != want {
		t.Errorf("want %q, got %q", want, gotError.Error())
	}
}

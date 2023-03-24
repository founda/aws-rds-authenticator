package authenticator_test

import (
	"bytes"
	"testing"

	authenticator "github.com/founda/aws-rds-authenticator"
	"github.com/founda/aws-rds-authenticator/pkg/authtoken/mock"
)

func TestPrintsConnectionStringToWriterPostgres(t *testing.T) {
	t.Parallel()
	fakeTerminal := &bytes.Buffer{}
	args := []string{"-host", "rds.amazon.com", "-user", "postgres", "-region", "eu-west-1", "-database", "prod-test"}

	mockTokenBuilder := mock.NewMockTokenBuilder()
	auth, err := authenticator.NewAuthenticator(
		authenticator.WithOutput(fakeTerminal),
		authenticator.FromArgs(args),
		authenticator.WithAuthTokenBuilder(mockTokenBuilder),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = auth.PrintConnectionString()
	if err != nil {
		t.Fatal(err)
	}
	want := "postgres://postgres:t0k3n@rds.amazon.com:5432/prod-test"
	got := fakeTerminal.String()
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestPrintsConnectionStringToWriterMySQL(t *testing.T) {
	t.Parallel()
	fakeTerminal := &bytes.Buffer{}
	args := []string{"-engine", "mysql", "-host", "rds.amazon.com", "-user", "maria", "-region", "eu-west-1", "-database", "prod-test"}

	mockTokenBuilder := mock.NewMockTokenBuilder()
	auth, err := authenticator.NewAuthenticator(
		authenticator.WithOutput(fakeTerminal),
		authenticator.FromArgs(args),
		authenticator.WithAuthTokenBuilder(mockTokenBuilder),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = auth.PrintConnectionString()
	if err != nil {
		t.Fatal(err)
	}
	want := "maria:t0k3n@tcp(rds.amazon.com:3306)/prod-test?tls=true&allowCleartextPasswords=true"
	got := fakeTerminal.String()
	if want != got {
		t.Errorf("want %q, got %q", want, got)
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
			name: "missing database",
			args: []string{"-host", "host", "-user", "user", "-region", "eu-west-1"},
			want: "missing required database",
		},
		{
			name: "incorrect engine",
			args: []string{"-host", "host", "-user", "user", "-region", "eu-west-1", "-database", "db", "-engine", "oracle"},
			want: "invalid engine: must be postgres or mysql",
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
				t.Errorf("want %q, got %q", tt.want, got)
			}
		})
	}
}

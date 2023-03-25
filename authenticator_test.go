package authenticator_test

import (
	"bytes"
	"testing"

	authenticator "github.com/founda/aws-rds-authenticator"
	"github.com/founda/aws-rds-authenticator/pkg/authtoken/mock"
	"github.com/google/go-cmp/cmp"
)

func TestPrintsConnectionStringToWriter(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "postgres with database name",
			args: []string{"-host", "rds.amazon.com", "-user", "postgres", "-region", "eu-west-1", "-database", "prod-test"},
			want: "postgres://postgres:t0k3n@rds.amazon.com:5432/prod-test",
		},
		{
			name: "postgres without database name",
			args: []string{"-host", "rds.amazon.com", "-user", "postgres", "-region", "eu-west-1"},
			want: "postgres://postgres:t0k3n@rds.amazon.com:5432/",
		},
		{
			name: "mysql with database name",
			args: []string{"-engine", "mysql", "-host", "rds.amazon.com", "-user", "maria", "-region", "eu-west-1", "-database", "prod-test"},
			want: "maria:t0k3n@tcp(rds.amazon.com:3306)/prod-test?tls=true&allowCleartextPasswords=true",
		},
		{
			name: "mysql without database name",
			args: []string{"-engine", "mysql", "-host", "rds.amazon.com", "-user", "maria", "-region", "eu-west-1"},
			want: "maria:t0k3n@tcp(rds.amazon.com:3306)/?tls=true&allowCleartextPasswords=true",
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
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("PrintConnectionString() mismatch (-want +got):\n%s", diff)
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

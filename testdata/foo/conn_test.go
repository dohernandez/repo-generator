package foo_test

import (
	"database/sql"
	"embed"
	"testing"

	"github.com/dohernandez/repo-generator/postgres"
)

//go:embed testdata/migrations/*.sql
var migrations embed.FS

func postgresConnect(t *testing.T) (*sql.DB, error) {
	return postgres.ConnectForTesting(t, migrations)
}

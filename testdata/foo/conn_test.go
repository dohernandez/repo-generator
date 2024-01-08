package foo_test

import (
	"database/sql"
	"embed"
	"testing"

	"github.com/dohernandez/storage/postgres"
)

//go:embed testdata/migrations/*.sql
var migrations embed.FS

func postgresConnect(t *testing.T) (*sql.DB, error) {
	t.Helper()

	return postgres.ConnectForTesting(t, migrations)
}

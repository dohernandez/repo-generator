package postgres_test

import (
	"context"
	"embed"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dohernandez/repo-generator/postgres"
)

//go:embed testdata/migrations/*.sql
var migrations embed.FS

func TestConnectForTesting(t *testing.T) {
	t.Parallel()

	t.Run("ping context", func(t *testing.T) {
		t.Parallel()

		conn, err := postgres.ConnectForTesting(t, migrations)
		require.NoError(t, err)
		require.NotNil(t, conn)

		err = conn.PingContext(context.Background())
		require.NoError(t, err)
	})

	t.Run("select version", func(t *testing.T) {
		t.Parallel()

		conn, err := postgres.ConnectForTesting(t, migrations)
		require.NoError(t, err)
		require.NotNil(t, conn)

		var version string

		err = conn.QueryRowContext(context.Background(), "SELECT VERSION()").Scan(&version)
		require.NoError(t, err)

		require.NotEmpty(t, version)
	})

	t.Run("cleanup occurs", func(t *testing.T) {
		t.Parallel()

		conn, err := postgres.ConnectForTesting(t, migrations)
		require.NoError(t, err)
		require.NotNil(t, conn)

		_, err = conn.ExecContext(context.Background(), "INSERT INTO tests (id, timestamp) VALUES ($1, $2)",
			uuid.NewString(), time.Now())
		require.NoError(t, err)

		var count int

		err = conn.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM tests").Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count)

		// Clean up the connection and the database it was using.
		err = conn.Close()
		require.NoError(t, err)

		// This should be a new connection, with an empty schema.
		conn, err = postgres.ConnectForTesting(t, migrations)
		require.NoError(t, err)
		require.NotNil(t, conn)

		err = conn.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM tests").Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})
}

func TestConfig_URI(t *testing.T) {
	t.Parallel()

	t.Run("simple config", func(t *testing.T) {
		t.Parallel()

		c := postgres.Config{
			Host:     "postgres",
			Port:     5432,
			User:     "tests",
			Password: "tests",
			Database: "postgres",
			Flags:    "sslmode=disable&client_encoding=UTF8",
		}

		assert.Equal(t, "postgresql://tests:tests@postgres:5432/postgres?sslmode=disable&client_encoding=UTF8", c.URI())

		_, err := pgx.ParseConfig(c.URI())
		require.NoError(t, err)
	})

	t.Run("port override in the host string", func(t *testing.T) {
		t.Parallel()

		c := postgres.Config{
			Host:     "postgres:4321",
			Port:     5432,
			User:     "tests",
			Password: "tests",
			Database: "postgres",
			Flags:    "sslmode=disable&client_encoding=UTF8",
		}

		assert.Equal(t, "postgresql://tests:tests@postgres:4321/postgres?sslmode=disable&client_encoding=UTF8", c.URI())

		_, err := pgx.ParseConfig(c.URI())
		require.NoError(t, err)
	})

	t.Run("with schema specified and no flags", func(t *testing.T) {
		t.Parallel()

		c := postgres.Config{
			Host:     "postgres:4321",
			Port:     5432,
			User:     "tests",
			Password: "tests",
			Schema:   "ethereum_mainnet",
			Database: "postgres",
		}

		assert.Equal(t, "postgresql://tests:tests@postgres:4321/postgres?search_path=ethereum_mainnet", c.URI())

		_, err := pgx.ParseConfig(c.URI())
		require.NoError(t, err)
	})

	t.Run("with schema specified and has flags", func(t *testing.T) {
		t.Parallel()

		c := postgres.Config{
			Host:     "postgres:4321",
			Port:     5432,
			User:     "tests",
			Password: "tests",
			Schema:   "ethereum_mainnet",
			Database: "postgres",
			Flags:    "sslmode=disable&client_encoding=UTF8",
		}

		assert.Equal(t, "postgresql://tests:tests@postgres:4321/postgres?search_path=ethereum_mainnet&sslmode=disable&client_encoding=UTF8", c.URI())

		_, err := pgx.ParseConfig(c.URI())
		require.NoError(t, err)
	})
}

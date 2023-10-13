package postgres

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host            string `default:"localhost"`
	Port            int    `default:"5432"`
	User            string
	Password        string
	Database        string `default:"postgres"`
	Schema          string `default:"public"`
	Flags           string
	MaxConnLifetime time.Duration `default:"4h"`
	MaxIdleConns    int           `default:"20"`
	MaxOpenConns    int           `default:"20"`
}

// URI parses a Config struct into a connection URI.
//
// postgresql://[user[:password]@][netloc][:port][,...][/dbname][?search_path=schema&param1=value1&...]
func (c *Config) URI() string {
	var builder strings.Builder

	builder.WriteString("postgresql://")

	if c.User != "" {
		builder.WriteString(c.User)

		if c.Password != "" {
			builder.WriteString(":" + c.Password)
		}

		builder.WriteString("@")
	}

	var uriParts []string

	if c.Host != "" {
		uriParts = strings.Split(c.Host, ":")

		builder.WriteString(uriParts[0])

		// A port provided in the host string takes precedence over one explicitly defined.
		if len(uriParts) > 1 {
			builder.WriteString(fmt.Sprintf(":%s", uriParts[1]))
		}
	}

	// We only write the port if one was not provided in the host string.
	if c.Port != 0 && len(uriParts) == 1 {
		builder.WriteString(fmt.Sprintf(":%d", c.Port))
	}

	if c.Database != "" {
		builder.WriteString("/" + c.Database)
	}

	flagsPrefix := "?"

	if c.Schema != "" {
		flagsPrefix = "&"

		builder.WriteString(fmt.Sprintf("?search_path=%s", c.Schema))
	}

	if c.Flags != "" {
		builder.WriteString(flagsPrefix)
		builder.WriteString(c.Flags)
	}

	return builder.String()
}

func ConnectWithRetry(ctx context.Context, conf *Config) (*sql.DB, error) {
	var (
		db  *sql.DB
		err error
	)

	fn := func() error {
		db, err = Connect(ctx, conf)
		if err != nil {
			return fmt.Errorf("connect: %w", err)
		}

		return nil
	}

	if err := backoff.Retry(fn, backoff.WithMaxRetries(backoff.NewConstantBackOff(5*time.Second), 3)); err != nil {
		return nil, err
	}

	return db, nil
}

func Connect(ctx context.Context, conf *Config) (*sql.DB, error) {
	connConfig, err := pgx.ParseConfig(conf.URI())
	if err != nil {
		return nil, fmt.Errorf("pgx parse config: %w", err)
	}

	connStr := stdlib.RegisterConnConfig(connConfig)

	driverName := "pgx"

	db, err := sql.Open(driverName, connStr)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	db.SetConnMaxLifetime(conf.MaxConnLifetime)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetMaxOpenConns(conf.MaxOpenConns)

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping context: %w", err)
	}

	return db, nil
}

type testConfig struct {
	Config *Config `envconfig:"POSTGRES_TEST" split_words:"true"`
}

func ConnectForTesting(t *testing.T, fs embed.FS) (*sql.DB, error) {
	t.Helper()

	ctx := context.Background()

	var conf testConfig
	if err := envconfig.Process("", &conf); err != nil {
		return nil, fmt.Errorf("process config: %w", err)
	}

	conn, err := ConnectWithRetry(ctx, conf.Config)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	testDatabase := fmt.Sprintf("%s_test_%d", conf.Config.Database, rand.Int63()) //nolint: gosec

	if err := createDatabase(ctx, t, conn, testDatabase); err != nil {
		return nil, fmt.Errorf("create database: %w", err)
	}

	conf.Config.Database = testDatabase

	testConn, err := Connect(ctx, conf.Config)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	t.Cleanup(func() {
		// Close the connection to the test database.
		_ = testConn.Close() //nolint: errcheck

		// Drop the database that we were currently testing against.
		_ = dropDatabase(ctx, t, conn, testDatabase) //nolint: errcheck

		// Close the main connection.
		_ = conn.Close() //nolint: errcheck
	})

	if err := MigrateUp(ctx, t, testConn, fs); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return testConn, nil
}

func createDatabase(ctx context.Context, t *testing.T, conn *sql.DB, database string) error {
	t.Helper()

	_, err := conn.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s;", database))
	if err != nil {
		return fmt.Errorf("%w: exec context", err)
	}

	return nil
}

func dropDatabase(ctx context.Context, t *testing.T, conn *sql.DB, database string) error {
	t.Helper()

	_, err := conn.ExecContext(ctx, fmt.Sprintf("DROP DATABASE %s;", database))
	if err != nil {
		return fmt.Errorf("%w: exec context", err)
	}

	return nil
}

// MigrateUp runs all migrations in the provided filesystem.
func MigrateUp(ctx context.Context, t *testing.T, conn *sql.DB, fs embed.FS) error {
	t.Helper()

	// Queue to keep track of the directories that we come across during the search.
	dirQueue := []string{"."}

	// Perform some sort of bfs across the directory tree to find all migration files, and execute them.
	for len(dirQueue) > 0 {
		dir := dirQueue[0]
		dirQueue = dirQueue[1:] // pop the current directory off the queue.

		entries, err := fs.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("read dir: %w", err)
		}

		for _, entry := range entries {
			// All files are read relative to the root, we need to keep track of the full path of the files here.
			path := filepath.Join(dir, entry.Name())

			if entry.IsDir() {
				// Store the directory for later, and move on.
				dirQueue = append(dirQueue, path)

				continue
			}

			if !strings.HasSuffix(entry.Name(), "up.sql") {
				// We only want "up" migrations.
				continue
			}

			m, err := fs.ReadFile(path)
			if err != nil {
				return fmt.Errorf("read migration file: %w", err)
			}

			if _, err := conn.ExecContext(ctx, string(m)); err != nil {
				return fmt.Errorf("exec context: %w", err)
			}
		}
	}

	return nil
}

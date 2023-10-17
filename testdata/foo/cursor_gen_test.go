package foo_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/dohernandez/repo-generator/testdata/foo"
)

const cursorTable = "cursor"

func TestCursorRepo_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("create successfully", func(t *testing.T) {
		t.Parallel()

		conn, err := postgresConnect(t)
		require.NoError(t, err)
		require.NotEmpty(t, conn)

		repo := foo.NewCursorRepo(conn, cursorTable)

		id := uuid.New()
		now := time.Now()

		c, err := repo.Create(ctx, &foo.Cursor{
			ID:        id,
			Name:      "test-cursor",
			CreatedAt: now,
			UpdatedAt: now,
		})
		require.NoError(t, err)
		require.NotEmpty(t, c)

		require.Equal(t, id, c.ID)
		require.Equal(t, "test-cursor", c.Name)
		require.Equal(t, c.CreatedAt, now)
		require.Equal(t, c.UpdatedAt, now)

		// Check that all transfers have been inserted.
		query := `SELECT count(*) FROM cursor`

		row := conn.QueryRowContext(ctx, query)
		require.NoError(t, err)

		var count int

		err = row.Scan(&count)
		require.NoError(t, err)

		require.Equal(t, 1, count)
	})
}

func TestCursorRepo_Insert(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("insert successfully", func(t *testing.T) {
		t.Parallel()

		conn, err := postgresConnect(t)
		require.NoError(t, err)
		require.NotEmpty(t, conn)

		repo := foo.NewCursorRepo(conn, cursorTable)

		err = repo.Insert(ctx,
			&foo.Cursor{
				ID:        uuid.New(),
				Name:      "test-cursor-1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			&foo.Cursor{
				ID:        uuid.New(),
				Name:      "test-cursor-2",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			&foo.Cursor{
				ID:        uuid.New(),
				Name:      "test-cursor-3",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		)
		require.NoError(t, err)

		// Check that all transfers have been inserted.
		query := `SELECT count(*) FROM cursor`

		row := conn.QueryRowContext(ctx, query)
		require.NoError(t, err)

		var count int

		err = row.Scan(&count)
		require.NoError(t, err)

		require.Equal(t, 3, count)
	})
}

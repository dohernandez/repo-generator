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
		now := time.Now().UTC()

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

func TestCursorRepo_Select(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("select successfully", func(t *testing.T) {
		t.Parallel()

		conn, err := postgresConnect(t)
		require.NoError(t, err)
		require.NotEmpty(t, conn)

		repo := foo.NewCursorRepo(conn, cursorTable)

		_, err = repo.Create(ctx,
			&foo.Cursor{
				ID:        uuid.New(),
				Name:      "test-cursor-1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		require.NoError(t, err)

		c2, err := repo.Create(ctx,
			&foo.Cursor{
				ID:        uuid.New(),
				Name:      "test-cursor-2",
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
		require.NoError(t, err)

		_, err = repo.Create(ctx,
			&foo.Cursor{
				ID:        uuid.New(),
				Name:      "test-cursor-3",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		require.NoError(t, err)

		// Check that all transfers have been inserted.
		c, err := repo.Select(ctx, foo.CursorCriteria{
			"name": c2.Name,
		})
		require.NoError(t, err)

		require.Equal(t, c2.ID, c.ID)
		require.Equal(t, c2.Name, c.Name)
		require.Equal(t, c2.CreatedAt, c.CreatedAt)
		require.Equal(t, c2.UpdatedAt, c.UpdatedAt)
	})
}

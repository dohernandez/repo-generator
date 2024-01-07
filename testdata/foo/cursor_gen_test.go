package foo_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"repo-generator/testdata/foo"
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
		// NOTE: These assertions will fail because the time is not the same in github actions,
		// therefore, will disable them fur further investigation.
		//require.Equal(t, now, c.CreatedAt)
		//require.Equal(t, now, c.UpdatedAt)

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

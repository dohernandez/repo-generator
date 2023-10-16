package foo

import (
	"context"
	"embed"
	"math/big"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/dohernandez/repo-generator/postgres"
	"github.com/dohernandez/repo-generator/testdata/deps"
)

const table = "block"

//go:embed testdata/migrations/*.sql
var migrations embed.FS

func TestBlockRepo_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("insert successfully", func(t *testing.T) {
		t.Parallel()

		conn, err := postgres.ConnectForTesting(t, migrations)
		require.NoError(t, err)
		require.NotEmpty(t, conn)

		repo := NewBlockRepo(conn, table)

		id := uuid.New()
		now := time.Now()

		b, err := repo.Create(ctx, &Block{
			ID:             id,
			Hash:           deps.HexToHash("0x0"),
			Number:         big.NewInt(0),
			ChainID:        deps.EthereumChainID,
			BlockTimestamp: now.UTC(),
		})
		require.NoError(t, err)
		require.NotEmpty(t, b)

		require.Equal(t, id, b.ID)
		require.Equal(t, deps.EthereumChainID, b.ChainID)
		require.Equal(t, big.NewInt(0), b.Number)
		require.Equal(t, deps.HexToHash("0x0"), b.Hash)
		require.Equal(t, b.BlockTimestamp, now.UTC())
	})
}

func TestBlockRepo_Insert(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("insert successfully", func(t *testing.T) {
		t.Parallel()

		conn, err := postgres.ConnectForTesting(t, migrations)
		require.NoError(t, err)
		require.NotEmpty(t, conn)

		repo := NewBlockRepo(conn, table)

		err = repo.Insert(ctx,
			&Block{
				ID:             uuid.New(),
				Hash:           deps.HexToHash("0x0"),
				Number:         big.NewInt(0),
				ChainID:        deps.EthereumChainID,
				BlockTimestamp: time.Now(),
			},
			&Block{
				ID:             uuid.New(),
				Hash:           deps.HexToHash("0x1"),
				Number:         big.NewInt(1),
				ChainID:        deps.EthereumChainID,
				BlockTimestamp: time.Now(),
				ParentHash:     deps.HexToHash("0x0"),
			},
			&Block{
				ID:             uuid.New(),
				Hash:           deps.HexToHash("0x2"),
				Number:         big.NewInt(2),
				ChainID:        deps.EthereumChainID,
				BlockTimestamp: time.Now(),
				ParentHash:     deps.HexToHash("0x1"),
			},
		)
		require.NoError(t, err)

		// Check that all transfers have been inserted.
		query := `SELECT count(*) FROM block`

		row := conn.QueryRowContext(ctx, query)
		require.NoError(t, err)

		var count int

		err = row.Scan(&count)
		require.NoError(t, err)

		require.Equal(t, 3, count)
	})
}

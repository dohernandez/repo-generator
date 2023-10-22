package foo_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/dohernandez/repo-generator/testdata/deps"
	"github.com/dohernandez/repo-generator/testdata/foo"
)

const blockTable = "block"

func TestBlockRepo_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("insert successfully", func(t *testing.T) {
		t.Parallel()

		conn, err := postgresConnect(t)
		require.NoError(t, err)
		require.NotEmpty(t, conn)

		repo := foo.NewBlockRepo(conn, blockTable)

		id := uuid.New()
		now := time.Now()

		b, err := repo.Create(ctx, &foo.Block{
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

		conn, err := postgresConnect(t)
		require.NoError(t, err)
		require.NotEmpty(t, conn)

		repo := foo.NewBlockRepo(conn, blockTable)

		err = repo.Insert(ctx,
			&foo.Block{
				ID:             uuid.New(),
				Hash:           deps.HexToHash("0x0"),
				Number:         big.NewInt(0),
				ChainID:        deps.EthereumChainID,
				BlockTimestamp: time.Now(),
			},
			&foo.Block{
				ID:             uuid.New(),
				Hash:           deps.HexToHash("0x1"),
				Number:         big.NewInt(1),
				ChainID:        deps.EthereumChainID,
				BlockTimestamp: time.Now(),
				ParentHash:     deps.HexToHash("0x0"),
			},
			&foo.Block{
				ID:             uuid.New(),
				Hash:           deps.HexToHash("0x2"),
				Number:         big.NewInt(2),
				ChainID:        deps.EthereumChainID,
				BlockTimestamp: time.Now(),
				ParentHash:     deps.HexToHash("0x1"),
			},
		)
		require.NoError(t, err)

		// Check that all blocks have been inserted.
		query := `SELECT count(*) FROM block`

		row := conn.QueryRowContext(ctx, query)
		require.NoError(t, err)

		var count int

		err = row.Scan(&count)
		require.NoError(t, err)

		require.Equal(t, 3, count)
	})
}

func TestBlockRepo_Update(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("update successfully", func(t *testing.T) {
		t.Parallel()

		conn, err := postgresConnect(t)
		require.NoError(t, err)
		require.NotEmpty(t, conn)

		repo := foo.NewBlockRepo(conn, blockTable)

		id := uuid.New()
		now := time.Now()

		b, err := repo.Create(ctx, &foo.Block{
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

		b.Hash = deps.HexToHash("0x2")
		b.Number = big.NewInt(2)

		err = repo.Update(ctx, b, true)
		require.NoError(t, err)

		// Check that the block has been updated.
		query := `SELECT hash, number FROM block WHERE id=$1`

		row := conn.QueryRowContext(ctx, query, id)
		require.NoError(t, err)

		var (
			hash string
			num  int64
		)

		err = row.Scan(&hash, &num)
		require.NoError(t, err)

		require.Equal(t, b.Hash.String(), hash)
		require.Equal(t, b.Number.Int64(), num)

		b.ParentHash = deps.HexToHash("0x1")

		err = repo.Update(ctx, b, true)
		require.NoError(t, err)

		// Check that the block has been updated.
		query = `SELECT parent_hash FROM block WHERE id=$1`

		row = conn.QueryRowContext(ctx, query, id)
		require.NoError(t, err)

		var (
			//chainID    int
			parentHash string
		)

		err = row.Scan(&parentHash)
		require.NoError(t, err)

		require.Equal(t, b.ParentHash.String(), parentHash)
	})
}

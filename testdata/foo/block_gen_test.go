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

		b, err := repo.Create(ctx, &Block{
			ID:             id,
			Hash:           deps.HexToHash("0x0"),
			Number:         big.NewInt(0),
			ChainID:        deps.EthereumChainID,
			BlockTimestamp: time.Now(),
		})
		require.NoError(t, err)
		require.NotEmpty(t, b)

		require.Equal(t, id, b.ID)
		require.Equal(t, deps.EthereumChainID, b.ChainID)
		require.Equal(t, big.NewInt(0), b.Number)
		require.Equal(t, deps.HexToHash("0x0"), b.Hash)
		//require.True(t, b.BlockTimestamp.IsZero())
		println(b.BlockTimestamp.String())
	})
}

package generator_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	generator "github.com/dohernandez/repo-generator"
)

func TestGenerate_Network(t *testing.T) {
	err := generator.Generate(
		"testdata/foo/network.go",
		"testdata/foo/network_gen.go",
		"Network",
	)
	require.NoError(t, err)
}

func TestGenerate_Block(t *testing.T) {
	err := generator.Generate(
		"testdata/foo/block.go",
		"testdata/foo/block_gen.go",
		"Block",
		generator.WithImports([]string{
			//"github.com/google/uuid",
			"github.com/dohernandez/repo-generator/testdata/deps",
		}),
		generator.WithCreateFunc(),
		generator.WithInsertFunc(),
		generator.WithUpdateFunc(),
		generator.WithDeleteFunc(),
	)
	require.NoError(t, err)
}

func TestGenerate_Asset(t *testing.T) {
	err := generator.Generate(
		"testdata/foo/asset.go",
		"testdata/foo/asset_gen.go",
		"Asset",
		generator.WithImports([]string{
			"github.com/dohernandez/repo-generator/testdata/deps",
			"github.com/lib/pq",
		}),
		generator.WithCreateFunc(),
	)
	require.NoError(t, err)
}

func TestGenerate_Transfer(t *testing.T) {
	err := generator.Generate(
		"testdata/foo/transfer.go",
		"testdata/foo/transfer_gen.go",
		"Transfer",
		generator.WithImports([]string{
			"github.com/dohernandez/repo-generator/testdata/deps",
		}),
		generator.WithCreateFunc(),
	)
	require.NoError(t, err)
}

func TestGenerate_Cursor(t *testing.T) {
	err := generator.Generate(
		"testdata/foo/cursor.go",
		"testdata/foo/cursor_gen.go",
		"Cursor",
		generator.WithImports([]string{
			"github.com/google/uuid",
		}),
		generator.WithCreateFunc(),
		generator.WithInsertFunc(),
	)
	require.NoError(t, err)
}

func TestGenerate_Sync(t *testing.T) {
	err := generator.Generate(
		"testdata/foo/sync.go",
		"testdata/foo/sync_gen.go",
		"Sync",
		generator.WithImports([]string{
			"github.com/dohernandez/repo-generator/testdata/deps",
		}),
		generator.WithCreateFunc(),
		generator.WithInsertFunc(),
		generator.WithUpdateFunc(),
		generator.WithDeleteFunc(),
	)
	require.NoError(t, err)
}

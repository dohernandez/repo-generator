package generator_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"

	generator "github.com/dohernandez/repo-generator"
)

func TestParseModel(t *testing.T) {
	model := "Block"
	filename := "testdata/foo/block.go"

	fset := token.NewFileSet()

	tree, err := parser.ParseFile(fset, filename, nil, 0)
	require.NoError(t, err)

	var gm generator.Model

	ast.Inspect(tree, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			// We should recurse on any node higher-level than a TypeSpec.
			// Invokes inspect f recursively
			return true
		}

		if ts.Name.Name != model {
			// We only care about the Model.
			// Invokes inspect f recursively
			return true
		}

		s, ok := ts.Type.(*ast.StructType)
		if !ok {
			t.Fatal("models must be structs")

			// Inspect stop
			return false
		}

		gm, err = generator.ParseModel(s, model)
		require.NoError(t, err)

		// Inspect stop
		return false
	})

	require.Equal(t, "Block", gm.Name)
	require.Len(t, gm.Fields, 6)

	require.Equal(t, "ID", gm.Fields[0].Name)
	require.Equal(t, "uuid.UUID", gm.Fields[0].Type)
	require.False(t, gm.Fields[0].IsPointer)
	require.Equal(t, "id", gm.Fields[0].ColName)
	require.True(t, gm.Fields[0].IsKey)
	require.False(t, gm.Fields[0].OmittedEmpty)
	require.False(t, gm.Fields[0].IsNullable)

	require.Equal(t, "ChainID", gm.Fields[1].Name)
	require.Equal(t, "deps.ChainID", gm.Fields[1].Type)
	require.False(t, gm.Fields[1].IsPointer)
	require.Equal(t, "chain_id", gm.Fields[1].ColName)
	require.False(t, gm.Fields[1].IsKey)
	require.False(t, gm.Fields[1].OmittedEmpty)
	require.False(t, gm.Fields[1].IsNullable)

	require.Equal(t, "Hash", gm.Fields[2].Name)
	require.Equal(t, "deps.Hash", gm.Fields[2].Type)
	require.False(t, gm.Fields[2].IsPointer)
	require.Equal(t, "hash", gm.Fields[2].ColName)
	require.False(t, gm.Fields[2].IsKey)
	require.False(t, gm.Fields[2].OmittedEmpty)
	require.True(t, gm.Fields[2].IsNullable)
	require.Equal(t, "sql.NullString", gm.Fields[2].SQLType)
	require.Equal(t, "deps", gm.Fields[2].Scan.Pkg)
	require.Equal(t, "HexToHash", gm.Fields[2].Scan.Name)
	require.Equal(t, "_", gm.Fields[2].Value.Pkg)
	require.Equal(t, "String", gm.Fields[2].Value.Name)

	require.Equal(t, "Number", gm.Fields[3].Name)
	require.Equal(t, "big.Int", gm.Fields[3].Type)
	require.True(t, gm.Fields[3].IsPointer)
	require.Equal(t, "number", gm.Fields[3].ColName)
	require.False(t, gm.Fields[3].IsKey)
	require.False(t, gm.Fields[3].OmittedEmpty)
	require.False(t, gm.Fields[3].IsNullable)
	require.Equal(t, "int64", gm.Fields[3].SQLType)
	require.Equal(t, "big", gm.Fields[3].Scan.Pkg)
	require.Equal(t, "NewInt", gm.Fields[3].Scan.Name)
	require.Equal(t, "_", gm.Fields[3].Value.Pkg)
	require.Equal(t, "Int64", gm.Fields[3].Value.Name)

	require.Equal(t, "ParentHash", gm.Fields[4].Name)
	require.Equal(t, "deps.Hash", gm.Fields[4].Type)
	require.False(t, gm.Fields[4].IsPointer)
	require.Equal(t, "parent_hash", gm.Fields[4].ColName)
	require.False(t, gm.Fields[4].IsKey)
	require.False(t, gm.Fields[4].OmittedEmpty)
	require.True(t, gm.Fields[4].IsNullable)
	require.Equal(t, "sql.NullString", gm.Fields[4].SQLType)
	require.Equal(t, "deps", gm.Fields[4].Scan.Pkg)
	require.Equal(t, "HexToHash", gm.Fields[4].Scan.Name)
	require.Equal(t, "_", gm.Fields[4].Value.Pkg)
	require.Equal(t, "String", gm.Fields[4].Value.Name)

	require.Equal(t, "BlockTimestamp", gm.Fields[5].Name)
	require.Equal(t, "time.Time", gm.Fields[5].Type)
	require.False(t, gm.Fields[5].IsPointer)
	require.Equal(t, "block_timestamp", gm.Fields[5].ColName)
	require.False(t, gm.Fields[5].IsKey)
	require.True(t, gm.Fields[5].OmittedEmpty)
	require.True(t, gm.Fields[5].IsNullable)
	require.Equal(t, "sql.NullTime", gm.Fields[5].SQLType)
	require.Equal(t, "_", gm.Fields[5].Scan.Pkg)
	require.Equal(t, "UTC", gm.Fields[5].Scan.Name)
	require.Equal(t, "", gm.Fields[5].Value.Pkg)
	require.Equal(t, "", gm.Fields[5].Value.Name)
}

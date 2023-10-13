package generator

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseRepo(t *testing.T) {
	model := "Block"
	filename := "testdata/foo/block.go"

	fset := token.NewFileSet()

	tree, err := parser.ParseFile(fset, filename, nil, 0)
	require.NoError(t, err)

	r, err := ParseRepo(tree, model)
	require.NoError(t, err)

	require.Equal(t, "Block", r.Model.Name)
	require.Len(t, r.Model.Fields, 6)

	require.Len(t, r.ColKeysFields, 1)
	require.Len(t, r.ColStatesFields, 5)
}

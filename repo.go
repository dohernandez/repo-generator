package generator

import (
	"go/ast"
)

type Repo struct {
	Receiver string

	Model Model

	// ColKeysFields are the fields that are keys in the table.
	ColKeysFields []Field
	// ColStatesFields are the fields that are not keys in the table.
	ColStatesFields []Field
}

func ParseRepo(tree *ast.File, model string) (Repo, error) {
	var (
		m   Model
		err error
	)

	ast.Inspect(tree, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			// We should recurse on any node higher-level than a TypeSpec.
			// Invokes inspect f recursively
			return true
		}

		if ts.Name.Name != model {
			// We only care about the Model.
			//
			// Continue inspecting.
			return true
		}

		s, ok := ts.Type.(*ast.StructType)
		if !ok {
			// Inspect stop.
			return false
		}

		m, err = ParseModel(s, model)

		// Inspect stop.
		return false
	})

	var (
		colKeys   []Field
		colStates []Field
	)

	for _, f := range m.Fields {
		if f.IsKey {
			colKeys = append(colKeys, f)

			continue
		}

		colStates = append(colStates, f)
	}

	return Repo{
		Receiver:        "repo",
		Model:           m,
		ColKeysFields:   colKeys,
		ColStatesFields: colStates,
	}, err
}

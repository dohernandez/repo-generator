package generator

import (
	"go/ast"
	"path"
	"strings"

	"github.com/dohernandez/errors"
)

var defaultImports = map[string]string{
	"errors":    "github.com/dohernandez/errors",
	"pgerrcode": "github.com/jackc/pgerrcode",
	"pgconn":    "github.com/jackc/pgx/v5/pgconn",
}

var ErrPackageNotFound = errors.New("package not found")

// PackageImport defines a package which was imported in a Go file.
type PackageImport struct {
	// Alias is the alias of the package. If the package was not imported with an alias, this field will be empty.
	Alias string
	// Path is the import path of the package.
	Path string
}

func parseImports(tree *ast.File, r Repo) (map[string]PackageImport, error) {
	mImports := make(map[string]PackageImport)

	ast.Inspect(tree, func(n ast.Node) bool {
		is, ok := n.(*ast.ImportSpec)
		if !ok {
			// We should recurse on any node higher-level than a TypeSpec.
			// Invokes inspect f recursively
			return true
		}

		ipath := is.Path.Value[1 : len(is.Path.Value)-1]

		k := path.Base(ipath)

		alias := ""

		if is.Name != nil && is.Name.Name != "_" {
			alias = is.Name.Name

			k = alias
		}

		mImports[k] = PackageImport{
			Path:  ipath,
			Alias: alias,
		}

		// Inspect stop.
		return false
	})

	pImports := make(map[string]PackageImport, len(defaultImports))

	for k, i := range defaultImports {
		pImports[k] = PackageImport{
			Path: i,
		}
	}

	// TODO: test
	// 1. arrayable, nullable, no scan method
	// 2. package type (type is defined in another pkg), nullable, scan method
	// 3. package type (type is defined in another pkg), nullable, nil method
	// 4. package type (type is defined in another pkg), nullable, value method
	// 5. package type (type is defined in another pkg), nullable, scan, nil, value method different pkg
	// 6. package type (type is defined in another pkg), nullable, scan is type (pkg is _), nil, value method different pkg
	for _, f := range r.Model.Fields {
		if !f.HasScanMethod && !f.HasValueMethod && !f.HasNilMethod {
			continue
		}

		if f.HasSqlArrayable {
			k := path.Base(arrayablePackageImport.Path)

			pImports[k] = arrayablePackageImport
		}

		//parts := strings.Split(f.Type, ".")

		for t, fn := range map[string]Method{
			"scan":  f.Scan,
			"value": f.Value,
			"nil":   f.Nil,
		} {
			var pkg string

			if t == "scan" {
				if fn.Pkg == "" {
					continue
				}

				if fn.Pkg == "_" {
					if f.IsNullable {
						continue
					}

					parts := strings.Split(f.Type, ".")

					if len(parts) <= 1 {
						continue
					}

					pkg = parts[0]
				} else {
					pkg = fn.Pkg
				}
			} else {
				if fn.Pkg == "" || fn.Pkg == "_" {
					continue
				}

				pkg = fn.Pkg
			}

			mi, ok := mImports[pkg]
			if !ok {
				continue
			}

			pImports[fn.Pkg] = mi
		}

	}

	// TODO: sort imports
	return pImports, nil
}

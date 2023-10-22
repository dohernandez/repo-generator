package generator

import (
	"strings"

	"github.com/dohernandez/errors"
)

var nativeImports = map[string]string{
	"time": "time",
	"big":  "math/big",
}

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

func parseImports(imports []string, r Repo) (map[string]PackageImport, error) {
	pImports := make(map[string]PackageImport, len(imports)+1)

	for k, i := range defaultImports {
		pImports[k] = PackageImport{
			Path: i,
		}
	}

	for _, f := range r.Model.Fields {
		if f.HasSqlNullable {
			continue
		}

		parts := strings.Split(f.Type, ".")

		if len(parts) <= 1 {
			continue
		}

		pkg := parts[0]

		path, ok := nativeImports[pkg]
		if !ok {
			continue
		}

		if !f.HasScanMethod {
			continue
		}

		pImports[pkg] = PackageImport{
			Path: path,
		}
	}

	for _, i := range imports {
		parts := strings.Split(i, ":")

		path := parts[len(parts)-1]

		pkSplit := strings.Split(path, "/")
		pk := pkSplit[len(pkSplit)-1]

		var alias string

		if len(parts) > 1 {
			alias = parts[0]
		}

		pImports[pk] = PackageImport{
			Alias: alias,
			Path:  path,
		}
	}

	// TODO: sort imports
	return pImports, nil
}

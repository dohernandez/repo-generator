package generator

import "strings"

var defaultImports = map[string]string{
	"errors": "github.com/dohernandez/repo-generator/errors",
}

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

	//for _, f := range r.Model.Fields {
	//	if f.Type == "time.Time" {
	//		pImports["time"] = PackageImport{
	//			Path: "time",
	//		}
	//
	//		break
	//	}
	//}

	// TODO: sort imports
	return pImports, nil
}

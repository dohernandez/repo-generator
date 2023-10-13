package generator

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"text/template"

	"github.com/iancoleman/strcase"

	"github.com/dohernandez/repo-generator/errors"
)

const repoTplFilename = "repo.tmpl"

//go:embed repo.tmpl
var repoTpl embed.FS

type Generator struct {
	Package string
	Imports map[string]PackageImport
	Repo    Repo
}

func Generate(sourcePath, outputPath string, model string, opts ...Option) error {
	var options Options

	for _, opt := range opts {
		opt(&options)
	}

	fset := token.NewFileSet()

	tree, err := parser.ParseFile(fset, sourcePath, nil, 0)
	if err != nil {
		return errors.Wrap(err, "parse source file")
	}

	r, err := ParseRepo(tree, model)
	if err != nil {
		return errors.Wrap(err, "parse repo")
	}

	imports, err := parseImports(options.imports, r)
	if err != nil {
		return errors.Wrap(err, "parse package imports")
	}

	g := Generator{
		Package: tree.Name.Name,
		Imports: imports,
		Repo:    r,
	}

	// Populate the functions which should be exposed to the template.
	funcMap := template.FuncMap{
		"toLowerCamel": strcase.ToLowerCamel,
		// The name "inc" is what the function will be called in the template text.
		"inc": func(i int) int {
			return i + 1
		},
		// The name "AND" is what the function will be called in the template text.
		"AND": func(i int) string {
			if i == len(r.ColKeysFields)-1 {
				return ""
			}

			return " AND"
		},
		"fieldToSql": fieldToSql,
		"fieldType":  fieldType,
		"scanField":  scanField,
		"sqlToField": sqlToField,
	}

	t, err := template.New(repoTplFilename).Funcs(funcMap).ParseFS(repoTpl, repoTplFilename)
	if err != nil {
		return errors.Wrap(err, "parse template")
	}

	var b bytes.Buffer

	if err = t.Execute(&b, g); err != nil {
		return errors.Wrap(err, "execute template")
	}

	// Format the source code before writing.
	output, err := format.Source(b.Bytes())
	if err != nil {
		return errors.Wrap(err, string(b.Bytes()))
	}

	if err = os.WriteFile(outputPath, output, 0o600); err != nil {
		return errors.Wrap(err, "write output file")
	}

	_, _ = fmt.Fprintf(os.Stdout, "successfully wrote %s\n", outputPath)

	return nil
}

func fieldToSql(f any) string {
	fd, ok := f.(Field)
	if !ok {
		return ""
	}

	if fd.Auto {
		fd.OmittedEmpty = true
	}

	if fd.IsArrayable {
		value := tmplFieldValueMethod(fd, fmt.Sprintf("m.%s[i]", fd.Name))

		tmpl := `
			cols = append(cols, "%[1]s")

			%[2]ss := make([]%[3]s, len(m.%[4]s))

			for i := range m.%[4]s {
				%[2]ss[i] = %[5]s
			}

			args = append(args, %[2]ss)
			`

		if fd.IsPointer && !fd.HasValueMethod {
			value = fmt.Sprintf("*%s", value)
		}

		return fmt.Sprintf(tmpl, fd.ColName, fd.LowerCaseName, fd.SqlType, fd.Name, value)

	}

	tmpl := `
			cols = append(cols, "%[1]s")
			args = append(args, %[2]s)
			`

	value := tmplFieldValueMethod(fd, fmt.Sprintf("m.%s", fd.Name))

	if fd.IsNullable || fd.OmittedEmpty {
		tnullable := fmt.Sprintf("%s.%s", fd.LowerCaseName, fd.SqlNullable.set)

		// It is only nullable.
		if !fd.OmittedEmpty {
			// Has a nullable type.
			if !fd.HasSqlNullable {
				return fmt.Sprintf(tmpl, fd.ColName, value)
			}

			tmpl = `
				var %[1]s %[2]s

				%[3]s = %[4]s
				%[1]s.Valid = true

				cols = append(cols, "%[5]s")
				args = append(args, %[1]s)
				`

			if fd.IsPointer && !fd.HasValueMethod {
				value = fmt.Sprintf("*%s", value)
			}

			return fmt.Sprintf(tmpl, fd.LowerCaseName, fd.SqlType, tnullable, value, fd.ColName)
		}

		var ifEmpty string

		if fd.IsPointer {
			ifEmpty = fmt.Sprintf("m.%s != nil", fd.Name)
		} else {
			switch fd.Type {
			case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
				ifEmpty = fmt.Sprintf("m.%s != 0", fd.Name)
			case "time.Time":
				ifEmpty = fmt.Sprintf("!m.%s.IsZero()", fd.Name)
			case "string":
				ifEmpty = fmt.Sprintf("m.%s != \"\"", fd.Name)
			}
		}

		// It is only omitted empty.
		if !fd.IsNullable {
			if ifEmpty != "" {
				tmpl = `
					if %[3]s {
						cols = append(cols, "%[1]s")
						args = append(args, %[2]s)
					}
					`

				if fd.IsPointer && !fd.HasValueMethod {
					value = fmt.Sprintf("*%s", value)
				}
			}

			return fmt.Sprintf(tmpl, fd.ColName, value, ifEmpty)
		}

		// It is nullable and omitted empty.
		if ifEmpty != "" {
			tmpl = `
				if %[1]s {
					var %[2]s %[3]s

					%[4]s = %[5]s
					%[2]s.Valid = true

					cols = append(cols, "%[6]s")
					args = append(args, %[2]s)
				}
				`

			if !fd.HasSqlNullable {
				tmpl = `
					if %[1]s {
						cols = append(cols, "%[6]s")
						args = append(args, %[2]s)
					}
					`
			}

			if fd.IsPointer && !fd.HasValueMethod {
				value = fmt.Sprintf("*%s", value)
			}

			return fmt.Sprintf(tmpl, ifEmpty, fd.LowerCaseName, fd.SqlType, tnullable, value, fd.ColName)
		}

		if fd.HasSqlNullable {
			tmpl = `
				var %[1]s %[2]s

				%[3]s = %[4]s
				%[1]s.Valid = true

				cols = append(cols, "%[5]s")
				args = append(args, %[1]s)
				`

			if fd.IsPointer && !fd.HasValueMethod {
				value = fmt.Sprintf("*%s", value)
			}

			return fmt.Sprintf(tmpl, fd.LowerCaseName, fd.SqlType, tnullable, value, fd.ColName)
		}

		if fd.IsPointer && !fd.HasValueMethod {
			value = fmt.Sprintf("*%s", value)
		}

		return fmt.Sprintf(tmpl, fd.ColName, value)
	}

	if fd.IsPointer {
		tmpl = `
			if m.%[3]s != nil {
				cols = append(cols, "%[1]s")
				args = append(args, %[2]s)
			}
			`

		if !fd.HasValueMethod {
			value = fmt.Sprintf("*m.%s", fd.Name)
		}
	}

	return fmt.Sprintf(tmpl, fd.ColName, value, fd.Name)
}

func tmplFieldValueMethod(f Field, a string) string {
	if !f.HasValueMethod {
		return a
	}

	scan := fmt.Sprintf("%s.%s()", a, f.Value.Name)

	if f.Value.Pkg == "_" {
		return scan
	}

	if f.Value.Pkg != "" {
		return fmt.Sprintf("%s.%s(%s)", f.Value.Pkg, f.Value.Name, a)
	}

	return fmt.Sprintf("%s(%s)", f.Value.Name, a)
}

func fieldType(f any) string {
	fd, ok := f.(Field)
	if !ok {
		return ""
	}

	if fd.IsArrayable {
		if fd.HasSqlArrayable {
			return fmt.Sprintf("%ss %s", fd.LowerCaseName, fd.SqlArrayable)
		}

		return fmt.Sprintf("%ss %s", fd.LowerCaseName, fd.SqlType)
	}

	if fd.IsNullable {
		return fmt.Sprintf("%s %s", fd.LowerCaseName, fd.SqlType)
	}

	if fd.IsPointer {
		if fd.SqlType != "" {
			return fmt.Sprintf("%s %s", fd.LowerCaseName, fd.SqlType)
		}

		return fmt.Sprintf("%s %s", fd.LowerCaseName, fd.Type)
	}

	if fd.SqlType != "" {
		return fmt.Sprintf("%s %s", fd.LowerCaseName, fd.SqlType)
	}

	if fd.HasScanMethod {
		return fmt.Sprintf("%s %s", fd.LowerCaseName, fd.Type)
	}

	return ""
}

func scanField(f any) string {
	fd, ok := f.(Field)
	if !ok {
		return ""
	}

	if fd.IsArrayable {
		return fmt.Sprintf("&%ss", fd.LowerCaseName)
	}

	if fd.IsNullable {
		return fmt.Sprintf("&%s", fd.LowerCaseName)
	}

	if fd.IsPointer {
		return fmt.Sprintf("&%s", fd.LowerCaseName)
	}

	if fd.SqlType != "" {
		return fmt.Sprintf("&%s", fd.LowerCaseName)
	}

	return fmt.Sprintf("&m.%s", fd.Name)
}

func sqlToField(f any) string {
	fd, ok := f.(Field)
	if !ok {
		return ""
	}

	if fd.IsArrayable {
		scan := tmplFieldSanMethod(fd, fmt.Sprintf("%ss[i]", fd.LowerCaseName))

		tmpl := `
			for i := range %ss {
				m.%s = append(m.%s, %s)
			}
			`

		return fmt.Sprintf(tmpl, fd.LowerCaseName, fd.Name, fd.Name, scan)
	}

	scan := tmplFieldSanMethod(fd, fd.LowerCaseName)

	if fd.IsNullable {
		a := fd.LowerCaseName

		if fd.HasSqlNullable {
			a = fmt.Sprintf("%s.%s", fd.LowerCaseName, fd.SqlNullable.set)
		}

		scan = tmplFieldSanMethod(fd, a)

		// Output:
		// if so.Valid {
		// 		m.SO = so.String
		// }
		tmpl := `
			if %[1]s.Valid {
				m.%[3]s = %[2]s
			}
			`

		if fd.IsPointer && !fd.HasScanMethod {
			// Output:
			// if so.Valid {
			// 		tmp = so.String
			//		m.SO = &tmp
			// }

			tmpl = `
			if %[1]s.Valid {
				tmp := %[2]s
				m.%[3]s = &tmp
			}
			`
		}

		return fmt.Sprintf(tmpl, fd.LowerCaseName, scan, fd.Name)
	}

	if fd.IsPointer {
		// Output:
		// 		tmp = so.String
		//		m.SO = &tmp
		// }

		tmpl := `
			tmp := %[2]s
			m.%[3]s = &tmp
			`

		if fd.HasScanMethod {
			tmpl = `
				m.%[3]s = %[2]s
				`
		}

		return fmt.Sprintf(tmpl, fd.LowerCaseName, scan, fd.Name)
	}

	if !fd.HasScanMethod && fd.SqlType == "" {
		return ""
	}

	return fmt.Sprintf("m.%s = %s", fd.Name, scan)
}

func tmplFieldSanMethod(f Field, a string) string {
	if !f.HasScanMethod {
		return a
	}

	scan := fmt.Sprintf("%s.%s()", a, f.Scan.Name)

	if f.Scan.Pkg == "_" {
		return scan
	}

	if f.Scan.Pkg != "" {
		return fmt.Sprintf("%s.%s(%s)", f.Scan.Pkg, f.Scan.Name, a)
	}

	return fmt.Sprintf("%s(%s)", f.Scan.Name, a)
}

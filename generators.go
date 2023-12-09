package generator

import (
	"bytes"
	"embed"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"text/template"

	"github.com/dohernandez/errors"
	"github.com/iancoleman/strcase"
	"mvdan.cc/gofumpt/format"
)

const version = "v0.1.0"

const repoTplFilename = "repo.tmpl"

//go:embed repo.tmpl
var repoTpl embed.FS

type Generator struct {
	Package string
	Imports map[string]PackageImport
	Repo    Repo
	Version string

	Funcs []repoFunc
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

	imports, err := parseImports(tree, r)
	if err != nil {
		return errors.Wrap(err, "parse package imports")
	}

	g := Generator{
		Package: tree.Name.Name,
		Imports: imports,
		Repo:    r,
		Version: version,
		Funcs:   options.funcs,
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
		"fieldToCreateSql": fieldToCreateSql(r),
		"fieldType":        fieldType(r),
		"scanField":        scanField(r),
		"sqlToField":       sqlToField(r),
		"fieldToInsertSql": fieldToInsertSql(r),
		"fieldToUpdateSql": fieldToUpdateSql(r),
		"fieldValueMethod": fieldValueMethod,
		"has": func(ls []repoFunc, c ...string) bool {
			for _, v := range ls {
				for _, s := range c {
					if v.String() == s {
						return true
					}
				}
			}

			return false
		},
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
	output, err := format.Source(b.Bytes(), format.Options{})
	if err != nil {
		return errors.Wrap(err, string(b.Bytes()))
	}

	if err = os.WriteFile(outputPath, output, 0o600); err != nil {
		return errors.Wrap(err, "write output file")
	}

	_, _ = fmt.Fprintf(os.Stdout, "successfully wrote %s\n", outputPath)

	return nil
}

func fieldToCreateSql(repo Repo) func(f any) string {
	return func(f any) string {
		fd, ok := f.(Field)
		if !ok {
			return ""
		}

		omittedEmpty := fd.OmittedEmpty

		if fd.Auto {
			omittedEmpty = true
		}

		colName := fmt.Sprintf("%s.col%s", repo.Receiver, fd.Name)

		if fd.IsArrayable {
			return fieldArrayable(fd, colName, false, func() string {
				return `
					cols = append(cols, %[1]s)
					args = append(args, %[2]ss)
				`
			})
		}

		value := fieldValueMethod(fd, fmt.Sprintf("m.%s", fd.Name))

		if fd.IsNullable || omittedEmpty {
			// It is only nullable.
			if !omittedEmpty {
				return fieldNullable(fd, colName, value, func() string {
					return `
						cols = append(cols, %[5]s)
						args = append(args, %[1]s)
					`
				})
			}

			// It is only omitted empty.
			if !fd.IsNullable {
				return fieldOmitEmpty(fd, colName, value, true, func() string {
					return `
						cols = append(cols, %[1]s)
						args = append(args, %[2]s)
					`
				})
			}

			// It is nullable and omitted empty.
			return fieldOmitEmptyNullable(fd, colName, value)
		}

		if fd.IsPointer {
			tmpl := `
			if m.%[3]s != nil {
				cols = append(cols, %[1]s)
				args = append(args, %[2]s)
			}`

			if !fd.HasValueMethod {
				value = fmt.Sprintf("*m.%s", fd.Name)
			}

			return fmt.Sprintf(tmpl, colName, value, fd.Name)
		}

		tmpl := `
			cols = append(cols, %[1]s)
			args = append(args, %[2]s)
		`

		return fmt.Sprintf(tmpl, colName, value)
	}
}

func fieldArrayable(f Field, colName string, skipZeroValues bool, tmplFunc func() string) string {
	value := fieldValueMethod(f, fmt.Sprintf("m.%s[i]", f.Name))

	if f.IsPointer && !f.HasValueMethod {
		value = fmt.Sprintf("*%s", value)
	}

	if colName != "" && (skipZeroValues || f.OmittedEmpty || f.Auto) {
		tmpl := `
				if len(m.%[4]s) > 0 {
					%[2]ss := make([]%[3]s, len(m.%[4]s))

					for i := range m.%[4]s {
						%[2]ss[i] = %[5]s
					}` +
			tmplFunc() +
			`}`

		return fmt.Sprintf(tmpl, colName, f.LowerCaseName, f.SqlType, f.Name, value)
	}

	tmpl := `
			%[2]ss := make([]%[3]s, len(m.%[4]s))

			for i := range m.%[4]s {
				%[2]ss[i] = %[5]s
			}
			` +
		tmplFunc()

	return fmt.Sprintf(tmpl, colName, f.LowerCaseName, f.SqlType, f.Name, value)
}

func fieldValueMethod(f Field, a string) string {
	if !f.HasValueMethod {
		return a
	}

	if f.Value.Pkg == "_" {
		return fmt.Sprintf("%s.%s()", a, f.Value.Name)
	}

	if f.Value.Pkg != "" {
		return fmt.Sprintf("%s.%s(%s)", f.Value.Pkg, f.Value.Name, a)
	}

	return fmt.Sprintf("%s(%s)", f.Value.Name, a)
}

func fieldNullable(f Field, colName, value string, tmplFunc func() string) string {
	if f.IsPointer && !f.HasValueMethod {
		value = fmt.Sprintf("*%s", value)
	}

	// Has a nullable type.
	if !f.HasSqlNullable {
		return fmt.Sprintf(tmplFunc(), colName, value)
	}

	tnullable := fmt.Sprintf("%s.%s", f.LowerCaseName, f.SqlNullable.set)

	tmpl := `
		var %[1]s %[2]s

		%[3]s = %[4]s
		%[1]s.Valid = true

		` + tmplFunc()

	return fmt.Sprintf(tmpl, f.LowerCaseName, f.SqlType, tnullable, value, colName)
}

func fieldOmitEmpty(f Field, colName, value string, skipZeroValues bool, tmplFunc func() string) string {
	if f.IsPointer && !f.HasValueMethod {
		value = fmt.Sprintf("*%s", value)
	}

	if !skipZeroValues {
		tmpl := tmplFunc()

		return fmt.Sprintf(tmpl, colName, value)
	}

	if f.IsKey && !f.Auto {
		tmpl := tmplFunc()

		return fmt.Sprintf(tmpl, colName, value)
	}

	ifEmpty := ifEmptyStatement(f, value)

	if ifEmpty == "" {
		tmpl := tmplFunc()

		return fmt.Sprintf(tmpl, colName, value)
	}

	tmpl := `
		if %[3]s {` +
		tmplFunc() +
		`}`

	return fmt.Sprintf(tmpl, colName, value, ifEmpty)
}

func ifEmptyStatement(f Field, value string) string {
	if f.IsPointer {
		return fmt.Sprintf("m.%s != nil", f.Name)
	}

	if !f.HasNilMethod {
		return ""
	}

	tmpl := tmplFieldNilMethod(f, fmt.Sprintf("m.%s", f.Name))

	if f.Nil.Name == "" {
		tmpl = value
	}

	switch f.Nil.CmpOperator {
	case MethodCmpOperatorNotEqual, MethodCmpOperatorGreater:
		return fmt.Sprintf(
			"%s %s %s",
			tmpl,
			f.Nil.CmpOperator.String(),
			f.Nil.EmptyValue,
		)
	case MethodCmpOperatorNot:
		return f.Nil.CmpOperator.String() + tmpl
	}

	return "operator not implemented"
}

func fieldOmitEmptyNullable(f Field, colName, value string) string {
	if f.IsPointer && !f.HasValueMethod {
		value = fmt.Sprintf("*%s", value)
	}

	tnullable := fmt.Sprintf("%s.%s", f.LowerCaseName, f.SqlNullable.set)
	ifEmpty := ifEmptyStatement(f, value)

	if ifEmpty == "" {
		if f.HasSqlNullable {
			tmpl := `
				var %[1]s %[2]s

				%[3]s = %[4]s
				%[1]s.Valid = true

				cols = append(cols, %[5]s)
				args = append(args, %[1]s)
				`

			return fmt.Sprintf(tmpl, f.LowerCaseName, f.SqlType, tnullable, value, colName)
		}

		tmpl := `
			cols = append(cols, %[1]s)
			args = append(args, %[2]s)
			`

		return fmt.Sprintf(tmpl, colName, value)
	}

	if !f.HasSqlNullable {
		tmpl := `
			if %[1]s {
				cols = append(cols, %[2]s)
				args = append(args, %[3]s)
			}
			`

		return fmt.Sprintf(tmpl, ifEmpty, colName, value)
	}

	tmpl := `
		if %[1]s {
			var %[2]s %[3]s

			%[4]s = %[5]s
			%[2]s.Valid = true

			cols = append(cols, %[6]s)
			args = append(args, %[2]s)
		}
			`
	return fmt.Sprintf(tmpl, ifEmpty, f.LowerCaseName, f.SqlType, tnullable, value, colName)
}

func fieldType(_ Repo) func(f any) string {
	return func(f any) string {
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
}

func scanField(_ Repo) func(f any) string {
	return func(f any) string {
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
}

func sqlToField(_ Repo) func(f any) string {
	return func(f any) string {
		fd, ok := f.(Field)
		if !ok {
			return ""
		}

		fieldCaseName := fd.LowerCaseName

		if fd.IsArrayable {
			scan := fieldSanMethod(fd, fmt.Sprintf("%ss[i]", fd.LowerCaseName))

			tmpl := `
			for i := range %ss {
				m.%s = append(m.%s, %s)
			}
			`

			return fmt.Sprintf(tmpl, fieldCaseName, fd.Name, fd.Name, scan)
		}

		scan := fieldSanMethod(fd, fd.LowerCaseName)

		if fd.IsNullable {
			fieldCaseName := fd.LowerCaseName

			if fd.HasSqlNullable {
				fieldCaseName = fmt.Sprintf("%s.%s", fd.LowerCaseName, fd.SqlNullable.set)
			}

			scan = fieldSanMethod(fd, fieldCaseName)

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
}

func fieldSanMethod(f Field, a string) string {
	if !f.HasScanMethod {
		return a
	}

	if f.Scan.Pkg == "_" {
		return fmt.Sprintf("%s.%s()", a, f.Scan.Name)
	}

	if f.Scan.Pkg != "" {
		return fmt.Sprintf("%s.%s(%s)", f.Scan.Pkg, f.Scan.Name, a)
	}

	return fmt.Sprintf("%s(%s)", f.Scan.Name, a)
}

func fieldToInsertSql(_ Repo) func(f any) string {
	return func(f any) string {
		fd, ok := f.(Field)
		if !ok {
			return ""
		}

		if fd.Auto {
			fd.OmittedEmpty = true
		}

		if fd.IsArrayable {
			return fieldArrayable(fd, "", false, func() string {
				return `
					args = append(args, %[2]ss)
				`
			})
		}

		value := fieldValueMethod(fd, fmt.Sprintf("m.%s", fd.Name))

		if fd.IsNullable {
			return fieldNullable(fd, "", value, func() string {
				return `
					args = append(args, %[1]s)
				`
			})
		}

		if fd.IsPointer {
			tmpl := `
			if m.%[2]s != nil {
				args = append(args, %[1]s)
			}`

			if !fd.HasValueMethod {
				value = fmt.Sprintf("*m.%s", fd.Name)
			}

			return fmt.Sprintf(tmpl, value, fd.Name)
		}

		tmpl := `
			args = append(args, %[1]s)
		`

		return fmt.Sprintf(tmpl, value)
	}
}

func tmplFieldNilMethod(f Field, a string) string {
	if !f.HasNilMethod || f.Nil.Name == "" {
		return a
	}

	if f.Nil.Pkg == "_" {
		return fmt.Sprintf("%s.%s()", a, f.Nil.Name)
	}

	if f.Nil.Pkg != "" {
		return fmt.Sprintf("%s.%s(%s)", f.Nil.Pkg, f.Nil.Name, a)
	}

	return fmt.Sprintf("%s(%s)", f.Nil.Name, a)
}

func fieldToUpdateSql(repo Repo) func(f Field, b bool) string {
	return func(f Field, b bool) string {
		colName := fmt.Sprintf("%s.col%s", repo.Receiver, f.Name)

		if f.IsArrayable {
			return fieldArrayable(f, colName, b, func() string {
				return `
					sets = append(sets, fmt.Sprintf("%%s = $%%d", %[1]s, offset))
					args = append(args, %[2]ss)
			
					offset++
				`
			})
		}

		value := fieldValueMethod(f, fmt.Sprintf("m.%s", f.Name))

		return fieldOmitEmpty(f, colName, value, b, func() string {
			return `
				sets = append(sets, fmt.Sprintf("%%s = $%%d", %[1]s, offset))
				args = append(args, %[2]s)
		
				offset++
			`
		})
	}
}

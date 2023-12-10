package generator

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"

	"github.com/dohernandez/errors"
	"github.com/iancoleman/strcase"
)

const (
	columnTag      = "db"
	columnTypeTag  = "type"
	columnScanTag  = "scan"
	columnValueTag = "value"
	columnNilTag   = "nil"

	columnNullableOpt  = "nullable"
	columnIsKeyOpt     = "key"
	columnOmitEmptyOpt = "omitempty"
	columnArrayableOpt = "arrayable"
	columnAutoOpt      = "auto"
)

type MethodCmpOperator string

func (o MethodCmpOperator) String() string {
	return string(o)
}

const (
	MethodCmpOperatorEqual    MethodCmpOperator = "=="
	MethodCmpOperatorNotEqual MethodCmpOperator = "!="
	MethodCmpOperatorNot      MethodCmpOperator = "!"
	MethodCmpOperatorGreater  MethodCmpOperator = ">"
)

type Method struct {
	Name                string
	Pkg                 string
	NotEqualCmpOperator MethodCmpOperator
	EqualCmpOperator    MethodCmpOperator
	EmptyValue          string
}

// Field describes a field within a builder struct.
type Field struct {
	// Name of the field as in the struct.
	Name          string
	LowerCaseName string

	// Type of the field. Represents the type of the field in the struct.
	Type string

	IsPointer bool

	IsArray bool

	// ColName is the name of the column in the database.
	ColName string

	// SqlType is the type of the field use to scan/value from the database.
	SqlType string

	SqlNullable sqltype

	HasSqlNullable bool

	SqlArrayable sqlarrayable

	HasSqlArrayable bool

	// IsNullable indicates if the column is nullable.
	IsNullable bool

	// IsKey indicates if the field is a column key.
	IsKey bool

	// Auto indicates if the field is a column key auto-generated by the database.
	Auto bool

	// OmittedEmpty indicates if the field is omitted if empty.
	OmittedEmpty bool

	// IsArrayable indicates if the field is a column arrayable.
	IsArrayable bool

	HasScanMethod bool

	// Scan is the method used to scan the column value.
	Scan Method

	HasValueMethod bool

	// Value is the method used to get the column value.
	Value Method

	HasNilMethod bool

	// Value is the method used to get the column value.
	Nil Method
}

type Model struct {
	Name string

	Auto bool

	Receiver string

	Fields []Field
}

func ParseModel(s *ast.StructType, name string) (Model, error) {
	fields, err := parseFields(s)
	if err != nil {
		return Model{}, err
	}

	var auto bool

	for _, f := range fields {
		if f.Auto {
			auto = true

			break
		}
	}

	return Model{
		Name:     name,
		Receiver: Receiver(name),
		Auto:     auto,
		Fields:   fields,
	}, nil
}

func parseFields(s *ast.StructType) ([]Field, error) {
	var (
		fields = make([]Field, 0)
	)

	for _, field := range s.Fields.List {
		// If a field has no names, then it's an anonymous / embedded field.
		if field.Names == nil {
			switch t := field.Type.(type) {
			case *ast.Ident:
				// This is specifically for the "Model" field. If panic here, it means that the Model field has
				// been declared differently.
				ifields, err := parseFields(t.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType))
				if err != nil {
					return nil, err
				}

				fields = append(fields, ifields...)

				continue
			case *ast.SelectorExpr:
				continue
			}
		}

		if field.Tag == nil {
			continue
		}

		// If the key is empty or "-" then we should skip this field.
		st := reflect.StructTag(strings.Trim(field.Tag.Value, "`")).Get(columnTag)
		if st == "" || st == "-" {
			continue
		}

		f, err := parseField(field)
		if err != nil {
			return nil, err
		}

		fields = append(fields, f)
	}

	return fields, nil
}

type tagColumn struct {
	name        string
	isKey       bool
	Auto        bool
	isNullable  bool
	oEmpty      bool
	isArrayable bool
}

func parseField(field *ast.Field) (Field, error) {
	var (
		isPointer bool
		isArray   bool

		ftype      string
		sType      string
		sNullable  sqltype
		sArrayable sqlarrayable

		hasValueMethod bool
		hasArrayable   bool
		hasSqlNullable bool
		hasScanMethod  bool
		hasNilMethod   bool
	)

	switch t := field.Type.(type) {
	case *ast.Ident:
		ftype = t.Name
	case *ast.SelectorExpr:
		ftype = fmt.Sprintf("%s.%s", t.X, t.Sel.Name)
	case *ast.StarExpr:
		isPointer = true

		switch tx := t.X.(type) {
		case *ast.SelectorExpr:
			ftype = fmt.Sprintf("%s.%s", tx.X, tx.Sel.Name)
		case *ast.Ident:
			ftype = tx.Name
		default:
			return Field{}, errors.Newf("unsupported type %T.Type %T", t, t.X)
		}
	case *ast.ArrayType:
		isArray = true

		switch tx := t.Elt.(type) {
		case *ast.SelectorExpr:
			ftype = fmt.Sprintf("%s.%s", tx.X, tx.Sel.Name)
		case *ast.Ident:
			ftype = tx.Name
		default:
			return Field{}, errors.Newf("unsupported type %T.Type %T", t, t.Elt)
		}
	}

	tagCol := parseTagColumn(field)
	sType = parseTagType(field)

	t := sType

	if t == "" {
		t = ftype
	}

	if tagCol.isNullable {
		ss, ok := sqlnullable[t]
		if ok {
			sType = ss.t
			sNullable = ss
			hasSqlNullable = true
		}
	}

	if tagCol.isArrayable {
		if sa, ok := sqlArrayable[t]; ok {
			sArrayable = sa
			hasArrayable = true
		}
	}

	sMethod := parseTagScan(field)
	if sMethod != (Method{}) {
		hasScanMethod = true
	}

	vMethod := parseTagValue(field)
	if vMethod != (Method{}) {
		hasValueMethod = true
	}

	nMethod := parseTagNil(field, ftype, sType, isArray, hasValueMethod)
	if nMethod != (Method{}) {
		hasNilMethod = true
	}

	name := field.Names[0].Name

	return Field{
		Name:          name,
		LowerCaseName: lowerCaseName(name),
		Type:          ftype,

		IsPointer: isPointer,
		IsArray:   isArray,

		ColName: tagCol.name,

		SqlType: sType,

		SqlNullable:    sNullable,
		HasSqlNullable: hasSqlNullable,

		SqlArrayable:    sArrayable,
		HasSqlArrayable: hasArrayable,

		IsNullable:   tagCol.isNullable,
		IsKey:        tagCol.isKey,
		Auto:         tagCol.Auto,
		OmittedEmpty: tagCol.oEmpty,
		IsArrayable:  tagCol.isArrayable,

		Scan:           sMethod,
		HasScanMethod:  hasScanMethod,
		Value:          vMethod,
		HasValueMethod: hasValueMethod,
		Nil:            nMethod,
		HasNilMethod:   hasNilMethod,
	}, nil
}

func parseTagColumn(field *ast.Field) tagColumn {
	tagCol := tagColumn{
		name: strcase.ToSnake(field.Names[0].Name),
	}

	st := reflect.StructTag(strings.Trim(field.Tag.Value, "`")).Get(columnTag)

	// Split the options from the name
	parts := strings.Split(st, ",")

	tagCol.name = parts[0]

	opts := parts[1:]

	if len(opts) == 0 {
		return tagCol
	}

	for _, opt := range opts {
		if strings.Contains(opt, "=") {
			//kv := strings.Split(opt, "=")
			continue
		}

		if opt == columnNullableOpt {
			tagCol.isNullable = true
		}

		if opt == columnIsKeyOpt {
			tagCol.isKey = true
		}

		if opt == columnAutoOpt {
			tagCol.Auto = true
		}

		if opt == columnOmitEmptyOpt {
			tagCol.oEmpty = true
		}

		if opt == columnArrayableOpt {
			tagCol.isArrayable = true
		}
	}

	return tagCol
}

func parseTagType(field *ast.Field) string {
	st := reflect.StructTag(strings.Trim(field.Tag.Value, "`")).Get(columnTypeTag)

	// Split the options from the name
	parts := strings.Split(st, ",")

	if len(parts) == 0 {
		return ""
	}

	return parts[0]
}

func parseTagScan(field *ast.Field) Method {
	return parseTagMethod(field, columnScanTag)
}

func parseTagValue(field *ast.Field) Method {
	return parseTagMethod(field, columnValueTag)
}

func parseTagNil(field *ast.Field, fType, sType string, isArray, hasValueMethod bool) Method {
	nMethod := parseTagMethod(field, columnNilTag)

	if nMethod != (Method{}) {
		nMethod.NotEqualCmpOperator = MethodCmpOperatorNot

		return nMethod
	}

	if isArray {
		return Method{
			Name:                "len",
			NotEqualCmpOperator: MethodCmpOperatorGreater,
			EqualCmpOperator:    MethodCmpOperatorEqual,
			EmptyValue:          "0",
		}
	}

	if hasValueMethod {
		if sType != "" {
			fType = sType
		}
	}

	switch fType {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
		nMethod = Method{
			NotEqualCmpOperator: MethodCmpOperatorNotEqual,
			EqualCmpOperator:    MethodCmpOperatorEqual,
			EmptyValue:          "0",
		}
	case "time.Time":
		nMethod = Method{
			Name:                "IsZero",
			Pkg:                 "_",
			NotEqualCmpOperator: MethodCmpOperatorNot,
		}
	case "string":
		nMethod = Method{
			NotEqualCmpOperator: MethodCmpOperatorNotEqual,
			EqualCmpOperator:    MethodCmpOperatorEqual,
			EmptyValue:          "\"\"",
		}
	}

	return nMethod
}

func parseTagMethod(field *ast.Field, tag string) Method {
	st := reflect.StructTag(strings.Trim(field.Tag.Value, "`")).Get(tag)

	// Split the options from the name
	parts := strings.Split(st, ".")

	if len(parts) == 0 {
		return Method{}
	}

	if len(parts) == 2 {
		if parts[0] == "" && strings.HasPrefix(st, ".") {
			return Method{
				Pkg:  "_",
				Name: parts[1],
			}
		}

		return Method{
			Pkg:  parts[0],
			Name: parts[1],
		}
	}

	return Method{
		Name: parts[0],
	}
}

func lowerCaseName(name string) string {
	n := strcase.ToLowerCamel(name)

	if n == "type" {
		return "typ"
	}

	return n
}

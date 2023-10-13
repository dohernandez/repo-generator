package generator

type sqltype struct {
	t    string
	set  string
	ifst string
}

var sqlnullable = map[string]sqltype{
	"bool": {
		t:    "sql.NullBool",
		set:  "Bool",
		ifst: "true",
	},
	"string": {
		t:    "sql.NullString",
		set:  "String",
		ifst: "\"\"",
	},
	"float64": {
		t:    "sql.NullFloat64",
		set:  "Float64",
		ifst: "0",
	},
	"int16": {
		t:    "sql.NullInt16",
		set:  "Int16",
		ifst: "0",
	},
	"int32": {
		t:    "sql.NullInt32",
		set:  "Int32",
		ifst: "0",
	},
	"int64": {
		t:    "sql.NullInt64",
		set:  "Int64",
		ifst: "0",
	},
	"time.Time": {
		t:    "sql.NullTime",
		set:  "Time",
		ifst: ".IsZero",
	},
}

type sqlarrayable string

var sqlArrayable = map[string]sqlarrayable{
	"string": "pq.StringArray",
	"bool":   "pq.BoolArray",
}

//var sqltypeSwapper = map[string]string{
//	"big.Int": "int64",
//}

//var sqltypeImporter = map[string]map[string]bool{
//	"time.Time": {"time": true},
//	"big.Int":   {"math/big": true},
//}

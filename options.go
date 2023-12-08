package generator

import "github.com/dohernandez/errors"

type Option func(*Options)

func WithImports(imports []string) Option {
	return func(opt *Options) {
		opt.imports = imports
	}
}

type repoFunc string

func (f repoFunc) String() string {
	return string(f)
}

const (
	scanFunc    repoFunc = "scan"
	scanAllFunc repoFunc = "scanAll"
	createFunc  repoFunc = "create"
	insertFunc  repoFunc = "insert"
	updateFunc  repoFunc = "update"
	deleteFunc  repoFunc = "delete"
)

func repoFuncFromString(s string) (repoFunc, error) {
	switch s {
	case scanFunc.String():
		return scanFunc, nil
	case scanAllFunc.String():
		return scanAllFunc, nil
	case createFunc.String():
		return createFunc, nil
	case insertFunc.String():
		return insertFunc, nil
	case deleteFunc.String():
		return deleteFunc, nil
	default:
		return "", errors.New("")
	}
}

func WithCreateFunc() Option {
	return withFunc(createFunc)
}

func withFunc(funcs ...repoFunc) Option {
	return func(opt *Options) {
		opt.funcs = append(opt.funcs, funcs...)
	}
}

func WithInsertFunc() Option {
	return withFunc(insertFunc)
}

func WithUpdateFunc() Option {
	return withFunc(updateFunc)
}

func WithDeleteFunc() Option {
	return withFunc(deleteFunc)
}

type Options struct {
	imports []string
	funcs   []repoFunc
}

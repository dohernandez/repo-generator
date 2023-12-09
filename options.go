package generator

type Option func(*Options)

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
	funcs []repoFunc
}

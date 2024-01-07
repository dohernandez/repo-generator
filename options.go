package generator

// Option sets up a Generate options.
type Option func(*Options)

type repoFunc string

// String returns the string representation of the repoFunc.
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

// WithCreateFunc sets the create function.
func WithCreateFunc() Option {
	return withFunc(createFunc)
}

func withFunc(funcs ...repoFunc) Option {
	return func(opt *Options) {
		opt.funcs = append(opt.funcs, funcs...)
	}
}

// WithInsertFunc sets the insert function.
func WithInsertFunc() Option {
	return withFunc(insertFunc)
}

// WithUpdateFunc sets the update function.
func WithUpdateFunc() Option {
	return withFunc(updateFunc)
}

// WithDeleteFunc sets the delete function.
func WithDeleteFunc() Option {
	return withFunc(deleteFunc)
}

// Options holds the Generate options.
type Options struct {
	funcs []repoFunc
}

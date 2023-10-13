package generator

type Option func(interface{})

func WithTag(tag string) Option {
	return func(opt interface{}) {
		opt.(*Options).tag = tag
	}
}

func WithImports(imports []string) Option {
	return func(opt interface{}) {
		opt.(*Options).imports = imports
	}
}

type Options struct {
	tag     string
	imports []string
}

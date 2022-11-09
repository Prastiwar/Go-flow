package config

import (
	"flag"
	"os"
)

type flagProvider struct {
	set *flag.FlagSet
}

// NewFlagProvider returns a new flag provider with defined flags which
// will be used for command line arguments parsing. It panics if any value field of flag
// does not implement flag.Getter interface.
func NewFlagProvider(flags ...flag.Flag) *flagProvider {
	set := flag.NewFlagSet("configs.FlagProvider", flag.ContinueOnError)
	for _, f := range flags {
		if _, ok := f.Value.(flag.Getter); !ok {
			panic(wrapErrMustImplementGetter(f))
		}
		set.Var(f.Value, f.Name, f.Usage)
	}

	return &flagProvider{
		set: set,
	}
}

// Load parses flag definitions from the argument list, which should not include the command name.
// Parsed flag value results are stored in matching v fields. If there is no matching field it
// will be ignored and it's value will not be overridden.
func (p *flagProvider) Load(v any, opts ...LoadOption) error {
	err := p.set.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	options := NewLoadOptions(opts...)
	setter := NewFieldSetter(FlagProviderName, *options)

	return setter.SetFields(v, func(key string) (any, error) {
		var f *flag.Flag
		p.set.Visit(func(visitFlag *flag.Flag) {
			if visitFlag.Name == key {
				f = visitFlag
			}
		})

		if f == nil {
			return nil, nil
		}

		g, ok := f.Value.(flag.Getter)
		if !ok {
			return nil, wrapErrMustImplementGetter(*f)
		}

		return g.Get(), nil
	})
}

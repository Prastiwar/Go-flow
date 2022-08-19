package config

import (
	"flag"
	"os"
)

type flagProvider struct {
	set *flag.FlagSet
}

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

func (p *flagProvider) Load(v any, opts ...LoadOption) error {
	options := NewLoadOptions(opts...)

	err := p.set.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	return setFields(v, *options, func(key string) (any, error) {
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

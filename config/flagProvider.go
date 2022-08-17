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
		set.Var(f.Value, f.Name, f.Usage)
	}

	return &flagProvider{
		set: set,
	}
}

func (p *flagProvider) Load(v any) error {
	return p.set.Parse(os.Args[1:])
}

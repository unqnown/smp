package ctl

import (
	"github.com/unqnown/smp/config"
	"github.com/unqnown/smp/pkg/check"
	"github.com/urfave/cli"
)

type ActionFunc = func(*cli.Context)

type CommandFunc func(conf config.Config, ctx *cli.Context) error

func NewAction(cmd CommandFunc) ActionFunc {
	return func(c *cli.Context) {
		conf, err := config.Open(c.GlobalString("config"))
		check.Fatalf(err, "open config: %v", err)

		err = conf.SetNamespace(c.GlobalString("namespace"))
		check.Fatal(err)

		check.Fatal(cmd(conf, c))
	}
}

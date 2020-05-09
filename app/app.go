package app

import (
	"encoding/hex"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/unqnown/smp/config"
	"github.com/unqnown/smp/pkg/check"
	"github.com/unqnown/smp/pkg/ctl"
	"github.com/unqnown/swca1"
	"github.com/urfave/cli"
)

func Run() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	app := cli.NewApp()

	app.Name = "smp"
	app.Version = "v0.1.0"
	app.Usage = "Say my password."
	app.Description = "To start using smp immediately run `smp init`."
	app.UseShortOptionHandling = true
	app.Action = ctl.NewAction(smp)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, conf",
			Usage:  "Path to config `FILE`",
			Value:  filepath.Join(home, ".smp/config.yaml"),
			EnvVar: "SMPCONFIG",
		},
		cli.StringFlag{
			Name:  "namespace, ns, n",
			Usage: "Namespace",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:                   "init",
			Usage:                  "Initialize or reinitialize smp",
			Description:            "Creates default config file",
			Action:                 _init,
			UseShortOptionHandling: true,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "Force smp reinitialization",
				},
				cli.StringFlag{
					Name:  "secret, s",
					Usage: "Sign secret",
				},
			},
		},
		{
			Name:                   "use",
			Usage:                  "Use specified namespace",
			Action:                 ctl.NewAction(use),
			UseShortOptionHandling: true,
			Flags:                  []cli.Flag{},
		},
		{
			Name:                   "namespace",
			Aliases:                []string{"ns", "n"},
			Usage:                  "Show current namespace",
			Action:                 ctl.NewAction(namespace),
			UseShortOptionHandling: true,
			Flags:                  []cli.Flag{},
		},
		{
			Name:                   "quiet",
			Aliases:                []string{"q"},
			Usage:                  "Quietly generates password without any configuration",
			Action:                 quiet,
			UseShortOptionHandling: true,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "u",
					Usage: "Unique",
				},
				cli.BoolFlag{
					Name:  "t",
					Usage: "No type repetition",
				},
				cli.BoolFlag{
					Name:  "c",
					Usage: "No characters repetition",
				},
				cli.StringFlag{
					Name:  "abc, a",
					Usage: "Alphabet",
				},
				cli.StringFlag{
					Name:  "secret, s",
					Usage: "Secret",
				},
				cli.IntFlag{
					Name:  "size",
					Usage: "Size",
				},
			},
		},
	}

	return app.Run(os.Args)
}

func use(conf config.Config, ctx *cli.Context) error {
	if !ctx.Args().Present() {
		log.Fatalf("namespace is not specified")
	}
	err := conf.SetNamespace(ctx.Args().First())
	check.Fatal(err)
	err = conf.Save(ctx.GlobalString("config"))
	check.Fatal(err)
	log.Printf("switched to namespace %q", conf.Namespace)
	return nil
}

func namespace(conf config.Config, _ *cli.Context) error {
	log.Printf("%q", conf.Namespace)
	return nil
}

func smp(conf config.Config, ctx *cli.Context) error {
	if !ctx.Args().Present() {
		log.Fatal("no hint")
	}
	ns, err := conf.Ns()
	check.Fatal(err)
	opts := []swca1.Option{
		swca1.Size(ns.Size),
		swca1.Alphabet(ns.Alphabet),
	}
	for _, c := range ns.Complexity {
		switch c {
		case 'u':
			opts = append(opts, swca1.Unique)
		case 't':
			opts = append(opts, swca1.NoTypeRepetition)
		case 'c':
			opts = append(opts, swca1.NoCharacterRepetition)
		}
	}
	h := swca1.New(opts...)
	_, err = h.Write([]byte(ns.Secret))
	check.Fatal(err)
	_, err = h.Write([]byte(ctx.Args().First()))
	check.Fatal(err)
	log.Printf("%s", h.Sum(nil))
	return nil
}

func quiet(ctx *cli.Context) {
	if !ctx.Args().Present() {
		log.Fatal("no hint")
	}
	var opts []swca1.Option
	if abc := ctx.String("abc"); abc != "" {
		opts = append(opts, swca1.Alphabet(abc))
	}
	if size := ctx.Int("size"); size > 0 {
		opts = append(opts, swca1.Size(size))
	}
	if ctx.Bool("u") {
		opts = append(opts, swca1.Unique)
	}
	if ctx.Bool("t") {
		opts = append(opts, swca1.NoTypeRepetition)
	}
	if ctx.Bool("c") {
		opts = append(opts, swca1.NoCharacterRepetition)
	}
	var err error
	h := swca1.New(opts...)
	_, err = h.Write([]byte(ctx.String("secret")))
	check.Fatal(err)
	_, err = h.Write([]byte(ctx.Args().First()))
	check.Fatal(err)
	log.Printf("%s", h.Sum(nil))
}

func _init(ctx *cli.Context) {
	path := ctx.GlobalString("config")

	dir, _ := filepath.Split(path)

	switch conf, err := config.Open(path); {
	case err == nil:
		if !ctx.Bool("force") {
			log.Fatalf("config already exists")
		}
		// save backup in case of force re-init
		backup := filepath.Join(dir, ".backup."+time.Now().Format("02.01.06-15:04:05")+".yaml")
		err := conf.Save(backup)
		check.Fatalf(err, "unable to backup config: %v", err)
		log.Printf("previous config is backuped in %s", backup)
	case os.IsNotExist(err):
		// go ahead
	default:
		log.Fatalf("check config exists")
	}

	err := os.MkdirAll(dir, os.ModePerm)
	check.Fatal(err)

	sec := ctx.String("secret")
	if sec == "" {
		rand.Seed(time.Now().UnixNano())
		b := make([]byte, 64)
		_, _ = rand.Read(b)
		sec = hex.EncodeToString(b)
	}

	conf := config.Config{
		Namespace: "default",
		Namespaces: map[string]config.Namespace{
			"default": {
				Secret:     sec,
				Alphabet:   swca1.NULS,
				Size:       swca1.Enough,
				Complexity: "utc",
			},
		},
	}

	err = conf.Save(path)
	check.Fatalf(err, "init config: %v", err)

	log.Printf("%s created", path)
}

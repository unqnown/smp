package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Secret = string

type Config struct {
	Namespace  string               `yaml:"namespace"`
	Namespaces map[string]Namespace `yaml:"namespaces"`
}

func (conf *Config) SetNamespace(n string) error {
	if n == "" {
		return nil
	}
	if _, found := conf.Namespaces[n]; !found {
		return fmt.Errorf("namespace %q not found", n)
	}
	conf.Namespace = n
	return nil
}

type Namespace struct {
	Secret     Secret `yaml:"secret"`
	Alphabet   string `yaml:"alphabet"`
	Size       int    `yaml:"size"`
	Complexity string `yaml:"complexity"`
}

func (conf Config) Secret() (Secret, error) {
	ns, err := conf.Ns()
	return ns.Secret, err
}

func (conf Config) Ns() (Namespace, error) {
	ns, found := conf.Namespaces[conf.Namespace]
	if !found {
		return Namespace{}, fmt.Errorf("namespace %q not found", conf.Namespace)
	}
	return ns, nil
}

func (conf Config) Validate() error {
	if _, err := conf.Ns(); err != nil {
		return err
	}
	return nil
}

func Open(path string) (conf Config, err error) { return conf, conf.Open(path) }

func (conf *Config) Open(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(conf); err != nil {
		return err
	}
	return conf.Validate()
}

func (conf Config) Save(path string) error {
	if err := conf.Validate(); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewEncoder(f).Encode(conf)
}

package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Bosh    *Bosh
	Credhub *Credhub
}

type Bosh struct {
	Environment  string `required:"true"`
	Client       string `required:"true"`
	ClientSecret string `required:"true" split_words:"true"`
	CaCert       string `required:"true" split_words:"true"`
}

type Credhub struct {
	Server string `required:"true"`
	Client string `required:"true"`
	Secret string `required:"true"`
	CaCert string `required:"true" split_words:"true"`
}

func LoadConfig() (*Config, error) {
	var b Bosh
	err := envconfig.Process("bosh", &b)
	if err != nil {
		return nil, err
	}

	var c Credhub
	err = envconfig.Process("credhub", &c)
	if err != nil {
		return nil, err
	}

	return &Config{&b, &c}, nil
}

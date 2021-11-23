package config

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Host           string `envconfig:"HOST" default:"127.0.0.1"`
	Port           string `envconfig:"PORT" default:"6666"`
	Environment    string `envconfig:"ENV" default:"development"`
	JaegerEndpoint string `envconfig:"JAEGER_ENDPOINT"`
}

var config *Config
var once sync.Once

// Load reads config file and ENV variables if set.
func Load() *Config {
	once.Do(func() {
		load()
	})

	return config
}

func load() {
	config = new(Config)
	if err := envconfig.Process("", config); err != nil {
		log.Fatal().Err(err).Send()
	}
}

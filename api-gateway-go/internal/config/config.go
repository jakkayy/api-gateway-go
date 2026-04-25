package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Route struct {
	Path    string   `yaml:"path"`
	Targets []string `yaml:"targets"`
}

type Config struct {
	Routes []Route `yaml:"routes"`
}

func Load(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}

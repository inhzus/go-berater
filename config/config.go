package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Database struct {
		Dialect string
		Url     string
	}
	Jwt struct {
		SecretKey string `yaml:"secret_key"`
		Timeout   int
		Prefix    string
	}
}

var config Config

func Setup() {
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Configuration setup error: %v", err)
	}
	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Fatalf("Configuration setup error: %v", err)
	}
	d, err := yaml.Marshal(config)
	log.Printf("\nOutput: %s", string(d))
}

func GetConfig() *Config {
	return &config
}

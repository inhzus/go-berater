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
		ExpireTime uint `yaml:"expire_time"`  // Token expire time (minute)
		Header     string                     // Header name in request header, e.g. "Authorization"
		Identity   string                     // Claim identity name in gin.Context, e.g. "claims"
		Prefix     string                     // Token prefix in header, e.g. "Bearer"
		SecretKey  string `yaml:"secret_key"` // JWT signing secret key
	}
	Redis struct {
		Addr string // e.g. "localhost:6379"
	}
	Code struct {
		AccessKey    string `yaml:"access_key"`
		AccessSecret string `yaml:"access_secret"`
		ExpireTime   uint   `yaml:"expire_time"` // Code expire time in redis
		Length       int                         // e.g. code is in range(1000, 9999) when length=4
		Suffix       string                      // Key name suffix stored in redis
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

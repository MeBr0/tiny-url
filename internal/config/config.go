package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

var (
	c *Config
)

type Config struct {
	Log struct {
		Level string `yaml:"level" envconfig:"LOG_LEVEL"`
	} `yaml:"log"`

	HTTP struct {
		Port               string        `yaml:"port" envconfig:"HTTP_PORT"`
		ReadTimeout        time.Duration `yaml:"read-timeout" envconfig:"HTTP_READ_TIMEOUT"`
		WriteTimeout       time.Duration `yaml:"write-timeout" envconfig:"HTTP_WRITE_TIMEOUT"`
		MaxHeaderMegabytes int           `yaml:"max-header-megabytes" envconfig:"HTTP_MAX_HEADER_MEGABYTES"`
	} `yaml:"http"`

	Mongo struct {
		URI      string `yaml:"uri" envconfig:"MONGO_URI"`
		User     string `yaml:"user" envconfig:"MONGO_USER"`
		Password string `yaml:"password" envconfig:"MONGO_PASSWORD"`
		Name     string `yaml:"name" envconfig:"MONGO_NAME"`
	} `yaml:"mongo"`

	Redis struct {
		URI      string        `yaml:"uri" envconfig:"REDIS_URI"`
		Password string        `yaml:"password" envconfig:"REDIS_PASSWORD"`
		Database int           `yaml:"db" envconfig:"REDIS_DB"`
		TTL      time.Duration `yaml:"ttl" envconfig:"REDIS_TTL"`
	} `yaml:"redis"`

	Auth struct {
		AccessTokenTTL time.Duration `yaml:"access-token-ttl" envconfig:"AUTH_ACCESS_TOKEN_TTL"`
		PasswordSalt   string        `yaml:"password-salt" envconfig:"AUTH_PASSWORD_SALT"`
		JWT            struct {
			Key string `yaml:"key" envconfig:"AUTH_JWT_KEY"`
		} `yaml:"jwt"`
	} `yaml:"auth"`

	URL struct {
		AliasLength       int `yaml:"alias-length" envconfig:"URL_ALIAS_LENGTH"`
		DefaultExpiration int `yaml:"default-expiration" envconfig:"URL_DEFAULT_EXPIRATION"`
		CountLimit        int `yaml:"count-limit" envconfig:"URL_COUNT_LIMIT"`
	} `yaml:"url"`
}

func LoadConfig(configPath string) *Config {
	if c == nil {
		c = &Config{}

		c.readFile(configPath)
		c.readEnv()
	}

	return c
}

// File configs with values from configs file
func (c *Config) readFile(path string) {
	f, err := os.Open(path)

	if err != nil {
		processError(err)
	}

	defer f.Close()

	err = yaml.NewDecoder(f).Decode(c)

	if err != nil {
		log.Println(c)
		processError(err)
	}
}

// Read configs with values from env variables
func (c *Config) readEnv() {
	loadFromEnvFile()

	err := envconfig.Process("", c)

	if err != nil {
		processError(err)
	}
}

// Load values from .env file to system
func loadFromEnvFile() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func processError(err error) {
	log.Error(err)
	os.Exit(2)
}

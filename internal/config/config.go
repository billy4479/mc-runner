package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

const (
	Version string = "dev"
)

type EnvConfig struct {
	// Private variables
	TgToken   string `split_words:"true"`
	TgChatIds string `split_words:"true"`

	// These you can commit publically if you want
	Environment string `default:"release"`
	VitePort    int    `default:"5173" split_words:"true"`
	Port        int    `default:"4479" split_words:"true"`
	ConfigPath  string `default:"config.yml" split_words:"true"`
}

type Config struct {
	WorldDir          string `json:"world_dir"`
	DbPath            string `json:"db_path"`
	ConnectUrl        string `json:"connect_url"`
	MinutesBeforeStop int    `json:"minutes_before_stop"`

	ResticPath string `json:"restic_path"`

	Java8  string `json:"java8"`
	Java17 string `json:"java17"`
	Java21 string `json:"java21"`
	Java25 string `json:"java25"`

	LogLevel string `json:"log_level"`

	EnvConfig EnvConfig `json:"-"`
	Debug     bool      `json:"-"`
}

func NewConfig() (*Config, error, error) {
	useDotenv := os.Getenv("DONT_LOAD_DOTENV") != "yes"

	var dotenvErr error
	if useDotenv {
		dotenvErr = godotenv.Load()
		if dotenvErr != nil {
			dotenvErr = fmt.Errorf("dotenv: %w", dotenvErr)
		}
	} else {
		dotenvErr = fmt.Errorf("")
	}

	config := &Config{}
	err := envconfig.Process("", &config.EnvConfig)
	if err != nil {
		return nil, dotenvErr, fmt.Errorf("envconfig: %w", err)
	}

	b, err := os.ReadFile(config.EnvConfig.ConfigPath)
	if err != nil {
		return nil, dotenvErr, fmt.Errorf("read config: %w", err)
	}

	err = yaml.Unmarshal(b, config)
	if err != nil {
		return nil, dotenvErr, fmt.Errorf("parse config: %w", err)
	}

	config.Debug = config.EnvConfig.Environment != "release"

	return config, dotenvErr, nil
}

func (c *Config) GetLogLevel() zerolog.Level {
	if c.Debug {
		return zerolog.DebugLevel
	}

	switch c.LogLevel {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "err":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "disabled":
		return zerolog.Disabled
	default:
		return zerolog.InfoLevel
	}
}

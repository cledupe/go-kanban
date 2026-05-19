package app

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
)

type Config struct {
	Host string
	Port string
}

func LoadConfigFromEnv() (Config, error) {
	return LoadConfig(os.Getenv)
}

func LoadConfig(getenv func(string) string) (Config, error) {
	cfg := Config{
		Host: getenv("HOST"),
		Port: getenv("PORT"),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	if c.Host == "" {
		return errors.New("HOST is required")
	}

	if c.Port == "" {
		return errors.New("PORT is required")
	}

	port, err := strconv.Atoi(c.Port)
	if err != nil {
		return fmt.Errorf("PORT must be numeric: %w", err)
	}

	if port < 1 || port > 65535 {
		return errors.New("PORT must be between 1 and 65535")
	}

	return nil
}

func (c Config) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

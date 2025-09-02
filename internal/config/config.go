package config

import (
	"fmt"
	"os"
)

type Config struct {
	Address  string
	LogLevel string
	Port     string
	DBUser   string
	DBPass   string
	DBName   string
	DBHost   string
	DBPort   string
	DSN      string
}

func LoadConfig() *Config {
	return &Config{
		Address:  os.Getenv("ADDRESS"),
		LogLevel: os.Getenv("LOGLEVEL"),
		Port:     os.Getenv("PORT"),
		DSN: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_DB"),
		),
	}
}

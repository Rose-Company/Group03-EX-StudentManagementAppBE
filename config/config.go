package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Postgres PostgresConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: could not load .env file, relying on environment variables")
	}
	config := &Config{
		Postgres: PostgresConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
	}
	return config, nil
}

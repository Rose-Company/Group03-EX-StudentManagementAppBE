package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Postgres                   PostgresConfig
	JWTSecret                  string
	GoogleDriveCredentialsFile string
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
		JWTSecret:                  os.Getenv("JWT_SECRET"),
		GoogleDriveCredentialsFile: os.Getenv("GOOGLE_DRIVE_CREDENTIALS_FILE"),
	}
	return config, nil
}

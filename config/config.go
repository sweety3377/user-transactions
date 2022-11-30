package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DbLogin    string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
	HttpPort   string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		DbLogin:    os.Getenv("DB_LOGIN"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbName:     os.Getenv("DB_NAME"),
		HttpPort:   os.Getenv("HTTP_PORT"),
	}

	return cfg, nil
}

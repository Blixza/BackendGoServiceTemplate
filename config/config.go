package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Level    string
	HttpPort string
}

func LoadConfig(path string) Config {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatal("Error loading env file")
	}

	cfg := Config{}

	cfg.Host = os.Getenv("DB_HOST")
	cfg.Port = os.Getenv("DB_PORT")
	cfg.User = os.Getenv("DB_USER")
	cfg.Password = os.Getenv("DB_PASSWORD")
	cfg.DBName = os.Getenv("DB_NAME")
	cfg.SSLMode = os.Getenv("DB_SSLNAME")
	cfg.Level = os.Getenv("LEVEL")
	cfg.HttpPort = os.Getenv("HTTP_PORT")

	return cfg
}

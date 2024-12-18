package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	WebSocketURL string
	APIBaseURL   string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		WebSocketURL: os.Getenv("WEBSOCKET_URL"),
		APIBaseURL:   os.Getenv("API_BASE_URL"),
	}
}

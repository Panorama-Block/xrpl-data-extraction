package config

import (
	"log"
	"os"
)

type Config struct {
	WebSocketURL string
	APIBaseURL   string
	MongoURI     string
}

func LoadConfig() *Config {
    wsURL := os.Getenv("WEBSOCKET_URL")
    apiURL := os.Getenv("API_BASE_URL")
    mongoURI := os.Getenv("MONGO_URI")

    // Validar variáveis obrigatórias
    if wsURL == "" || apiURL == "" || mongoURI == "" {
        log.Fatalf("❌ Variáveis obrigatórias não definidas! Certifique-se de definir WEBSOCKET_URL, API_BASE_URL e MONGO_URI")
    }

    return &Config{
        WebSocketURL: wsURL,
        APIBaseURL:   apiURL,
        MongoURI:     mongoURI,
    }
}



package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Panorama-Block/xrpl-data-extraction/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

// Conectar ao MongoDB usando URI do .env
func ConnectMongoDB(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return fmt.Errorf("Erro ao conectar ao MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("Erro ao fazer ping no MongoDB: %v", err)
	}

	log.Println("✅ Conectado ao MongoDB com sucesso!")
	Client = client
	return nil
}

// Retorna a coleção de ledgers
func GetLedgerCollection() *mongo.Collection {
	return Client.Database("xrpl").Collection("ledger")
}

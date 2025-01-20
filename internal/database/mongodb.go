package database

import (
	"context"
	"log"
	"time"

	"github.com/Panorama-Block/xrpl-data-extraction/config"
	"go.mongodb.org/mongo-driver/bson"
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
		log.Printf("❌ Erro ao conectar ao MongoDB: %v", err)
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Printf("❌ Erro ao pingar o MongoDB: %v", err)
		return err
	}

	log.Println("✅ Conectado ao MongoDB com sucesso!")
	Client = client
	return nil
}

// Retorna a coleção de ledgers
func GetLedgerCollection() *mongo.Collection {
	return Client.Database("xrpl").Collection("ledger")
}
// Retorna a coleção de accounts
func GetAccountCollection() *mongo.Collection {
	return Client.Database("xrpl").Collection("accounts")
}
// GetTransactionCollection retorna a coleção transactions do banco de dados MongoDB
func GetTransactionCollection() *mongo.Collection {
	return Client.Database("xrpl").Collection("transactions")
}

func CreateIndexes() error {
    collection := GetLedgerCollection()

    // Índice único para ledger_index
    _, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
        Keys:    bson.D{{Key: "ledger_index", Value: 1}},
        Options: options.Index().SetUnique(true),
    })
    if err != nil {
        log.Printf("⚠️ Índice para ledger_index já existe: %v", err)
				return nil
    }

    // Índice TTL para expirar dados antigos (30 dias)
    _, err = collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
        Keys:    bson.D{{Key: "created_at", Value: 1}},
        Options: options.Index().SetExpireAfterSeconds(30 * 24 * 60 * 60),
    })
    if err != nil {
        log.Printf("⚠️ Índice TTL para created_at já existe: %v", err)
				
    }

    log.Println("✅ Índices criados com sucesso!")
    return nil
}


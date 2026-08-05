package database

import (
	"context"
	"time"

	"github.com/olabanji12-ojo/church-backend/config"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

// ConnectDB initializes the MongoDB connection
func ConnectDB() *mongo.Database {
	uri := config.GetEnv("MONGO_URI", "mongodb://localhost:27017")
	dbName := config.GetEnv("DB_NAME", "church_match")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logrus.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		logrus.Fatalf("❌ Failed to ping MongoDB: %v", err)
	}

	logrus.Println("✅ Successfully connected to MongoDB")
	MongoClient = client
	return client.Database(dbName)
}

package main

import (
	"context"
	"log"
	"time"

	"github.com/olabanji12-ojo/church-backend/config"
	"github.com/olabanji12-ojo/church-backend/database"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	log.Println("⚠️ Wiping MongoDB collections to start fresh...")

	// 1. Load environment variables
	config.LoadEnv()

	// 2. Connect to Database
	db := database.ConnectDB()

	// List of all collections in the application
	collections := []string{"users", "swipes", "matches", "messages", "prayers", "reports"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, colName := range collections {
		collection := db.Collection(colName)
		_, err := collection.DeleteMany(ctx, bson.M{})
		if err != nil {
			log.Printf("❌ Failed to clear collection %s: %v", colName, err)
		} else {
			log.Printf("🗑️ Cleared all documents from collection: %s", colName)
		}
	}

	log.Println("🎉 Database reset successfully! You can now sign up/create accounts fresh.")
}

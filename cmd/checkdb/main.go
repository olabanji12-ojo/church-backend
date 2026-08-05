package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/olabanji12-ojo/church-backend/config"
	"github.com/olabanji12-ojo/church-backend/database"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	log.Println("🔍 Checking live MongoDB database users and embeddings (RAW MAP)...")

	// 1. Load env
	config.LoadEnv()

	// 2. Connect
	db := database.ConnectDB()
	collection := db.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("❌ Failed to query users: %v", err)
	}
	defer cursor.Close(ctx)

	var users []bson.M
	if err = cursor.All(ctx, &users); err != nil {
		log.Fatalf("❌ Failed to decode users: %v", err)
	}

	fmt.Printf("\n--- 📋 RAW DATABASE SUMMARY: %d USERS FOUND ---\n", len(users))
	for i, u := range users {
		jsonBytes, _ := json.MarshalIndent(u, "", "  ")
		fmt.Printf("[%d] User:\n%s\n", i+1, string(jsonBytes))
		fmt.Println("--------------------------------------------------")
	}
}

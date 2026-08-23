package main

import (
	"context"
	"log"
	"time"

	"github.com/olabanji12-ojo/church-backend/config"
	"github.com/olabanji12-ojo/church-backend/database"
	"github.com/olabanji12-ojo/church-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	log.Println("🔑 Ensuring all users in MongoDB have password_hash set to 'password123'...")

	config.LoadEnv()
	db := database.ConnectDB()
	collection := db.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	passHash, err := utils.HashPassword("password123")
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Update all users whose password_hash is empty or missing
	filter := bson.M{
		"$or": []bson.M{
			{"password_hash": ""},
			{"password_hash": bson.M{"$exists": false}},
		},
	}
	update := bson.M{"$set": bson.M{"password_hash": passHash}}

	res, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Fatalf("Failed to update passwords: %v", err)
	}

	log.Printf("🎉 Done! Updated %d users to have valid password_hash for 'password123'.", res.ModifiedCount)
}

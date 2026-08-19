package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/olabanji12-ojo/church-backend/config"
	"github.com/olabanji12-ojo/church-backend/database"
	"github.com/olabanji12-ojo/church-backend/models"
	"github.com/olabanji12-ojo/church-backend/repositories"
	"github.com/olabanji12-ojo/church-backend/services"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	config.LoadEnv()
	db := database.ConnectDB()

	userRepo := repositories.NewUserRepository(db)
	matchRepo := repositories.NewMatchRepository(db)
	msgRepo := repositories.NewMessageRepository(db)

	swipeService := services.NewSwipeService(userRepo, matchRepo, msgRepo, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find Olabanji Ojo (man)
	var banji models.User
	err := db.Collection("users").FindOne(ctx, bson.M{"email": "ojoolabanji59@gmail.com"}).Decode(&banji)
	if err != nil {
		log.Fatalf("Failed to find Banji: %v", err)
	}

	fmt.Printf("=== User 1: %s %s (Gender: %s, InterestedIn: %s) ===\n", banji.FirstName, banji.LastName, banji.Gender, banji.InterestedIn)

	// Call GetDiscoveryFeed for Banji
	feed, err := swipeService.GetDiscoveryFeed(banji.ID)
	if err != nil {
		fmt.Printf("GetDiscoveryFeed error: %v\n", err)
	} else {
		fmt.Printf("Discovery feed candidates count for %s: %d\n", banji.FirstName, len(feed))
		for i, cand := range feed {
			fmt.Printf("  %d. %s %s | Gender: %s | Denomination: %s | Photos: %d\n",
				i+1, cand.FirstName, cand.LastName, cand.Gender, cand.Denomination, len(cand.Photos))
		}
	}

	// Find Barakat (woman)
	var barakat models.User
	err = db.Collection("users").FindOne(ctx, bson.M{"email": "alamubarakatabisola@gmail.com"}).Decode(&barakat)
	if err != nil {
		log.Fatalf("Failed to find Barakat: %v", err)
	}

	fmt.Printf("\n=== User 2: %s %s (Gender: %s, InterestedIn: %s) ===\n", barakat.FirstName, barakat.LastName, barakat.Gender, barakat.InterestedIn)

	// Call GetDiscoveryFeed for Barakat
	feed2, err := swipeService.GetDiscoveryFeed(barakat.ID)
	if err != nil {
		fmt.Printf("GetDiscoveryFeed error: %v\n", err)
	} else {
		fmt.Printf("Discovery feed candidates count for %s: %d\n", barakat.FirstName, len(feed2))
		for i, cand := range feed2 {
			fmt.Printf("  %d. %s %s | Gender: %s | Denomination: %s | Photos: %d\n",
				i+1, cand.FirstName, cand.LastName, cand.Gender, cand.Denomination, len(cand.Photos))
		}
	}
}

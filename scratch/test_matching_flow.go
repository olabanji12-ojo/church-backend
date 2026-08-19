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

	var banji models.User
	err := db.Collection("users").FindOne(ctx, bson.M{"email": "ojoolabanji59@gmail.com"}).Decode(&banji)
	if err != nil {
		log.Fatalf("Failed to find Banji: %v", err)
	}

	var barakat models.User
	err = db.Collection("users").FindOne(ctx, bson.M{"email": "alamubarakatabisola@gmail.com"}).Decode(&barakat)
	if err != nil {
		log.Fatalf("Failed to find Barakat: %v", err)
	}

	fmt.Printf("Initial State:\n  Banji ID: %s\n  Barakat ID: %s\n", banji.ID.Hex(), barakat.ID.Hex())

	// Clear any prior matches between Banji and Barakat for testing
	_, _ = db.Collection("matches").DeleteMany(ctx, bson.M{
		"users": bson.M{"$all": []interface{}{banji.ID, barakat.ID}},
	})

	// STEP 1: Banji likes Barakat (SwipeRight)
	fmt.Println("\n--- STEP 1: Banji likes Barakat (SwipeRight) ---")
	match1, err := swipeService.SwipeRight(banji.ID, barakat.ID)
	if err != nil {
		log.Fatalf("Banji SwipeRight failed: %v", err)
	}
	fmt.Printf("Match Status: %s (Initiator: %s)\n", match1.Status, match1.Users[0].Hex())

	// STEP 2: Check Barakat's pending requests
	fmt.Println("\n--- STEP 2: Fetch Barakat's Pending Requests ---")
	requests, err := swipeService.GetPendingLikes(barakat.ID)
	if err != nil {
		log.Fatalf("GetPendingLikes for Barakat failed: %v", err)
	}
	fmt.Printf("Barakat pending requests count: %d\n", len(requests))
	for _, req := range requests {
		fmt.Printf("  Incoming request from: %s %s (%s)\n", req.FirstName, req.LastName, req.ID.Hex())
	}

	// STEP 3: Barakat likes Banji (SwipeRight back to form mutual match)
	fmt.Println("\n--- STEP 3: Barakat likes Banji back (SwipeRight) ---")
	match2, err := swipeService.SwipeRight(barakat.ID, banji.ID)
	if err != nil {
		log.Fatalf("Barakat SwipeRight failed: %v", err)
	}
	fmt.Printf("Match Status after mutual like: %s\n", match2.Status)

	// STEP 4: Check GetMatches for both users
	fmt.Println("\n--- STEP 4: Check GetMatches for both users ---")
	banjiMatches, err := swipeService.GetMatches(banji.ID)
	if err != nil {
		log.Fatalf("GetMatches for Banji failed: %v", err)
	}
	fmt.Printf("Banji matches count: %d\n", len(banjiMatches))

	barakatMatches, err := swipeService.GetMatches(barakat.ID)
	if err != nil {
		log.Fatalf("GetMatches for Barakat failed: %v", err)
	}
	fmt.Printf("Barakat matches count: %d\n", len(barakatMatches))
}

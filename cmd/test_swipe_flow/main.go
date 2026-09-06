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
	log.Println("🧪 Testing Swipe Right -> Discovery Feed exclusion flow...")

	config.LoadEnv()
	db := database.ConnectDB()
	userRepo := repositories.NewUserRepository(db)
	matchRepo := repositories.NewMatchRepository(db)
	msgRepo := repositories.NewMessageRepository(db)
	swipeSvc := services.NewSwipeService(userRepo, matchRepo, msgRepo, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var maleUser, femaleUser models.User
	err1 := db.Collection("users").FindOne(ctx, bson.M{"email": "david.emmanuel@covenant.app"}).Decode(&maleUser)
	err2 := db.Collection("users").FindOne(ctx, bson.M{"email": "grace.adebayo@covenant.app"}).Decode(&femaleUser)

	if err1 != nil || err2 != nil {
		log.Fatalf("Users not found: err1=%v, err2=%v", err1, err2)
	}

	fmt.Printf("1. Initial Discovery Feed for David Emmanuel (%s):\n", maleUser.Email)
	feed1, _ := swipeSvc.GetDiscoveryFeed(maleUser.ID)
	for i, u := range feed1 {
		fmt.Printf("   [%d] %s (%s)\n", i+1, u.FirstName, u.Email)
	}

	fmt.Printf("\n2. David Emmanuel swipes RIGHT on Grace Adebayo (%s)...\n", femaleUser.Email)
	match, err := swipeSvc.SwipeRight(maleUser.ID, femaleUser.ID)
	if err != nil {
		log.Fatalf("Swipe failed: %v", err)
	}
	fmt.Printf("   Swipe Result Status: %s, Match ID: %s\n", match.Status, match.ID.Hex())

	fmt.Printf("\n3. Refetched Discovery Feed for David Emmanuel after swipe:\n")
	feed2, _ := swipeSvc.GetDiscoveryFeed(maleUser.ID)
	graceFound := false
	for i, u := range feed2 {
		fmt.Printf("   [%d] %s (%s)\n", i+1, u.FirstName, u.Email)
		if u.Email == femaleUser.Email {
			graceFound = true
		}
	}

	if graceFound {
		fmt.Println("\n❌ FAIL: Grace Adebayo STILL APPEARED in David's Discovery Feed after swipe!")
	} else {
		fmt.Println("\n✅ SUCCESS: Grace Adebayo was correctly EXCLUDED from David's Discovery Feed!")
	}

	fmt.Printf("\n4. Checking Discovery / Pending Requests Feed for Grace Adebayo:\n")
	feedGrace, _ := swipeSvc.GetDiscoveryFeed(femaleUser.ID)
	fmt.Printf("   Grace Discovery Feed Count: %d\n", len(feedGrace))
	pendingGrace, _ := swipeSvc.GetPendingLikes(femaleUser.ID)
	fmt.Printf("   Grace Pending Requests Count: %d\n", len(pendingGrace))
	for i, u := range pendingGrace {
		fmt.Printf("   [%d] Connection Request From: %s %s (%s)\n", i+1, u.FirstName, u.LastName, u.Email)
	}

	// Clean up test match document so we leave DB clean
	_ = matchRepo.DeleteMatch(match.ID)
	fmt.Println("\n🧹 Test match cleaned up.")
}

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
	log.Println("🔍 Diagnosing Discovery Feed candidates for all users...")

	config.LoadEnv()
	db := database.ConnectDB()
	userRepo := repositories.NewUserRepository(db)
	matchRepo := repositories.NewMatchRepository(db)
	msgRepo := repositories.NewMessageRepository(db)
	swipeSvc := services.NewSwipeService(userRepo, matchRepo, msgRepo, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to fetch users: %v", err)
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		log.Fatalf("Failed to parse users: %v", err)
	}

	fmt.Printf("Total users in DB: %d\n\n", len(users))

	for _, u := range users {
		if u.IsGuest {
			continue
		}

		feed, err := swipeSvc.GetDiscoveryFeed(u.ID)
		if err != nil {
			fmt.Printf("❌ ERROR for %s (%s): %v\n", u.FirstName, u.Email, err)
			continue
		}

		fmt.Printf("👤 USER: %s %s | Email: %s | Gender: '%s' | InterestedIn: '%s' | MinAge: %d | MaxAge: %d | PrefDenom: '%s'\n",
			u.FirstName, u.LastName, u.Email, u.Gender, u.InterestedIn, u.MinAgePref, u.MaxAgePref, u.PreferredDenomination)
		fmt.Printf("   👉 Discovery Feed Count: %d candidates\n", len(feed))
		for i, cand := range feed {
			fmt.Printf("      [%d] %s %s (%s, Genotype: %s, Age: %d)\n",
				i+1, cand.FirstName, cand.LastName, cand.Gender, cand.Genotype, calculateAge(cand.DateOfBirth))
		}
		fmt.Println("--------------------------------------------------------------------------------")
	}
}

func calculateAge(dob time.Time) int {
	if dob.IsZero() {
		return 0
	}
	now := time.Now()
	age := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		age--
	}
	return age
}

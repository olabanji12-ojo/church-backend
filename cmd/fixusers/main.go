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
	log.Println("⚡ Fixing profile data for test users in MongoDB...")

	// 1. Load env
	config.LoadEnv()

	// 2. Connect
	db := database.ConnectDB()
	collection := db.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// User 1: Emmanuel (Man looking for Woman)
	dob1, _ := time.Parse("2006-01-02", "1995-05-15")
	_, err := collection.UpdateOne(ctx,
		bson.M{"email": "emmanuelojo291@gmail.com"},
		bson.M{"$set": bson.M{
			"gender":                 "man",
			"interested_in":          "woman",
			"dob":                    dob1,
			"preferred_denomination": "Any",
			"preferred_church_freq":  "Any",
			"min_age_pref":           20,
			"max_age_pref":           35,
		}},
	)
	if err != nil {
		log.Printf("❌ Failed to update Emmanuel: %v", err)
	} else {
		log.Println("✅ Fixed profile for emmanuelojo291@gmail.com")
	}

	// User 5: Emma Kreator (Woman looking for Man)
	dob5, _ := time.Parse("2006-01-02", "1997-08-20")
	_, err = collection.UpdateOne(ctx,
		bson.M{"email": "emmakreator@gmail.com"},
		bson.M{"$set": bson.M{
			"gender":                 "woman",
			"interested_in":          "man",
			"dob":                    dob5,
			"preferred_denomination": "Any",
			"preferred_church_freq":  "Any",
		}},
	)
	if err != nil {
		log.Printf("❌ Failed to update Emma Kreator: %v", err)
	} else {
		log.Println("✅ Fixed profile for emmakreator@gmail.com")
	}

	// User 6: Olabanji Tofunmie (Woman looking for Man)
	dob6, _ := time.Parse("2006-01-02", "1998-12-10")
	_, err = collection.UpdateOne(ctx,
		bson.M{"email": "tofunmieolabanji@gmail.com"},
		bson.M{"$set": bson.M{
			"gender":                 "woman",
			"interested_in":          "man",
			"dob":                    dob6,
			"preferred_denomination": "Any",
			"preferred_church_freq":  "Any",
		}},
	)
	if err != nil {
		log.Printf("❌ Failed to update Tofunmie: %v", err)
	} else {
		log.Println("✅ Fixed profile for tofunmieolabanji@gmail.com")
	}

	log.Println("🎉 Profiles fixed successfully! Now let's run backfill to regenerate embeddings.")
}

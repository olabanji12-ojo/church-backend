package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/olabanji12-ojo/church-backend/config"
	"github.com/olabanji12-ojo/church-backend/database"
	"github.com/olabanji12-ojo/church-backend/models"
	"github.com/olabanji12-ojo/church-backend/services"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	log.Println("⚡ Starting embedding backfill script...")

	// 1. Load environment variables
	config.LoadEnv()
	log.Printf("DEBUG: Hugging_Face_Key = '%s' (len: %d)", os.Getenv("Hugging_Face_Key"), len(os.Getenv("Hugging_Face_Key")))

	// 2. Connect to Database
	db := database.ConnectDB()
	collection := db.Collection("users")

	// 3. Initialize Embedding Service
	embedService := services.NewEmbeddingService()

	// 4. Fetch all users
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("❌ Failed to query users: %v", err)
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		log.Fatalf("❌ Failed to parse users: %v", err)
	}

	log.Printf("📋 Found %d total users in the database.", len(users))

	successCount := 0
	skipCount := 0

	for _, user := range users {
		// Skip if user is guest or already has embedding populated
		if user.IsGuest {
			log.Printf("⏭️ Skipping guest user: %s", user.Email)
			skipCount++
			continue
		}

		// Generate embeddings for all active non-guest users

		log.Printf("🌀 Generating embedding for: %s %s (%s)", user.FirstName, user.LastName, user.Email)
		
		// 1. Compile profile text representation
		text := embedService.GenerateUserText(&user)
		embedding, err := embedService.GetEmbedding(text)
		if err != nil {
			log.Printf("❌ Failed to generate profile embedding for %s: %v", user.Email, err)
			continue
		}

		// 2. Compile partner preference text representation
		prefText := embedService.GeneratePartnerPreferenceText(&user)
		prefEmbedding, err := embedService.GetEmbedding(prefText)
		if err != nil {
			log.Printf("❌ Failed to generate partner preference embedding for %s: %v", user.Email, err)
			continue
		}

		// Save both embeddings back to MongoDB
		updateCtx, updateCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, err = collection.UpdateOne(
			updateCtx,
			bson.M{"_id": user.ID},
			bson.M{"$set": bson.M{
				"profile_embedding":      embedding,
				"partner_pref_embedding": prefEmbedding,
				"updated_at":             time.Now(),
			}},
		)
		updateCancel()

		if err != nil {
			log.Printf("❌ Failed to update MongoDB for %s: %v", user.Email, err)
			continue
		}

		log.Printf("✅ Saved both embeddings for %s successfully!", user.Email)
		successCount++

		// CRITICAL: Sleep for 500ms to respect Hugging Face free API Rate Limits
		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("🎉 Backfill completed! Successful: %d, Skipped/Failed: %d", successCount, len(users)-successCount)
}

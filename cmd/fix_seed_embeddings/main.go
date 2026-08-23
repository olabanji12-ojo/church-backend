package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/olabanji12-ojo/church-backend/config"
	"github.com/olabanji12-ojo/church-backend/database"
	"github.com/olabanji12-ojo/church-backend/models"
	"github.com/olabanji12-ojo/church-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func generateDummyVector(dim int) []float32 {
	vec := make([]float32, dim)
	var sumSq float32
	for i := 0; i < dim; i++ {
		val := rand.Float32()*2 - 1
		vec[i] = val
		sumSq += val * val
	}
	norm := float32(1.0 / (1.0 + float64(sumSq)*0.001))
	for i := 0; i < dim; i++ {
		vec[i] *= norm
	}
	return vec
}

func main() {
	log.Println("🛠️ Fixing Seeded Users: Setting 384-dim embeddings & standard password hashes...")

	config.LoadEnv()
	db := database.ConnectDB()
	collection := db.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Hash password123
	passHash, err := utils.HashPassword("password123")
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	seedEmails := []string{
		"grace.adebayo@covenant.app",
		"joy.okonkwo@covenant.app",
		"esther.danjuma@covenant.app",
		"blessing.nwachukwu@covenant.app",
		"hannah.bello@covenant.app",
		"david.emmanuel@covenant.app",
		"samuel.ezekiel@covenant.app",
		"daniel.ogundipe@covenant.app",
		"caleb.martins@covenant.app",
		"joshua.chukwu@covenant.app",
	}

	rand.Seed(42)

	for _, email := range seedEmails {
		profEmb := generateDummyVector(384)
		prefEmb := generateDummyVector(384)

		res, err := collection.UpdateOne(
			ctx,
			bson.M{"email": strings.ToLower(email)},
			bson.M{"$set": bson.M{
				"password_hash":          passHash,
				"profile_embedding":      profEmb,
				"partner_pref_embedding": prefEmb,
				"is_verified":            true,
				"updated_at":             time.Now(),
			}},
		)
		if err != nil {
			log.Printf("❌ Failed to update %s: %v", email, err)
		} else {
			log.Printf("✅ Updated embeddings & password for %s (Matched: %d, Modified: %d)", email, res.MatchedCount, res.ModifiedCount)
		}
	}

	// 2. Test Login API endpoint on local running backend (http://localhost:8081/api/v1/auth/login)
	log.Println("\n🧪 Testing Login API for david.emmanuel@covenant.app...")
	loginBody, _ := json.Marshal(map[string]string{
		"email":    "david.emmanuel@covenant.app",
		"password": "password123",
	})
	resp, err := http.Post("http://localhost:8081/api/v1/auth/login", "application/json", bytes.NewBuffer(loginBody))
	if err != nil {
		log.Printf("⚠️ Login request failed (is backend running?): %v", err)
	} else {
		defer resp.Body.Close()
		var resMap map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&resMap)
		log.Printf("Status: %d, Response: %+v", resp.StatusCode, resMap)

		if data, ok := resMap["data"].(map[string]interface{}); ok {
			if token, ok := data["token"].(string); ok && token != "" {
				log.Println("🎉 LOGIN SUCCESSFUL! Token generated.")

				// Test discovery feed endpoint for this token!
				req, _ := http.NewRequest("GET", "http://localhost:8081/api/v1/users/discovery", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				client := &http.Client{}
				discResp, err := client.Do(req)
				if err == nil {
					defer discResp.Body.Close()
					var discData map[string]interface{}
					json.NewDecoder(discResp.Body).Decode(&discData)
					if feed, ok := discData["data"].([]interface{}); ok {
						log.Printf("📱 Discovery Feed returned %d candidates for David Emmanuel!", len(feed))
						for i, item := range feed {
							if cand, ok := item.(map[string]interface{}); ok {
								fmt.Printf("   [%d] Candidate: %v %v (%v%% match)\n", i+1, cand["first_name"], cand["last_name"], cand["match_score"])
							}
						}
					}
				}
			}
		}
	}

	// 3. Test discovery feed for user 'tofunmieolabanji@gmail.com'
	var tofunmieUser models.User
	err = collection.FindOne(ctx, bson.M{"email": "tofunmieolabanji@gmail.com"}).Decode(&tofunmieUser)
	if err == nil {
		token, err := utils.GenerateJWT(tofunmieUser.ID)
		if err == nil {
			req, _ := http.NewRequest("GET", "http://localhost:8081/api/v1/users/discovery", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			client := &http.Client{}
			discResp, err := client.Do(req)
			if err == nil {
				defer discResp.Body.Close()
				var discData map[string]interface{}
				json.NewDecoder(discResp.Body).Decode(&discData)
				if feed, ok := discData["data"].([]interface{}); ok {
					log.Printf("\n📱 Discovery Feed returned %d female candidates for tofunmieolabanji@gmail.com!", len(feed))
					for i, item := range feed {
						if cand, ok := item.(map[string]interface{}); ok {
							fmt.Printf("   [%d] Candidate: %v %v (%v%% match)\n", i+1, cand["first_name"], cand["last_name"], cand["match_score"])
						}
					}
				}
			}
		}
	}
}

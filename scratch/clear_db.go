package main

import (

	"context"
	"fmt"
	"log"
	"time"

	"github.com/olabanji12-ojo/church-backend/config"
	"github.com/olabanji12-ojo/church-backend/database"
	"go.mongodb.org/mongo-driver/bson"
	
)

func main() {
	config.LoadEnv()
	db := database.ConnectDB()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collections := []string{"users", "matches", "messages", "prayers", "reports", "notifications"}

	fmt.Println("🧹 Starting database cleanup...")
	for _, colName := range collections {
		res, err := db.Collection(colName).DeleteMany(ctx, bson.M{})
		if err != nil {
			log.Printf("❌ Failed to clear collection %s: %v", colName, err)
		} else {
			fmt.Printf("✅ Cleared collection '%s': %d documents deleted.\n", colName, res.DeletedCount)
		}
	}
	fmt.Println("✨ Database cleanup complete! All users, matches, and messages removed.")
}

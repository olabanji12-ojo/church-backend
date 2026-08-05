package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Prayer struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	AuthorID  primitive.ObjectID   `bson:"author_id" json:"author_id"`
	Content   string               `bson:"content" json:"content"`
	AmenCount int                  `bson:"amen_count" json:"amen_count"`
	AmensBy   []primitive.ObjectID `bson:"amens_by" json:"amens_by"` // Array of User IDs who clicked "Amen"
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
}

type PrayerFeedItem struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	AuthorID  primitive.ObjectID   `bson:"author_id" json:"author_id"`
	Content   string               `bson:"content" json:"content"`
	AmenCount int                  `bson:"amen_count" json:"amen_count"`
	AmensBy   []primitive.ObjectID `bson:"amens_by" json:"amens_by"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	
	// Joined User Details
	AuthorName  string `bson:"author_name" json:"author_name"`
	AuthorPhoto string `bson:"author_photo" json:"author_photo"`
}

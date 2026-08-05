package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MatchID   primitive.ObjectID `bson:"match_id" json:"match_id"`
	SenderID  primitive.ObjectID `bson:"sender_id" json:"sender_id"`
	Content   string             `bson:"content" json:"content"`
	Type      string             `bson:"type,omitempty" json:"type,omitempty"` // e.g. "text", "prayer"
	IsRead    bool               `bson:"is_read" json:"is_read"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

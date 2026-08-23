package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Match struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Users     []primitive.ObjectID `bson:"users" json:"users"`   // Array of the two user IDs
	Status    string               `bson:"status" json:"status"` // "pending", "matched", "rejected"
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
}

// MatchResponse is a composite struct used by the frontend Inbox to display a match along with the last message.
type MatchResponse struct {
	MatchID     primitive.ObjectID `json:"match_id"`
	User        User               `json:"user"`
	LastMessage *Message           `json:"last_message"`
}

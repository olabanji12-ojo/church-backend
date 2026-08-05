package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Report struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ReporterID  primitive.ObjectID `bson:"reporter_id" json:"reporter_id"`
	ReportedID  primitive.ObjectID `bson:"reported_id" json:"reported_id"`
	Reason      string             `bson:"reason" json:"reason"`
	Status      string             `bson:"status" json:"status"` // "pending", "reviewed", "banned"
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

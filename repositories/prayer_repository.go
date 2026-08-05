package repositories

import (
	"context"
	"time"

	"github.com/olabanji12-ojo/church-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PrayerRepository struct {
	db *mongo.Database
}

func NewPrayerRepository(db *mongo.Database) *PrayerRepository {
	return &PrayerRepository{db: db}
}

// CreatePrayer inserts a new prayer request
func (pr *PrayerRepository) CreatePrayer(prayer models.Prayer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prayer.CreatedAt = time.Now()
	prayer.AmenCount = 0
	prayer.AmensBy = []primitive.ObjectID{}

	_, err := pr.db.Collection("prayers").InsertOne(ctx, prayer)
	return err
}

// GetRecentPrayers fetches the most recent prayers with author details
func (pr *PrayerRepository) GetRecentPrayers(limit int64) ([]models.PrayerFeedItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "author_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "author"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$author"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "author_name", Value: "$author.first_name"},
			{Key: "author_photo", Value: bson.D{
				{Key: "$arrayElemAt", Value: []interface{}{"$author.photos", 0}},
			}},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "author", Value: 0},
		}}},
	}

	cursor, err := pr.db.Collection("prayers").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var prayers []models.PrayerFeedItem
	if err = cursor.All(ctx, &prayers); err != nil {
		return nil, err
	}

	return prayers, nil
}

// AddAmen adds a user's ID to the AmensBy array and increments the count
func (pr *PrayerRepository) AddAmen(prayerID, userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Only match if the user is NOT in the amens_by array to prevent double incrementing
	filter := bson.M{
		"_id":      prayerID,
		"amens_by": bson.M{"$ne": userID},
	}

	update := bson.M{
		"$addToSet": bson.M{"amens_by": userID},
		"$inc":      bson.M{"amen_count": 1},
	}

	_, err := pr.db.Collection("prayers").UpdateOne(ctx, filter, update)
	return err
}

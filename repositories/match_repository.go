package repositories

import (
	"context"
	"time"

	"github.com/olabanji12-ojo/church-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MatchRepository struct {
	db *mongo.Database
}

func NewMatchRepository(db *mongo.Database) *MatchRepository {
	return &MatchRepository{db: db}
}

// CreateMatch explicitly creates a new match document
func (mr *MatchRepository) CreateMatch(match models.Match) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	match.CreatedAt = time.Now()
	match.UpdatedAt = time.Now()

	_, err := mr.db.Collection("matches").InsertOne(ctx, match)
	return err
}

// FindMatchBetweenUsers finds if a match document already exists between two users
func (mr *MatchRepository) FindMatchBetweenUsers(userID1, userID2 primitive.ObjectID) (*models.Match, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"users": bson.M{
			"$all": []primitive.ObjectID{userID1, userID2},
		},
	}

	var match models.Match
	err := mr.db.Collection("matches").FindOne(ctx, filter).Decode(&match)
	if err != nil {
		return nil, err
	}

	return &match, nil
}

// GetMatchByID fetches a match by its ObjectID
func (mr *MatchRepository) GetMatchByID(matchID primitive.ObjectID) (*models.Match, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var match models.Match
	err := mr.db.Collection("matches").FindOne(ctx, bson.M{"_id": matchID}).Decode(&match)
	if err != nil {
		return nil, err
	}

	return &match, nil
}

// UpdateMatchStatus updates the status of a match (e.g. from pending to matched)
func (mr *MatchRepository) UpdateMatchStatus(matchID primitive.ObjectID, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := mr.db.Collection("matches").UpdateOne(
		ctx,
		bson.M{"_id": matchID},
		bson.M{
			"$set": bson.M{
				"status":     status,
				"updated_at": time.Now(),
			},
		},
	)

	return err
}

// GetUserMatches gets all successful matches for a user
func (mr *MatchRepository) GetUserMatches(userID primitive.ObjectID) ([]models.Match, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"users":  userID,
		"status": "matched",
	}

	cursor, err := mr.db.Collection("matches").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var matches []models.Match
	if err = cursor.All(ctx, &matches); err != nil {
		return nil, err
	}

	return matches, nil
}

// GetAllUserInteractions gets all match documents where the user is involved, regardless of status
func (mr *MatchRepository) GetAllUserInteractions(userID primitive.ObjectID) ([]models.Match, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"users": userID,
	}

	cursor, err := mr.db.Collection("matches").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var matches []models.Match
	if err = cursor.All(ctx, &matches); err != nil {
		return nil, err
	}

	return matches, nil
}

// DeleteMatch removes a match document from the database
func (mr *MatchRepository) DeleteMatch(matchID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := mr.db.Collection("matches").DeleteOne(ctx, bson.M{"_id": matchID})
	return err
}

// GetPendingLikesForUser gets pending match documents where userID is involved
func (mr *MatchRepository) GetPendingLikesForUser(userID primitive.ObjectID) ([]models.Match, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"users":  userID,
		"status": "pending",
	}

	cursor, err := mr.db.Collection("matches").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var matches []models.Match
	if err = cursor.All(ctx, &matches); err != nil {
		return nil, err
	}

	return matches, nil
}


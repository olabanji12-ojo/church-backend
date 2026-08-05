package repositories

import (
	"context"
	"time"

	"github.com/olabanji12-ojo/church-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepository struct {
	db *mongo.Database
}

func NewMessageRepository(db *mongo.Database) *MessageRepository {
	return &MessageRepository{db: db}
}

// CreateMessage inserts a new chat message
func (mr *MessageRepository) CreateMessage(message models.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	message.CreatedAt = time.Now()

	_, err := mr.db.Collection("messages").InsertOne(ctx, message)
	return err
}

// GetMessagesForMatch fetches all messages for a specific match, sorted by creation time
func (mr *MessageRepository) GetMessagesForMatch(matchID primitive.ObjectID) ([]models.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"match_id": matchID}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}) // Sort ascending

	cursor, err := mr.db.Collection("messages").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

// GetLastMessageForMatch fetches the single most recent message for a match
func (mr *MessageRepository) GetLastMessageForMatch(matchID primitive.ObjectID) (*models.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"match_id": matchID}
	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort descending

	var message models.Message
	err := mr.db.Collection("messages").FindOne(ctx, filter, opts).Decode(&message)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No message yet, which is totally fine
		}
		return nil, err
	}

	return &message, nil
}

// DeleteMessagesForMatch deletes all messages associated with a given match
func (mr *MessageRepository) DeleteMessagesForMatch(matchID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := mr.db.Collection("messages").DeleteMany(ctx, bson.M{"match_id": matchID})
	return err
}

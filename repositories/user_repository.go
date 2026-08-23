package repositories

import (
	"context"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/olabanji12-ojo/church-backend/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser inserts a new user into the MongoDB users collection
func (ur *UserRepository) CreateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user.ID = primitive.NewObjectID() // Explicitly generate the ID before inserting!
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := ur.db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		logrus.Error("Error inserting user: ", err)
		return err
	}
	logrus.Info("User inserted successfully")
	return nil
}

// FindUserByEmail searches for a user by their email (case-insensitive)
func (ur *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cleanEmail := strings.TrimSpace(strings.ToLower(email))
	filter := bson.M{"email": primitive.Regex{Pattern: "^" + cleanEmail + "$", Options: "i"}}

	var user models.User
	err := ur.db.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		logrus.Warn("User not found with email: ", email)
		return nil, err
	}
	return &user, nil
}

// FindUserByID fetches a user by MongoDB ObjectID
func (ur *UserRepository) FindUserByID(userID primitive.ObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := ur.db.Collection("users").FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserByID updates specific fields for a user
func (ur *UserRepository) UpdateUserByID(userID primitive.ObjectID, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Automatically append updated_at
	update["updated_at"] = time.Now()

	_, err := ur.db.Collection("users").UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$set": update},
	)

	return err
}

// BlockUser adds a user ID to the current user's blocked list
func (ur *UserRepository) BlockUser(reporterID, blockedID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := ur.db.Collection("users").UpdateOne(
		ctx,
		bson.M{"_id": reporterID},
		bson.M{"$addToSet": bson.M{"blocked_users": blockedID}},
	)

	return err
}

// GetUsersByIDs fetches multiple users by their IDs
func (ur *UserRepository) GetUsersByIDs(userIDs []primitive.ObjectID) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	cursor, err := ur.db.Collection("users").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// GetTargetGender determines the target gender string for discovery (Female or Male)
func GetTargetGender(user *models.User) string {
	interestedIn := strings.TrimSpace(strings.ToLower(user.InterestedIn))
	if interestedIn == "female" || interestedIn == "women" || interestedIn == "woman" {
		return "Female"
	}
	if interestedIn == "male" || interestedIn == "men" || interestedIn == "man" {
		return "Male"
	}

	// Fallback based on user's gender
	gender := strings.TrimSpace(strings.ToLower(user.Gender))
	if gender == "male" || gender == "man" || gender == "men" {
		return "Female"
	}
	if gender == "female" || gender == "woman" || gender == "women" {
		return "Male"
	}
	return ""
}

// GetTargetGenderRegexPattern returns a regex pattern matching all variations of the target gender
func GetTargetGenderRegexPattern(user *models.User) string {
	target := GetTargetGender(user)
	if target == "Female" {
		return "^(female|woman|women)"
	}
	if target == "Male" {
		return "^(male|man|men)"
	}
	return ""
}

// FindPotentialMatches fetches users for the discovery feed, respecting user preferences
func (ur *UserRepository) FindPotentialMatches(currentUser *models.User, excludedUserIDs []primitive.ObjectID, limit int64) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Exclude current user, interacted users, and blocked users
	allExcluded := append([]primitive.ObjectID{currentUser.ID}, excludedUserIDs...)
	allExcluded = append(allExcluded, currentUser.BlockedUsers...)
	filter := bson.M{"_id": bson.M{"$nin": allExcluded}}

	// 1b. Apply Gender Filtering (Men see Women, Women see Men)
	targetGenderPattern := GetTargetGenderRegexPattern(currentUser)
	if targetGenderPattern != "" {
		filter["gender"] = primitive.Regex{Pattern: targetGenderPattern, Options: "i"}
	}

	// 2. Apply Preferences if they exist!
	if currentUser.PreferredDenomination != "" && currentUser.PreferredDenomination != "Any" {
		filter["denomination"] = currentUser.PreferredDenomination
	}
	if currentUser.PreferredChurchFreq != "" && currentUser.PreferredChurchFreq != "Any" {
		filter["church_freq"] = currentUser.PreferredChurchFreq
	}

	// 3. Apply Age Preferences
	if currentUser.MinAgePref > 0 || currentUser.MaxAgePref > 0 {
		now := time.Now()
		dobFilter := bson.M{}
		
		// If MaxAgePref is 35, they must be born AFTER (Now - 35 years - 1 year)
		if currentUser.MaxAgePref > 0 {
			earliestDOB := now.AddDate(-(currentUser.MaxAgePref + 1), 0, 0)
			dobFilter["$gte"] = earliestDOB
		}
		
		// If MinAgePref is 25, they must be born BEFORE (Now - 25 years)
		if currentUser.MinAgePref > 0 {
			latestDOB := now.AddDate(-currentUser.MinAgePref, 0, 0)
			dobFilter["$lte"] = latestDOB
		}
		
		if len(dobFilter) > 0 {
			filter["dob"] = dobFilter
		}
	}

	opts := options.Find().SetLimit(limit)

	cursor, err := ur.db.Collection("users").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	// If strict preference filters returned 0 candidates, fallback without denomination/church_freq filters
	if len(users) == 0 && (filter["denomination"] != nil || filter["church_freq"] != nil) {
		delete(filter, "denomination")
		delete(filter, "church_freq")
		fallbackCursor, fallbackErr := ur.db.Collection("users").Find(ctx, filter, opts)
		if fallbackErr == nil {
			_ = fallbackCursor.All(ctx, &users)
			fallbackCursor.Close(ctx)
		}
	}
	
	// Clear password hashes before returning
	for i := range users {
		users[i].PasswordHash = ""
	}

	return users, nil
}

// FindPotentialMatchesVector fetches users using MongoDB Atlas Vector Search based on profile embedding similarity
func (ur *UserRepository) FindPotentialMatchesVector(currentUser *models.User, excludedUserIDs []primitive.ObjectID, limit int64) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Exclude current user, interacted users, and blocked users
	allExcluded := append([]primitive.ObjectID{currentUser.ID}, excludedUserIDs...)
	allExcluded = append(allExcluded, currentUser.BlockedUsers...)

	// Create matches filter
	matchFilter := bson.M{
		"_id": bson.M{"$nin": allExcluded},
	}

	// 1b. Apply Gender Filtering (Men see Women, Women see Men)
	targetGenderPattern := GetTargetGenderRegexPattern(currentUser)
	if targetGenderPattern != "" {
		matchFilter["gender"] = primitive.Regex{Pattern: targetGenderPattern, Options: "i"}
	}

	// 2. Apply Age Preferences if they exist
	if currentUser.MinAgePref > 0 || currentUser.MaxAgePref > 0 {
		now := time.Now()
		dobFilter := bson.M{}
		if currentUser.MaxAgePref > 0 {
			earliestDOB := now.AddDate(-(currentUser.MaxAgePref + 1), 0, 0)
			dobFilter["$gte"] = earliestDOB
		}
		if currentUser.MinAgePref > 0 {
			latestDOB := now.AddDate(-currentUser.MinAgePref, 0, 0)
			dobFilter["$lte"] = latestDOB
		}
		if len(dobFilter) > 0 {
			matchFilter["dob"] = dobFilter
		}
	}

	// Choose query vector: PartnerPrefEmbedding if present, else fallback to ProfileEmbedding
	var queryVector []float32
	if len(currentUser.PartnerPrefEmbedding) == 384 {
		queryVector = currentUser.PartnerPrefEmbedding
	} else {
		queryVector = currentUser.ProfileEmbedding
	}

	// Build MongoDB Aggregation Pipeline
	pipeline := mongo.Pipeline{
		// First Stage: Atlas Vector Search
		{{Key: "$vectorSearch", Value: bson.D{
			{Key: "index", Value: "vector_index"},
			{Key: "path", Value: "profile_embedding"},
			{Key: "queryVector", Value: queryVector},
			{Key: "numCandidates", Value: limit * 5}, // Look at more candidates for filtering
			{Key: "limit", Value: limit},
		}}},
		// Second Stage: Filter excluded users & age range
		{{Key: "$match", Value: matchFilter}},
	}

	cursor, err := ur.db.Collection("users").Aggregate(ctx, pipeline)
	if err != nil {
		logrus.Errorf("❌ Vector search failed: %v. Falling back to standard query.", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	// Rank users by Mutual Compatibility Score: (Score1 + Score2) / 2
	sort.Slice(users, func(i, j int) bool {
		scoreI := calculateMutualScore(currentUser, &users[i])
		scoreJ := calculateMutualScore(currentUser, &users[j])
		return scoreI > scoreJ
	})

	// Clear password hashes before returning
	for i := range users {
		users[i].PasswordHash = ""
	}

	return users, nil
}

// cosineSimilarity calculates the cosine similarity between two float vectors
func cosineSimilarity(v1, v2 []float32) float32 {
	if len(v1) != len(v2) || len(v1) == 0 {
		return 0
	}
	var dotProduct float32
	var normA float32
	var normB float32
	for i := 0; i < len(v1); i++ {
		dotProduct += v1[i] * v2[i]
		normA += v1[i] * v1[i]
		normB += v2[i] * v2[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / float32(math.Sqrt(float64(normA))*math.Sqrt(float64(normB)))
}

// calculateMutualScore calculates the mutual score between current user and a candidate user
func calculateMutualScore(currentUser, candidate *models.User) float32 {
	var score1 float32
	// A's preference vs B's profile
	if len(currentUser.PartnerPrefEmbedding) == 384 {
		score1 = cosineSimilarity(currentUser.PartnerPrefEmbedding, candidate.ProfileEmbedding)
	} else {
		score1 = cosineSimilarity(currentUser.ProfileEmbedding, candidate.ProfileEmbedding)
	}

	// B's preference vs A's profile
	var score2 float32
	if len(candidate.PartnerPrefEmbedding) == 384 {
		score2 = cosineSimilarity(candidate.PartnerPrefEmbedding, currentUser.ProfileEmbedding)
	} else {
		// Fallback: If B hasn't set preferences, assume it matches score1
		score2 = score1
	}

	return (score1 + score2) / 2.0
}


package services

import (
	"github.com/olabanji12-ojo/church-backend/models"
	"github.com/olabanji12-ojo/church-backend/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileService struct {
	userRepository   *repositories.UserRepository
	embeddingService *EmbeddingService
}

func NewProfileService(userRepository *repositories.UserRepository) *ProfileService {
	return &ProfileService{
		userRepository:   userRepository,
		embeddingService: NewEmbeddingService(),
	}
}

// TriggerEmbeddingUpdate generates the user embedding asynchronously in the background.
func (ps *ProfileService) TriggerEmbeddingUpdate(userID primitive.ObjectID) {
	go func() {
		user, err := ps.userRepository.FindUserByID(userID)
		if err != nil {
			return
		}

		if user.IsGuest {
			return
		}

		text := ps.embeddingService.GenerateUserText(user)
		embedding, err := ps.embeddingService.GetEmbedding(text)
		if err != nil {
			return
		}

		_ = ps.userRepository.UpdateUserByID(userID, bson.M{"profile_embedding": embedding})
	}()
}

// TriggerPreferenceEmbeddingUpdate generates the user preference embedding asynchronously in the background.
func (ps *ProfileService) TriggerPreferenceEmbeddingUpdate(userID primitive.ObjectID) {
	go func() {
		user, err := ps.userRepository.FindUserByID(userID)
		if err != nil {
			return
		}

		if user.IsGuest {
			return
		}

		text := ps.embeddingService.GeneratePartnerPreferenceText(user)
		embedding, err := ps.embeddingService.GetEmbedding(text)
		if err != nil {
			return
		}

		_ = ps.userRepository.UpdateUserByID(userID, bson.M{"partner_pref_embedding": embedding})
	}()
}

// UpdateFaithProfile updates the faith-specific attributes of a user
func (ps *ProfileService) UpdateFaithProfile(userID primitive.ObjectID, denomination, churchFreq, prayerFreq, bibleFreq string) error {
	update := bson.M{
		"denomination": denomination,
		"church_freq":  churchFreq,
		"prayer_freq":  prayerFreq,
		"bible_freq":   bibleFreq,
	}

	err := ps.userRepository.UpdateUserByID(userID, update)
	if err == nil {
		ps.TriggerEmbeddingUpdate(userID)
	}
	return err
}

// UpdateIntentions updates what the user is looking for
func (ps *ProfileService) UpdateIntentions(userID primitive.ObjectID, intention, interestedIn string) error {
	update := bson.M{
		"intention":     intention,
		"interested_in": interestedIn,
	}

	err := ps.userRepository.UpdateUserByID(userID, update)
	if err == nil {
		ps.TriggerEmbeddingUpdate(userID)
	}
	return err
}

// UpdatePreferences updates the discovery filters
func (ps *ProfileService) UpdatePreferences(userID primitive.ObjectID, minAge, maxAge, maxDistance int) error {
	update := bson.M{
		"min_age_pref": minAge,
		"max_age_pref": maxAge,
		"max_distance": maxDistance,
	}

	return ps.userRepository.UpdateUserByID(userID, update)
}

// GetUserProfile fetches the complete profile of the user
func (ps *ProfileService) GetUserProfile(userID primitive.ObjectID) (*models.User, error) {
	return ps.userRepository.FindUserByID(userID)
}

// UpdateUserProfile updates generic profile fields
func (ps *ProfileService) UpdateUserProfile(userID primitive.ObjectID, updateData bson.M) error {
	err := ps.userRepository.UpdateUserByID(userID, updateData)
	if err == nil {
		// Only trigger update if we modified fields that impact embeddings
		impactsEmbedding := false
		fields := []string{"first_name", "gender", "denomination", "church_freq", "prayer_freq", "bible_freq", "intention", "bio"}
		for _, field := range fields {
			if _, exists := updateData[field]; exists {
				impactsEmbedding = true
				break
			}
		}
		if impactsEmbedding {
			ps.TriggerEmbeddingUpdate(userID)
		}

		// Trigger preference embedding update if preference fields are modified
		impactsPrefEmbedding := false
		prefFields := []string{"partner_pref_text", "preferred_denomination", "preferred_church_freq", "min_age_pref", "max_age_pref"}
		for _, field := range prefFields {
			if _, exists := updateData[field]; exists {
				impactsPrefEmbedding = true
				break
			}
		}
		if impactsPrefEmbedding {
			ps.TriggerPreferenceEmbeddingUpdate(userID)
		}
	}
	return err
}

// AddPhotos adds Cloudinary URLs to the user's photos array
func (ps *ProfileService) AddPhotos(userID primitive.ObjectID, photoURLs []string) error {
	setUpdate := bson.M{"photos": photoURLs}
	return ps.userRepository.UpdateUserByID(userID, setUpdate)
}


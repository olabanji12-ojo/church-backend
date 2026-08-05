package services

import (
	"github.com/olabanji12-ojo/church-backend/models"
	"github.com/olabanji12-ojo/church-backend/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PrayerService struct {
	prayerRepository *repositories.PrayerRepository
}

func NewPrayerService(prayerRepo *repositories.PrayerRepository) *PrayerService {
	return &PrayerService{prayerRepository: prayerRepo}
}

// PostPrayer saves a new prayer request
func (ps *PrayerService) PostPrayer(authorID primitive.ObjectID, content string) error {
	newPrayer := models.Prayer{
		ID:       primitive.NewObjectID(),
		AuthorID: authorID,
		Content:  content,
	}
	return ps.prayerRepository.CreatePrayer(newPrayer)
}

// GetPrayerFeed fetches prayers for the community feed
func (ps *PrayerService) GetPrayerFeed() ([]models.PrayerFeedItem, error) {
	return ps.prayerRepository.GetRecentPrayers(50) // limit to 50
}

// SayAmen allows a user to "like" a prayer
func (ps *PrayerService) SayAmen(prayerID, userID primitive.ObjectID) error {
	return ps.prayerRepository.AddAmen(prayerID, userID)
}

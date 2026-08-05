package services

import (
	"github.com/olabanji12-ojo/church-backend/models"
	"github.com/olabanji12-ojo/church-backend/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportService struct {
	reportRepository *repositories.ReportRepository
	userRepository   *repositories.UserRepository
}

func NewReportService(reportRepo *repositories.ReportRepository, userRepo *repositories.UserRepository) *ReportService {
	return &ReportService{
		reportRepository: reportRepo,
		userRepository:   userRepo,
	}
}

// SubmitReport handles the business logic when a user reports another user
func (rs *ReportService) SubmitReport(reporterID, reportedID primitive.ObjectID, reason string) error {
	
	newReport := models.Report{
		ID:         primitive.NewObjectID(),
		ReporterID: reporterID,
		ReportedID: reportedID,
		Reason:     reason,
	}

	// Block the user instantly
	err := rs.userRepository.BlockUser(reporterID, reportedID)
	if err != nil {
		return err // Or log and continue
	}

	return rs.reportRepository.CreateReport(newReport)
}

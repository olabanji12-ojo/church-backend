package repositories

import (
	"context"
	"time"

	"github.com/olabanji12-ojo/church-backend/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportRepository struct {
	db *mongo.Database
}

func NewReportRepository(db *mongo.Database) *ReportRepository {
	return &ReportRepository{db: db}
}

// CreateReport logs a new abuse or spam report into MongoDB
func (rr *ReportRepository) CreateReport(report models.Report) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	report.CreatedAt = time.Now()
	report.Status = "pending"

	_, err := rr.db.Collection("reports").InsertOne(ctx, report)
	return err
}

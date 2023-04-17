package survey

import (
	"context"
	"github.com/jmoiron/sqlx"

	"github.com/remotestate/golang/internal"
	"github.com/remotestate/golang/models"
)

type Service struct {
	logger *internal.Logger
	repo   *repository
}

// NewSurveyService creates a new survey service
func NewSurveyService(db *internal.Database, logger *internal.Logger) Service {
	return Service{
		logger: logger,
		repo:   newSurveyRepository(db, logger),
	}
}

func (s Service) FindAttributeWeightForUser(ctx context.Context, trxHandle *sqlx.Tx, surveyInviteID string) ([]models.SurveyAnswerAttributeWeight, error) {
	return s.repo.findAttributeWeightForUser(ctx, trxHandle, surveyInviteID)
}

func (s Service) FindGenderAttribute(ctx context.Context, trxHandle *sqlx.Tx, surveyInviteID string) ([]string, error) {
	return s.repo.findGenderAttribute(ctx, trxHandle, surveyInviteID)
}

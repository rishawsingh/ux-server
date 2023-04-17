package services

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/remotestate/golang/models"
)

type Survey interface {
	FindAttributeWeightForUser(ctx context.Context, trxHandle *sqlx.Tx, surveyInviteID string) ([]models.SurveyAnswerAttributeWeight, error)
	FindGenderAttribute(ctx context.Context, trxHandle *sqlx.Tx, surveyInviteID string) ([]string, error)
}

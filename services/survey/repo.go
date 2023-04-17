package survey

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/remotestate/golang/internal"
	"github.com/remotestate/golang/models"
)

type repository struct {
	*internal.Database
	logger *internal.Logger
}

func newSurveyRepository(db *internal.Database, logger *internal.Logger) *repository {
	return &repository{
		Database: db,
		logger:   logger,
	}
}

func (r *repository) findAttributeWeightForUser(ctx context.Context, tx *sqlx.Tx, surveyInviteID string) ([]models.SurveyAnswerAttributeWeight, error) {
	// language=SQL
	SQL := `SELECT 
    			ua.attribute_id as id,
	        	count(ua.attribute_id) as weight,
	        	a.name as name
            FROM survey_attribute ua
	 		JOIN attribute a on ua.attribute_id = a.id AND a.archived_at IS NULL AND a.is_enabled = TRUE
            WHERE ua.survey_invite_id = $1
            AND ua.archived_at IS NULL
            GROUP BY ua.attribute_id, a.name`
	attributeWeights := make([]models.SurveyAnswerAttributeWeight, 0)
	var err error
	if tx != nil {
		err = tx.SelectContext(ctx, &attributeWeights, SQL, surveyInviteID)
	} else {
		err = r.DB.SelectContext(ctx, &attributeWeights, SQL, surveyInviteID)
	}
	return attributeWeights, errors.WithStack(err)
}

func (r *repository) findGenderAttribute(ctx context.Context, tx *sqlx.Tx, surveyInviteID string) ([]string, error) {
	// language=SQL
	SQL := `SELECT 
    			a.id
			FROM survey_invite s
			JOIN gender_attribute_v2 g ON s.gender = g.gender_id AND g.archived_at IS NULL
			JOIN attribute a on g.attribute_id = a.id AND a.archived_at IS NULL AND a.is_enabled IS TRUE
			WHERE s.id = $1
			AND s.archived_at IS NULL`
	attributes := make([]string, 0)
	var err error
	if tx != nil {
		err = tx.SelectContext(ctx, &attributes, SQL, surveyInviteID)
	} else {
		err = r.DB.SelectContext(ctx, &attributes, SQL, surveyInviteID)
	}
	return attributes, err
}

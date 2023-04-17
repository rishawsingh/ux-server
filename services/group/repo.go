package group

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/remotestate/golang/internal"
	"github.com/remotestate/golang/models"
)

type repository struct {
	*internal.Database
	logger *internal.Logger
}

func newGroupRepository(db *internal.Database, logger *internal.Logger) *repository {
	return &repository{
		Database: db,
		logger:   logger,
	}
}

func (r *repository) CreateGroup(groupDetails models.GroupDetails, userID string, tx *sqlx.Tx) (string, error) {
	// language=SQL
	SQL := `INSERT INTO groups(name, description, visibility, is_verified, location_id, created_by, group_image_id)
                       VALUES ($1, $2, $3, $4, $5, $6, $7)
           RETURNING id`
	var groupID string
	var err error
	if tx != nil {
		err = tx.Get(&groupID, SQL, groupDetails.Name, groupDetails.Description, groupDetails.Visibility, groupDetails.Visibility, groupDetails.IsVerified, "locationID", userID, "groupID")
	} else {
		err = r.DB.SelectContext(nil, &groupID, SQL)
	}
	return groupID, errors.WithStack(err)
}

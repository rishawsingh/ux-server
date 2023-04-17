package location

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

func newLocationRepository(db *internal.Database, logger *internal.Logger) *repository {
	return &repository{
		Database: db,
		logger:   logger,
	}
}

func (r *repository) CreateLocation(groupDetails models.GroupDetails, userID string, tx *sqlx.Tx) (string, error) {
	// language=SQL
	SQL := `INSERT INTO location(name)
                        VALUES ($1)
            RETURNING id`
	var locationID string
	var err error
	if tx != nil {
		err = tx.Get(&locationID, SQL, "group image", groupDetails.UploadType, groupDetails.Url, userID)
	} else {
		err = r.DB.SelectContext(nil, &locationID, SQL)
	}
	return locationID, errors.WithStack(err)
}

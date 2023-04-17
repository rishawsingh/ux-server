package attachment

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

func newAttachmentRepository(db *internal.Database, logger *internal.Logger) *repository {
	return &repository{
		Database: db,
		logger:   logger,
	}
}

func (r *repository) CreateAttachment(groupDetails models.GroupDetails, userID string, tx *sqlx.Tx) (string, error) {
	// language=SQL
	SQL := `INSERT INTO uploads(name, type, url, uploaded_by)
                        VALUES ($1, $2, $3, $4)
            RETURNING id`
	var attachmentID string
	var err error
	if tx != nil {
		err = tx.Get(&attachmentID, SQL, "group image", groupDetails.UploadType, groupDetails.Url, userID)
	} else {
		err = r.DB.SelectContext(nil, &attachmentID, SQL)
	}
	return attachmentID, errors.WithStack(err)
}

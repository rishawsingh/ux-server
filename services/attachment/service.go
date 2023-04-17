package attachment

import (
	"github.com/jmoiron/sqlx"
	"github.com/remotestate/golang/internal"
	"github.com/remotestate/golang/models"
)

type Service struct {
	logger *internal.Logger
	repo   *repository
}

// NewAttachmentService creates a new attachment service
func NewAttachmentService(db *internal.Database, logger *internal.Logger) Service {
	return Service{
		logger: logger,
		repo:   newAttachmentRepository(db, logger),
	}
}

func (s Service) CreateAttachment(groupDetails models.GroupDetails, userID string, trxHandle *sqlx.Tx) (string, error) {
	return s.repo.CreateAttachment(groupDetails, userID, trxHandle)
}

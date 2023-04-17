package location

import (
	"github.com/jmoiron/sqlx"
	"github.com/remotestate/golang/internal"
	"github.com/remotestate/golang/models"
)

type Service struct {
	logger *internal.Logger
	repo   *repository
}

// NewLocationService creates a new group service
func NewLocationService(db *internal.Database, logger *internal.Logger) Service {
	return Service{
		logger: logger,
		repo:   newLocationRepository(db, logger),
	}
}

func (s Service) CreateLocation(groupDetails models.GroupDetails, userID string, trxHandle *sqlx.Tx) (string, error) {
	return s.repo.CreateLocation(groupDetails, userID, trxHandle)
}

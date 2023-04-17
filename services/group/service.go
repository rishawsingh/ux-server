package group

import (
	"github.com/jmoiron/sqlx"
	"github.com/remotestate/golang/internal"
	"github.com/remotestate/golang/models"
)

type Service struct {
	logger *internal.Logger
	repo   *repository
}

// NewGroupService creates a new group service
func NewGroupService(db *internal.Database, logger *internal.Logger) Service {
	return Service{
		logger: logger,
		repo:   newGroupRepository(db, logger),
	}
}

func (s Service) CreateGroup(groupDetails models.GroupDetails, userID string, trxHandle *sqlx.Tx) (string, error) {
	return s.repo.CreateGroup(groupDetails, userID, trxHandle)
}

package services

import (
	"github.com/jmoiron/sqlx"
	"github.com/remotestate/golang/models"
)

type Location interface {
	CreateLocation(groupDetails models.GroupDetails, userID string, trxHandle *sqlx.Tx) (string, error)
}

package services

import (
	"github.com/jmoiron/sqlx"
	"github.com/remotestate/golang/models"
)

type Group interface {
	CreateGroup(groupDetails models.GroupDetails, userID string, trxHandle *sqlx.Tx) (string, error)
}

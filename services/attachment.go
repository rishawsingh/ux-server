package services

import (
	"github.com/jmoiron/sqlx"
	"github.com/remotestate/golang/models"
)

type Attachment interface {
	CreateAttachment(groupDetails models.GroupDetails, userID string, trxHandle *sqlx.Tx) (string, error)
}

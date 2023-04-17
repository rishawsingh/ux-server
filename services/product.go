package services

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/remotestate/golang/models"
)

type Product interface {
	GetAllProductWithAttributes(ctx context.Context, trxHandle *sqlx.Tx) ([]models.ProductAttributeList, error)
}

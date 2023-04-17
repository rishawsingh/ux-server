package product

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/remotestate/golang/internal"
	"github.com/remotestate/golang/models"
)

type Service struct {
	logger *internal.Logger
	repo   *repository
}

// NewProductService creates a new product service
func NewProductService(db *internal.Database, logger *internal.Logger) Service {
	return Service{
		logger: logger,
		repo:   newProductRepository(db, logger),
	}
}

func (s Service) GetAllProductWithAttributes(ctx context.Context, trxHandle *sqlx.Tx) ([]models.ProductAttributeList, error) {
	return s.repo.getAllProductWithAttributes(ctx, trxHandle)
}

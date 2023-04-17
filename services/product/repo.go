package product

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/remotestate/golang/internal"
	"github.com/remotestate/golang/models"
)

type repository struct {
	*internal.Database
	logger *internal.Logger
}

func newProductRepository(db *internal.Database, logger *internal.Logger) *repository {
	return &repository{
		Database: db,
		logger:   logger,
	}
}

func (r *repository) getAllProductWithAttributes(ctx context.Context, tx *sqlx.Tx) ([]models.ProductAttributeList, error) {
	// language=SQL
	SQL := `SELECT
				p.id,
				sc.id as subcategory_id,
				b.id as brand_id,
                ARRAY_AGG(DISTINCT pa.attribute_id)
                    || 
                ARRAY_AGG(DISTINCT ba.attribute_id) FILTER ( WHERE ba.attribute_id IS NOT NULL ) 
                AS attribute_id
            FROM product p
            JOIN product_attribute pa ON p.id = pa.product_id AND pa.archived_at IS NULL
            JOIN attribute a ON pa.attribute_id = a.id AND a.archived_at IS NULL AND a.is_enabled IS TRUE
            JOIN brand_v2 b ON b.id = p.brand_id AND b.archived_at IS NULL AND b.is_enabled IS TRUE
            LEFT JOIN brand_attribute_v2 ba ON ba.brand_id = b.id AND ba.archived_at IS NULL
            LEFT JOIN attribute baa ON baa.id = ba.attribute_id AND baa.archived_at IS NULL AND baa.is_enabled IS TRUE
			JOIN categories c ON c.id = p.category_id  AND c.archived_at IS NULL AND c.is_enabled IS TRUE
			JOIN sub_category sc ON p.sub_category_id = sc.id AND sc.archived_at IS NULL AND sc.is_enabled IS TRUE
			JOIN product_type pt ON p.product_type_id = pt.id AND pt.archived_at IS NULL AND pt.is_enabled IS TRUE
            WHERE p.archived_at IS NULL
            AND p.is_enabled IS TRUE
            AND p.in_stock IS TRUE
            GROUP BY p.id, sc.id, b.id`
	productAttribute := make([]models.ProductAttributeList, 0)
	var err error
	if tx != nil {
		err = tx.SelectContext(ctx, &productAttribute, SQL)
	} else {
		err = r.DB.SelectContext(ctx, &productAttribute, SQL)
	}
	return productAttribute, errors.WithStack(err)
}

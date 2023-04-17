package internal

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Transactor struct {
	db     *Database
	logger *Logger
}

func NewTransactor(db *Database, logger *Logger) *Transactor {
	return &Transactor{db: db, logger: logger}
}

func (t *Transactor) Wrap(ctx context.Context, fn func(tx *sqlx.Tx) error) (err error) {
	tx, txErr := t.db.Beginx()
	if txErr != nil {
		return errors.WithMessage(txErr, "failed to start a transaction")
	}
	defer func(err error, tx *sqlx.Tx, t *Transactor) {
		if panicErr := recover(); panicErr != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = errors.Wrap(err, fmt.Sprintf("panic and rollback error in transaction: %v", panicErr))
			} else {
				err = errors.Errorf("panic error in transaction: %v", panicErr)
			}
		}
		if err != nil {
			if rollBackErr := tx.Rollback(); rollBackErr != nil {
				t.logger.WithCtx(ctx).With(rollBackErr).Error("failed to rollback database transaction")
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			t.logger.WithCtx(ctx).With(commitErr).Error("failed to commit database transaction")
		}
	}(err, tx, t)
	err = fn(tx)
	return err
}

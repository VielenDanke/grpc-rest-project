package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CompanyRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewCompanyRepository(pool *pgxpool.Pool) CompanyRepository {
	return &CompanyRepositoryImpl{pool: pool}
}

func (c CompanyRepositoryImpl) GetByBin(ctx context.Context, bin string) (string, error) {
	tx, txErr := c.pool.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})

	if txErr != nil {
		return "", txErr
	}
	var name string
	if rowErr := tx.QueryRow(ctx, "select c.name from companies c where c.bin = $1", bin).Scan(&name); rowErr != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return "", rbErr
		}
		return "", rowErr
	}
	if commErr := tx.Commit(ctx); commErr != nil {
		return "", commErr
	}
	return name, nil
}

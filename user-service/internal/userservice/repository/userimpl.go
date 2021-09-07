package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	u "github.com/vielendanke/grpc-rest-project/user-service/user"
	"log"
)

type UserRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &UserRepositoryImpl{pool: pool}
}

func (u UserRepositoryImpl) SaveUser(ctx context.Context, sr *u.SaveUserRequest) (string, error) {
	log.Println(fmt.Sprintf("User saved: %s. Company INN %s", sr.FullName, sr.CompanyBin))

	tx, txErr := u.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if txErr != nil {
		return "", txErr
	}
	var resultId string

	if scErr := tx.QueryRow(ctx,
		"insert into users(iin, fullname, company_bin) values($1, $2, $3) returning id",
		sr.GetIin(), sr.GetFullName(), sr.GetCompanyBin(),
	).Scan(&resultId); scErr != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return "", rbErr
		}
		return "", scErr
	}
	if cErr := tx.Commit(ctx); cErr != nil {
		return "", cErr
	}
	return resultId, nil
}

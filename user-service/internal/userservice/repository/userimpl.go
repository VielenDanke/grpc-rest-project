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

func (ur UserRepositoryImpl) FindAll(ctx context.Context) ([]*u.UserResponse, error) {
	tx, txErr := ur.pool.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if txErr != nil {
		return nil, txErr
	}
	rows, rowsErr := tx.Query(ctx, "select u.id, u.iin, u.fullname, u.company_bin from users u")

	if rowsErr != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return nil, rbErr
		}
		return nil, rowsErr
	}
	users := make([]*u.UserResponse, 0)
	for rows.Next() {
		uResp := &u.UserResponse{}
		if scErr := rows.Scan(&uResp.Id, &uResp.Iin, &uResp.FullName, &uResp.CompanyBin); scErr != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				return nil, rbErr
			}
			return nil, scErr
		}
		users = append(users, uResp)
	}
	if cErr := tx.Commit(ctx); cErr != nil {
		return nil, cErr
	}
	return users, nil
}

func (ur UserRepositoryImpl) SaveUser(ctx context.Context, sr *u.SaveUserRequest) (string, error) {
	log.Println(fmt.Sprintf("User saved: %s. Company INN %s", sr.FullName, sr.CompanyBin))

	tx, txErr := ur.pool.BeginTx(ctx, pgx.TxOptions{
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

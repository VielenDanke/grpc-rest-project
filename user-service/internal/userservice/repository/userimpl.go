package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	u "github.com/vielendanke/grpc-rest-project/user-service/user"
	"log"
)

type UserRepositoryImpl struct {

}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (u UserRepositoryImpl) SaveUser(ctx context.Context, sr *u.SaveUserRequest) (string, error) {
	log.Println(fmt.Sprintf("User saved: %s. Company IIN %s. Company full name %s", sr.FullName, sr.CompanyIin, sr.CompanyFullName))
	return uuid.NewString(), nil
}
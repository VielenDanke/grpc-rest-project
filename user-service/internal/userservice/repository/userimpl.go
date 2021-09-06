package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/vielendanke/grpc-rest-project/user-service/proto"
	"log"
)

type UserRepositoryImpl struct {

}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (u UserRepositoryImpl) SaveUser(ctx context.Context, sr *pb.SaveUserRequest) (string, error) {
	log.Println(fmt.Sprintf("User saved: %s. Company IIN %s. Company full name %s", sr.FullName, sr.CompanyIin, sr.CompanyFullName))
	return uuid.NewString(), nil
}

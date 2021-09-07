package service

import (
	"context"
	"fmt"
	cp "github.com/vielendanke/grpc-rest-project/user-service/company"
	"github.com/vielendanke/grpc-rest-project/user-service/internal/userservice/repository"
	u "github.com/vielendanke/grpc-rest-project/user-service/user"
	"log"
)

type UserServiceImpl struct {
	ur repository.UserRepository
	cs cp.CompanyServiceClient
}

func NewUserService(ur repository.UserRepository, cs cp.CompanyServiceClient) UserService {
	return &UserServiceImpl{ur: ur, cs: cs}
}

func (u UserServiceImpl) SaveUser(ctx context.Context, sr *u.SaveUserRequest) (string, error) {
	resp, respErr := u.cs.CompanyByBin(ctx, &cp.CompanyByBinRequest{Bin: sr.CompanyBin})

	log.Println(fmt.Sprintf("Name received %s", resp.GetName()))

	if respErr != nil {
		return "", respErr
	}
	return u.ur.SaveUser(ctx, sr)
}

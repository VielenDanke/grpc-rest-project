package service

import (
	"context"
	cp "github.com/vielendanke/grpc-rest-project/user-service/company"
	"github.com/vielendanke/grpc-rest-project/user-service/internal/userservice/repository"
	u "github.com/vielendanke/grpc-rest-project/user-service/user"
)

type UserServiceImpl struct {
	ur repository.UserRepository
	cs cp.CompanyServiceClient
}

func NewUserService(ur repository.UserRepository, cs cp.CompanyServiceClient) UserService {
	return &UserServiceImpl{ur: ur, cs: cs}
}

func (u UserServiceImpl) SaveUser(ctx context.Context, sr *u.SaveUserRequest) (string, error) {
	iin, respErr := u.cs.CompanyByIIN(ctx, &cp.CompanyByIINRequest{Inn: sr.CompanyIin})

	if respErr != nil {
		return "", respErr
	}
	sr.CompanyFullName = iin.GetFullName()
	return u.ur.SaveUser(ctx, sr)
}
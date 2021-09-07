package handler

import (
	"context"
	"github.com/vielendanke/grpc-rest-project/user-service/internal/userservice/service"
	u "github.com/vielendanke/grpc-rest-project/user-service/user"
)

type User struct {
	srv service.UserService
	u.UserServer
}

func NewUserHandler(srv service.UserService) *User {
	return &User{srv: srv}
}

func (us *User) SaveUser(ctx context.Context, request *u.SaveUserRequest) (*u.SaveUserResponse, error) {
	id, err := us.srv.SaveUser(ctx, request)
	if err != nil {
		return nil, err
	}
	return &u.SaveUserResponse{Id: id}, nil
}

package handler

import (
	"context"
	"github.com/vielendanke/grpc-rest-project/user-service/internal/userservice/service"
	pb "github.com/vielendanke/grpc-rest-project/user-service/proto"
)

type User struct {
	srv service.UserService
	pb.UserServer
}

func NewUserHandler(srv service.UserService) *User {
	return &User{srv: srv}
}

func (u User) SaveUser(ctx context.Context, request *pb.SaveUserRequest) (*pb.SaveUserResponse, error) {
	id, err := u.srv.SaveUser(ctx, request)
	if err != nil {
		return nil, err
	}
	return &pb.SaveUserResponse{Id: id}, nil
}

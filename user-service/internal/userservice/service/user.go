package service

import (
	"context"
	u "github.com/vielendanke/grpc-rest-project/user-service/user"
)

type UserService interface {
	SaveUser(ctx context.Context, sr *u.SaveUserRequest) (string, error)
	FindAll(ctx context.Context) ([]*u.UserResponse, error)
}

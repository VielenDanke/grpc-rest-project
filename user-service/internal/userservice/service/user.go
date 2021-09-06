package service

import (
	"context"
	pb "github.com/vielendanke/grpc-rest-project/user-service/proto"
)

type UserService interface {
	SaveUser(ctx context.Context, sr *pb.SaveUserRequest) (string, error)
}

package repository

import (
	"context"
	pb "github.com/vielendanke/grpc-rest-project/user-service/proto"
)

type UserRepository interface {
	SaveUser(ctx context.Context, sr *pb.SaveUserRequest) (string, error)
}

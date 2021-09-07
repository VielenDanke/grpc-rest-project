package service

import (
	"context"

	pb "github.com/vielendanke/grpc-rest-project/company-service/proto"
)

type TaskService interface {
	CompanyByBin(context.Context, *pb.CompanyByBinRequest) (*pb.CompanyByBinResponse, error)
}

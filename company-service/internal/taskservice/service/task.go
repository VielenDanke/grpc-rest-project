package service

import (
	"context"

	pb "github.com/vielendanke/grpc-rest-project/company-service/proto"
)

type TaskService interface {
	CompanyByIIN(context.Context, *pb.CompanyByIINRequest) (*pb.CompanyByIINResponse, error)
}

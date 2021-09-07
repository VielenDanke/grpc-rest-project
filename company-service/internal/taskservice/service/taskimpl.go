package service

import (
	"context"
	"fmt"
	pb "github.com/vielendanke/grpc-rest-project/company-service/proto"
	"log"
)

type TaskServiceImpl struct {
}

func NewTaskService() TaskService {
	return &TaskServiceImpl{}
}

func (ts *TaskServiceImpl) CompanyByIIN(ctx context.Context, req *pb.CompanyByIINRequest) (*pb.CompanyByIINResponse, error) {
	log.Println(fmt.Sprintf("Company by IIN %s is requested", req.GetInn()))
	return &pb.CompanyByIINResponse{Inn: req.Inn, Kpp: "Kpp", Name: "Test", FullName: "Test"}, nil
}

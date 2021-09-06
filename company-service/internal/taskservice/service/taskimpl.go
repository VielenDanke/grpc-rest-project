package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	pb "github.com/vielendanke/grpc-rest-project/company-service/proto"
)

type TaskServiceImpl struct {
	cli *http.Client
}

func NewTaskService(cli *http.Client) TaskService {
	return &TaskServiceImpl{cli: cli}
}

func (ts *TaskServiceImpl) CompanyByIIN(ctx context.Context, req *pb.CompanyByIINRequest) (*pb.CompanyByIINResponse, error) {
	log.Println(fmt.Sprintf("Company by IIN %s is requested", req.GetInn()))
	return &pb.CompanyByIINResponse{Inn: req.Inn}, nil
}

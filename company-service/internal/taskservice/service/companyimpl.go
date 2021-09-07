package service

import (
	"context"
	"fmt"
	"github.com/vielendanke/grpc-rest-project/company-service/internal/taskservice/repository"
	pb "github.com/vielendanke/grpc-rest-project/company-service/proto"
	"log"
)

type TaskServiceImpl struct {
	r repository.CompanyRepository
}

func NewTaskService(r repository.CompanyRepository) TaskService {
	return &TaskServiceImpl{r: r}
}

func (ts *TaskServiceImpl) CompanyByBin(ctx context.Context, req *pb.CompanyByBinRequest) (*pb.CompanyByBinResponse, error) {
	log.Println(fmt.Sprintf("Company by Bin %s is requested", req.GetBin()))
	name, err := ts.r.GetByBin(ctx, req.GetBin())
	if err != nil {
		return nil, err
	}
	return &pb.CompanyByBinResponse{Name: name}, nil
}

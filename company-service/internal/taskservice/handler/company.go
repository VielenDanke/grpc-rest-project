package handler

import (
	"context"
	"github.com/vielendanke/grpc-rest-project/company-service/internal/taskservice/service"

	pb "github.com/vielendanke/grpc-rest-project/company-service/proto"
)

type TaskHandler struct {
	ts service.TaskService
	pb.CompanyServiceServer
}

func (th *TaskHandler) CompanyByBin(ctx context.Context, request *pb.CompanyByBinRequest) (*pb.CompanyByBinResponse, error) {
	return th.ts.CompanyByBin(ctx, request)
}

func NewTaskHandler(ts service.TaskService) *TaskHandler {
	return &TaskHandler{ts: ts}
}

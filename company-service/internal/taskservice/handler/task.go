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

func NewTaskHandler(ts service.TaskService) *TaskHandler {
	return &TaskHandler{ts: ts}
}

func (th *TaskHandler) CompanyByIIN(ctx context.Context, req *pb.CompanyByIINRequest) (*pb.CompanyByIINResponse, error) {
	return th.ts.CompanyByIIN(ctx, req)
}

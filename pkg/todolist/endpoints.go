package todolist

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateTaskEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(svc Service) Endpoints {
	return Endpoints{
		CreateTaskEndpoint: MakeCreateTaskEndpoint(svc),
	}
}

type CreateTaskRequest struct {
	S string `json:"s"`
}

type CreateTaskResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"`
}

func MakeCreateTaskEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateTaskRequest)
		v, err := svc.CreateTask(ctx, req.S)
		if err != nil {
			return CreateTaskResponse{v, err.Error()}, nil
		}
		return CreateTaskResponse{v, ""}, nil
	}
}

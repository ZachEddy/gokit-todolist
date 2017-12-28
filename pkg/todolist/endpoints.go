package todolist

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateTaskEndpoint endpoint.Endpoint
	GetTaskEndpoint    endpoint.Endpoint
	ListTasksEndpoint  endpoint.Endpoint
}

func MakeServerEndpoints(svc Service) Endpoints {
	return Endpoints{
		CreateTaskEndpoint: MakeCreateTaskEndpoint(svc),
		GetTaskEndpoint:    MakeGetTaskEndpoint(svc),
		ListTasksEndpoint:  MakeListTasksEndpoint(svc),
	}
}

type CreateTaskRequest struct {
	payload TaskPayload
}

type CreateTaskResponse struct {
	T   *Task  `json:"task,omitempty"`
	Err string `json:"error,omitempty"`
}

type ListTasksRequest struct {
}

type ListTasksResponse struct {
	T   *[]Task `json:"task_list,omitempty"`
	Err string  `json:"error,omitempty"`
}

type GetTaskRequest struct {
	ID uint
}

type GetTaskResponse struct {
	T   *Task  `json:"task,omitempty"`
	Err string `json:"error,omitempty"`
}

func MakeCreateTaskEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateTaskRequest)
		task, err := svc.CreateTask(ctx, req.payload)
		if err != nil {
			return CreateTaskResponse{task, err.Error()}, nil
		}
		return CreateTaskResponse{task, ""}, nil
	}
}

func MakeListTasksEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// req := request.(ListTasksRequest)
		tasks, err := svc.ListTasks(ctx)
		if err != nil {
			return ListTasksResponse{tasks, err.Error()}, nil
		}
		return ListTasksResponse{tasks, ""}, nil
	}
}

func MakeGetTaskEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetTaskRequest)
		task, err := svc.GetTask(ctx, req.ID)
		if err != nil {
			return GetTaskResponse{task, err.Error()}, nil
		}
		return GetTaskResponse{task, ""}, nil
	}
}

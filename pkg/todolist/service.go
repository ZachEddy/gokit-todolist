package todolist

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type Service interface {
	CreateTask(context.Context, string) (string, error)
}

type TodoListService struct{}

func (svc TodoListService) CreateTask(ctx context.Context, str string) (string, error) {
	fmt.Println("Hello world!")
	return "nothing to see here", nil
}

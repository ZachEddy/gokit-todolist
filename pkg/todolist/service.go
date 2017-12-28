package todolist

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type Service interface {
	CreateTask(context.Context, TaskPayload) (*Task, error)
	UpdateTask(context.Context, string) (string, error)
	DeleteTask(context.Context, string) (string, error)
	ListTasks(context.Context) (*[]Task, error)
	GetTask(context.Context, uint) (*Task, error)
}

type Model struct {
	ID        uint `json:"id" gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Task struct {
	gorm.Model
	Name string `json:"name"`
	Body string `json:"body"`
}

type TaskPayload struct {
	Name *string `json:"name"`
	Body *string `json:"body"`
}

type ModifyTaskPayload struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

type TodoListService struct {
	DB *gorm.DB
}

func (svc TodoListService) CreateTask(ctx context.Context, payload TaskPayload) (*Task, error) {
	if payload.Name == nil || *payload.Name == "" {
		return nil, errors.New("request body parameter \"name\" is required")
	}
	if payload.Body == nil || *payload.Body == "" {
		return nil, errors.New("request body parameter \"body\" is required")
	}
	task := Task{Name: *payload.Name, Body: *payload.Body}
	result := svc.DB.Create(&task)
	if err := result.Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (svc TodoListService) UpdateTask(ctx context.Context, str string) (string, error) {
	fmt.Println("update task not implemented!")
	return "nothing to see here", nil
}

func (svc TodoListService) DeleteTask(ctx context.Context, str string) (string, error) {
	fmt.Println("delete task not implemented!")
	return "nothing to see here", nil
}

func (svc TodoListService) ListTasks(ctx context.Context) (*[]Task, error) {
	tasks := make([]Task, 0)
	query := svc.DB.Find(&tasks)
	if errors := query.GetErrors(); len(errors) > 0 {
		// return first error
		return nil, errors[0]
	}
	return &tasks, nil
}

func (svc TodoListService) GetTask(ctx context.Context, id uint) (*Task, error) {
	task := Task{}
	query := svc.DB.First(&task, id)
	if errors := query.GetErrors(); len(errors) > 0 {
		// return first error
		return nil, errors[0]
	}
	return &task, nil
}

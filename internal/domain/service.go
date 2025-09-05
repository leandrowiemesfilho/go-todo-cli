package domain

import (
	"context"

	"github.com/google/uuid"
)

type TodoService interface {
	FindAllTodos(ctx context.Context) ([]*Todo, error)
	FindTodoByID(ctx context.Context, id uuid.UUID) (*Todo, error)
	CreateTodo(ctx context.Context, request CreateTodoRequest) (*Todo, error)
	UpdateTodo(ctx context.Context, request UpdateTodoRequest) (*Todo, error)
	DeleteTodo(ctx context.Context, id uuid.UUID) error
	ToggleTodo(ctx context.Context, id uuid.UUID) (*Todo, error)
}

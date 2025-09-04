package domain

import (
	"context"

	"github.com/google/uuid"
)

type TodoRepository interface {
	FindAll(ctx context.Context) ([]*Todo, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Todo, error)
	Create(ctx context.Context, todo *Todo) error
	Update(ctx context.Context, todo *Todo) error
	Delete(ctx context.Context, id uuid.UUID) error
}

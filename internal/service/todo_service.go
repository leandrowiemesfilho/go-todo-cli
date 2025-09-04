package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/leandrowiemesfilho/go-todo-cli/internal/domain"
)

type TodoService interface {
	FindAllTodos(ctx context.Context) ([]*domain.Todo, error)
	FindTodoByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error)
	CreateTodo(ctx context.Context, request domain.CreateTodoRequest) (*domain.Todo, error)
	UpdateTodo(ctx context.Context, request domain.UpdateTodoRequest) (*domain.Todo, error)
	DeleteTodo(ctx context.Context, id uuid.UUID) error
	ToggleTodo(ctx context.Context, id uuid.UUID) (*domain.Todo, error)
}

type todoServiceImpl struct {
	repo domain.TodoRepository
}

func NewTodoService(repo domain.TodoRepository) TodoService {
	return &todoServiceImpl{repo: repo}
}

func (s todoServiceImpl) FindAllTodos(ctx context.Context) ([]*domain.Todo, error) {
	return s.repo.FindAll(ctx)
}

func (s todoServiceImpl) FindTodoByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error) {
	return s.repo.FindByID(ctx, id)
}

func (s todoServiceImpl) CreateTodo(ctx context.Context, request domain.CreateTodoRequest) (*domain.Todo, error) {
	todo := &domain.Todo{
		ID:          uuid.New(),
		Title:       request.Title,
		Description: request.Description,
		Completed:   false,
		CreatedDate: time.Now(),
		UpdatedDate: time.Now(),
	}

	if err := s.repo.Create(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s todoServiceImpl) UpdateTodo(ctx context.Context, request domain.UpdateTodoRequest) (*domain.Todo, error) {
	todo, err := s.repo.FindByID(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	todo.Title = request.Title
	todo.Description = request.Description
	todo.UpdatedDate = time.Now()

	if err = s.repo.Update(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s todoServiceImpl) DeleteTodo(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s todoServiceImpl) ToggleTodo(ctx context.Context, id uuid.UUID) (*domain.Todo, error) {
	todo, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	todo.Completed = !todo.Completed

	if err = s.repo.Update(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

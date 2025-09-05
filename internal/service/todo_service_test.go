package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/leandrowiemesfilho/go-todo-cli/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTodoRepository struct {
	mock.Mock
}

func (mock *MockTodoRepository) FindAll(ctx context.Context) ([]*domain.Todo, error) {
	args := mock.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.Todo), args.Error(1)
}

func (mock *MockTodoRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error) {
	args := mock.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Todo), args.Error(1)
}

func (mock *MockTodoRepository) Create(ctx context.Context, todo *domain.Todo) error {
	args := mock.Called(ctx, todo)
	return args.Error(0)
}

func (mock *MockTodoRepository) Update(ctx context.Context, todo *domain.Todo) error {
	args := mock.Called(ctx, todo)
	return args.Error(0)
}

func (mock *MockTodoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := mock.Called(ctx, id)
	return args.Error(0)
}

func TestTodoService_FindAllTodos(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	service := NewTodoService(mockRepo)
	ctx := context.Background()

	expectedTodos := []*domain.Todo{
		{
			ID:        uuid.New(),
			Title:     "Test 1",
			Completed: false,
		},
		{
			ID:        uuid.New(),
			Title:     "Test 2",
			Completed: true,
		},
	}

	mockRepo.On("FindAll").Return(expectedTodos, nil)

	result, err := service.FindAllTodos(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, len(expectedTodos))

	mockRepo.AssertExpectations(t)
}

func TestTodoService_FindTodoByID(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	service := NewTodoService(mockRepo)
	ctx := context.Background()

	testID := uuid.New()
	expectedResul := &domain.Todo{
		ID:        testID,
		Title:     "Test 1",
		Completed: false,
	}

	mockRepo.On("FindByID", ctx, testID).Return(expectedResul, nil)

	result, err := service.FindTodoByID(ctx, testID)
	assert.NoError(t, err)
	assert.Equal(t, expectedResul.ID, result.ID)

	mockRepo.AssertExpectations(t)
}

func TestTodoService_FindTodoByID_NotFound(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	service := NewTodoService(mockRepo)
	ctx := context.Background()

	testID := uuid.New()
	mockRepo.On("FindByID", ctx, testID).Return(nil, errors.New("not found"))

	todo, err := service.FindTodoByID(ctx, testID)
	assert.Error(t, err)
	assert.Nil(t, todo)

	mockRepo.AssertExpectations(t)
}

func TestTodoService_CreateTodo(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	service := NewTodoService(mockRepo)
	ctx := context.Background()

	request := domain.CreateTodoRequest{
		Title:       "Test 1",
		Description: "Description 1",
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Todo")).
		Return(nil).
		Run(func(args mock.Arguments) {
			todo := args.Get(1).(*domain.Todo)

			assert.NotNil(t, todo.ID)
			assert.Equal(t, request.Title, todo.Title)
			assert.Equal(t, request.Description, todo.Description)
			assert.False(t, todo.Completed)
		})

	result, err := service.CreateTodo(ctx, request)

	assert.NoError(t, err)
	assert.Equal(t, request.Title, result.Title)
	assert.Equal(t, request.Description, result.Description)
	assert.False(t, result.Completed)

	mockRepo.AssertExpectations(t)
}

func TestTodoService_UpdateTodo(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	service := NewTodoService(mockRepo)
	ctx := context.Background()

	testID := uuid.New()
	existingTodo := &domain.Todo{
		ID:          testID,
		Title:       "Test 1",
		Description: "Description 1",
		Completed:   false,
	}
	request := domain.UpdateTodoRequest{
		ID:          testID,
		Title:       "Test Updated",
		Description: "Description Updated",
	}

	mockRepo.On("FindByID", ctx, testID).Return(existingTodo, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Todo")).
		Return(nil).
		Run(func(args mock.Arguments) {
			todo := args.Get(1).(*domain.Todo)

			assert.Equal(t, request.Title, todo.Title)
			assert.Equal(t, request.Description, todo.Description)
		})

	result, err := service.UpdateTodo(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, request.Title, result.Title)
	assert.Equal(t, request.Description, result.Description)

	mockRepo.AssertExpectations(t)
}

func TestTodoService_DeleteTodo(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	service := NewTodoService(mockRepo)
	ctx := context.Background()

	testID := uuid.New()

	mockRepo.On("Delete", ctx, testID).Return(nil)

	err := service.DeleteTodo(ctx, testID)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestTodoServiceImpl_ToggleTodo(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	service := NewTodoService(mockRepo)
	ctx := context.Background()

	testID := uuid.New()
	existingTodo := &domain.Todo{
		ID:          testID,
		Title:       "Test 1",
		Description: "Description 1",
		Completed:   false,
	}

	mockRepo.On("FindByID", ctx, testID).Return(existingTodo, nil)
	mockRepo.On("Update", ctx, existingTodo).
		Return(nil).
		Run(func(args mock.Arguments) {
			todo := args.Get(1).(*domain.Todo)

			assert.Equal(t, todo.ID, existingTodo.ID)
			assert.Equal(t, todo.Completed, existingTodo.Completed)
		})

	result, err := service.ToggleTodo(ctx, testID)
	assert.NoError(t, err)
	assert.Equal(t, result.ID, existingTodo.ID)
	assert.Equal(t, result.Completed, existingTodo.Completed)

	mockRepo.AssertExpectations(t)
}

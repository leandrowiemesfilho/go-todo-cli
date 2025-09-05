package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leandrowiemesfilho/go-todo-cli/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TodoRepositoryTestSuite struct {
	suite.Suite
	pool     *pgxpool.Pool
	repo     *TodoRepository
	ctx      context.Context
	testTodo *domain.Todo
}

func (suite *TodoRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	// Create a test database connection
	// In real tests, you'd use a test container or mock
	config, err := pgxpool.ParseConfig("postgres://todo_user:todo_password@localhost:5432/todo_test?sslmode=disable")
	assert.NoError(suite.T(), err)

	suite.pool, err = pgxpool.NewWithConfig(suite.ctx, config)
	assert.NoError(suite.T(), err)

	suite.repo = NewTodoRepository(suite.pool)

	// Create test table
	_, err = suite.pool.Exec(suite.ctx, `
		CREATE TABLE IF NOT EXISTS todos_test (
			id UUID PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			completed BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
    `)
	assert.NoError(suite.T(), err)
}

func (suite *TodoRepositoryTestSuite) SetupTest() {
	// Clear the test table before each test
	_, err := suite.pool.Exec(suite.ctx, "TRUNCATE TABLE todos_test RESTART IDENTITY")
	assert.NoError(suite.T(), err)

	// Create a test TODO
	suite.testTodo = &domain.Todo{
		ID:          uuid.New(),
		Title:       "Test TODO",
		Description: "Test Description",
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (suite *TodoRepositoryTestSuite) TearDownTest() {
	if suite.pool != nil {
		suite.pool.Close()
	}
}

func (suite *TodoRepositoryTestSuite) TestFindAll() {
	// Create multiple TODOs
	todos := []*domain.Todo{
		{
			ID:          uuid.New(),
			Title:       "Todo 1",
			Description: "Description 1",
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Title:       "Todo 2",
			Description: "Description 2",
			Completed:   true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, v := range todos {
		err := suite.repo.Create(suite.ctx, v)
		assert.NoError(suite.T(), err)
	}

	foundTodos, err := suite.repo.FindAll(suite.ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), len(todos), len(foundTodos))
}

func (suite *TodoRepositoryTestSuite) TestFindById() {
	// Create a TODO
	err := suite.repo.Create(suite.ctx, suite.testTodo)
	assert.NoError(suite.T(), err)

	// Find the TODO created
	todo, err := suite.repo.FindByID(suite.ctx, suite.testTodo.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), todo.ID, suite.testTodo.ID)
}

func (suite *TodoRepositoryTestSuite) TestFindById_NotFound() {
	// Try to find an nonexistent TODO
	nonExistentID := uuid.New()
	_, err := suite.repo.FindByID(suite.ctx, nonExistentID)
	assert.Error(suite.T(), err)
}

func (suite *TodoRepositoryTestSuite) TestCreateTodo() {
	err := suite.repo.Create(suite.ctx, suite.testTodo)
	assert.NoError(suite.T(), err)

	// Validate TODO was created
	var count int
	suite.pool.QueryRow(suite.ctx, `
		SELECT 
		    COUNT(*) 
		FROM todos_test 
		WHERE id = $1
	`, suite.testTodo.ID).Scan(&count)
	assert.Equal(suite.T(), 1, count)
}

func (suite *TodoRepositoryTestSuite) TestCreateTodo_DuplicatedID() {
	err := suite.repo.Create(suite.ctx, suite.testTodo)
	assert.Error(suite.T(), err)

	// Try to create a TODO with the same ID
	err = suite.repo.Create(suite.ctx, suite.testTodo)
	assert.Error(suite.T(), err)
}

func (suite *TodoRepositoryTestSuite) TestUpdateTodo() {
	// Create a TODO
	err := suite.repo.Create(suite.ctx, suite.testTodo)
	assert.NoError(suite.T(), err)

	// Update TODO information
	suite.testTodo.Description = "Updated description"
	suite.testTodo.UpdatedAt = time.Now()

	err = suite.repo.Update(suite.ctx, suite.testTodo)
	assert.NoError(suite.T(), err)

	// Validate TODO was updated
	todo, err := suite.repo.FindByID(suite.ctx, suite.testTodo.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), todo.Description, suite.testTodo.Description)
}

func (suite *TodoRepositoryTestSuite) TestUpdateTodo_NotFound() {
	nonExistentTodo := &domain.Todo{
		ID:          uuid.New(),
		Title:       "Non-existent",
		Description: "Should not exists",
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := suite.repo.Update(suite.ctx, nonExistentTodo)
	assert.Error(suite.T(), err)
}

func (suite *TodoRepositoryTestSuite) TestDeleteTodo() {
	// Create a TODO to be deleted
	err := suite.repo.Create(suite.ctx, suite.testTodo)

	// Delete TODO
	err = suite.repo.Delete(suite.ctx, suite.testTodo.ID)
	assert.NoError(suite.T(), err)

	// Validate TODO was deleted
	_, err = suite.repo.FindByID(suite.ctx, suite.testTodo.ID)
	assert.Error(suite.T(), err)
}

func (suite *TodoRepositoryTestSuite) TestDeleteTodo_NotFound() {
	nonExistentID := uuid.New()
	err := suite.repo.Delete(suite.ctx, nonExistentID)
	assert.Error(suite.T(), err)
}

func TestTodoRepositoryTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, new(TodoRepositoryTestSuite))
}

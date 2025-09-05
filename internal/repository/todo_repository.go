package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leandrowiemesfilho/go-todo-cli/internal/domain"
)

type TodoRepository struct {
	db *pgxpool.Pool
}

func NewTodoRepository(db *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) FindAll(ctx context.Context) ([]*domain.Todo, error) {
	query := `
			SELECT 
    			id, 
    			title, 
    			description, 
    			completed, 
    			created_at, 
				updated_at 
			FROM todos
			ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*domain.Todo
	for rows.Next() {
		var todo domain.Todo

		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		todos = append(todos, &todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *TodoRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error) {
	query := `
			SELECT
				id, 
    			title, 
    			description, 
    			completed, 
    			created_at, 
				updated_at 
			FROM todos
			WHERE id = $1
	`
	var todo domain.Todo

	err := r.db.QueryRow(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *TodoRepository) Create(ctx context.Context, todo *domain.Todo) error {
	query := `
			INSERT INTO todos (id, title, description, completed, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query,
		todo.ID, todo.Title, todo.Description, todo.Completed, todo.CreatedAt, todo.UpdatedAt)
	return err
}

func (r *TodoRepository) Update(ctx context.Context, todo *domain.Todo) error {
	query := `
			UPDATE todos
			SET title = $1, description = $2, updated_at = $3
			WHERE id = $4
	`
	_, err := r.db.Exec(ctx, query,
		todo.Title, todo.Description, todo.UpdatedAt, todo.ID)
	return err
}

func (r *TodoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
			DELETE FROM todos
			WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

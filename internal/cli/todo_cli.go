package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leandrowiemesfilho/go-todo-cli/config"
	"github.com/leandrowiemesfilho/go-todo-cli/internal/domain"
	"github.com/leandrowiemesfilho/go-todo-cli/internal/repository"
	"github.com/leandrowiemesfilho/go-todo-cli/internal/service"
	"github.com/spf13/cobra"
)

type CLI struct {
	rootCmd     *cobra.Command
	todoService service.TodoService
	dbPool      *pgxpool.Pool
}

func NewCLI() *CLI {
	// Load configuration
	cfg := config.LoadConfig()

	// Create database connection pool
	dbPool, err := createDBPool(cfg)
	if err != nil {
		log.Fatalf("Unable to create database connection pool: %v\n", err)
	}

	// Initialize repository and service
	repo := repository.NewTodoRepository(dbPool)
	todoService := service.NewTodoService(repo)

	cli := &CLI{
		todoService: todoService,
		dbPool:      dbPool,
	}

	cli.setupRootCommand()
	return cli
}

func createDBPool(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := cfg.GetPostgresDSN()
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %v", err)
	}

	// Set connection pool settings
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	ctx := context.Background()
	dbPool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	// Test the connection
	if err := dbPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	fmt.Println("✅ Successfully connected to PostgreSQL database")
	return dbPool, nil
}

func (cli *CLI) setupRootCommand() {
	cli.rootCmd = &cobra.Command{
		Use:   "todo",
		Short: "A simple CLI todo application",
		Long:  "A command-line interface for managing your todos with persistence",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Ensure database connection is healthy
			if err := cli.dbPool.Ping(context.Background()); err != nil {
				fmt.Printf("❌ Database connection lost: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cli.rootCmd.AddCommand(
		cli.findAllCommand(),
		cli.findByIDCommand(),
		cli.createCommand(),
		cli.updateCommand(),
		cli.deleteCommand(),
		cli.toggleCommand(),
	)
}

func (cli *CLI) Execute() error {
	defer cli.dbPool.Close()
	return cli.rootCmd.Execute()
}

func (cli *CLI) findAllCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all todos",
		Run: func(cmd *cobra.Command, args []string) {
			filterCompleted, _ := cmd.Flags().GetBool("completed")
			filterPending, _ := cmd.Flags().GetBool("pending")

			todos, err := cli.todoService.FindAllTodos(context.Background())
			if err != nil {
				fmt.Printf("Error getting TODOs: %v\n", err)
				return
			}

			if len(todos) == 0 {
				fmt.Println("No TODOs found")
			}

			var filteredTodos []*domain.Todo
			for _, todo := range todos {
				if filterCompleted && !todo.Completed {
					continue
				} else if filterPending && todo.Completed {
					continue
				}

				filteredTodos = append(filteredTodos, todo)

				cli.printTodoTable(filteredTodos)
			}
		},
	}
}

func (cli *CLI) findByIDCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "find [id]",
		Short: "Find a specific todo by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := uuid.Parse(args[0])
			if err != nil {
				fmt.Printf("Error parsing id: %v\n", err)
				return
			}

			todo, err := cli.todoService.FindTodoByID(context.Background(), id)
			if err != nil {
				fmt.Printf("Error getting TODO: %v\n", err)
				return
			}

			cli.printTodo(todo)
		},
	}
}

func (cli *CLI) createCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [title]",
		Short: "Create a new TODO item",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			desc, _ := cmd.Flags().GetString("description")
			request := domain.CreateTodoRequest{
				Title:       args[0],
				Description: desc,
			}

			todo, err := cli.todoService.CreateTodo(context.Background(), request)
			if err != nil {
				fmt.Printf("Error creating TODO %v\n", err)
				return
			}

			fmt.Printf("TODO created successfully!\n")
			cli.printTodo(todo)
		},
	}

	cmd.Flags().StringP("title", "t", "", "New title for the todo")
	cmd.Flags().StringP("description", "d", "", "New description for the todo")

	return cmd
}

func (cli *CLI) updateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [id]",
		Short: "Update a TODO item",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := uuid.Parse(args[0])
			if err != nil {
				fmt.Printf("Error parsing id: %v\n", err)
				return
			}

			title, _ := cmd.Flags().GetString("title")
			description, _ := cmd.Flags().GetString("description")
			request := domain.UpdateTodoRequest{
				ID:          id,
				Title:       title,
				Description: description,
			}

			todo, err := cli.todoService.UpdateTodo(context.Background(), request)
			if err != nil {
				fmt.Printf("Error trying to update TOD item %v\n", err)
				return
			}

			fmt.Printf("TODO updated successfully!\n")
			cli.printTodo(todo)
		},
	}

	cmd.Flags().StringP("title", "t", "", "New title for the todo")
	cmd.Flags().StringP("description", "d", "", "New description for the todo")

	return cmd
}

func (cli *CLI) deleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete a TODO item by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := uuid.Parse(args[0])
			if err != nil {
				fmt.Printf("Error parsing id: %v\n", err)
				return
			}

			if err = cli.todoService.DeleteTodo(context.Background(), id); err != nil {
				fmt.Printf("Error trying to delete a TODO item %v\n", err)
				return
			}

			fmt.Printf("Todo deleted successfully!\n")
		},
	}
}

func (cli *CLI) toggleCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "toggle [id]",
		Short: "Toggle todo completion status",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := uuid.Parse(args[0])
			if err != nil {
				fmt.Printf("Error parsing id: %v\n", err)
				return
			}

			todo, err := cli.todoService.ToggleTodo(context.Background(), id)
			if err != nil {
				fmt.Printf("Error toggling TODO: %v\n", err)
				return
			}

			status := "completed"
			if !todo.Completed {
				status = "pending"
			}
			fmt.Printf("Todo marked as %s!\n", status)
			cli.printTodo(todo)
		},
	}
}

func (cli *CLI) printTodoTable(todos []*domain.Todo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTITLE\tSTATUS\tCREATION DATE")

	for _, todo := range todos {
		status := "❌ Pending"
		if todo.Completed {
			status = "✅ Completed"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			todo.ID.String()[:8],
			truncate(todo.Title, 20),
			status,
			todo.CreatedDate.Format("2006-01-02 15:04"),
		)
	}
	w.Flush()
}

func (cli *CLI) printTodo(todo *domain.Todo) {
	status := "❌ Pending"
	if todo.Completed {
		status = "✅ Completed"
	}

	fmt.Printf("\nTodo Details:\n")
	fmt.Printf("  ID:          %s\n", todo.ID)
	fmt.Printf("  Title:       %s\n", todo.Title)
	if todo.Description != "" {
		fmt.Printf("  Description: %s\n", todo.Description)
	}
	fmt.Printf("  Status:      %s\n", status)
	fmt.Printf("  Created:     %s\n", todo.CreatedDate.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Updated:     %s\n", todo.UpdatedDate.Format("2006-01-02 15:04:05"))
	fmt.Println()
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

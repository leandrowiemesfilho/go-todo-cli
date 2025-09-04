# TODO CLI Application with PostgreSQL

A command-line todo application built with Go and PostgreSQL, following clean architecture principles.

## Features

- ✅ Create, read, update, and delete todos
- ✅ Toggle todo completion status
- ✅ PostgreSQL persistence with connection pooling
- ✅ Docker and Docker Compose support
- ✅ Database migrations
- ✅ Clean, tabular output
- ✅ Filter todos by status
- ✅ UUID-based identification

## Quick Start with Docker

```bash
# Start the application with Docker Compose
docker compose up -d

# List TODOs
docker compose exec go-todo-cli ./go-todo-cli list

# Add a new TODO
docker compose exec go-todo-cli ./go-todo-cli create "Learn Go" --description "Study Go programming language"
```

## Local Development

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Docker (optional)

### Installation
```bash
git clone https://github.com/leandrowiemesfilho/go-todo-cli.git
cd go-todo-cli
go mod tidy
go build -o go-todo-cli cmd/main.go
```

### Database setup
1. Start PostgreSQL:
```bash
docker compose up -d postgres 
```

### Environment variables
Create a `.env` file:
```env
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=todo_user
POSTGRES_PASSWORD=todo_password
POSTGRES_DB=todo_db
POSTGRES_SSL_MODE=disable
```

## Usage
```bash
# Create a new TODO
./go-todo-cli create "Learn Go" --description "Study Go programming language"

# List all TODOs
./go-todo-cli list

# List only completed TODOs
./go-todo-cli list --completed

# List only pending TODOs
./go-todo-cli list --pending

# Get a specific TODO
./go-todo-cli find <todo-id>

# Update a TODO
./go-todo-cli update <todo-id> --title "New title" --description "New description" --completed

# Toggle TODO completion
./go-todo-cli toggle <todo-id>

# Delete a TODO
./go-todo-cli delete <todo-id>
```
## Docker commands
```bash
# Build the image
docker compose build

# Start all services
docker compose up -d

# Run specific commands
docker compose exec go-todo-cli ./go-todo-cli list
docker compose exec go-todo-cli ./go-todo-cli add "Docker todo"

# View logs
docker compose logs go-todo-cli
docker compose logs postgres

# Stop services
docker compose down
```

## Project structure
```text
go-todo-cli/
├── cmd/                # Application entry point
├── internal/           # Private application code
│   ├── domain/         # Business entities and interfaces
│   ├── repository/     # PostgreSQL data access layer
│   ├── service/        # Business logic
│   └── cli/            # CLI command handlers
├── migrations/         # Database migration files
├── config/             # Configuration management
├── pkg/                # Public utility packages
├── Dockerfile          # Docker build instructions
├── docker compose.yml  # Docker Compose configuration
└── go.mod              # Go module definition
```

## Database schema
The application uses a simple todos table:

| Column          | Type      | Description           |
|:----------------|:----------|:----------------------|
| **ID**          | UUID      | Primary key           |
| **Title**       | VARCHAR   | TODO title            |
| **Description** | TEXT      | Optional description  |
| **Completed**   | BOOLEAN   | Completion status     |
| **Created at**  | TIMESTAMP | Creation timestamp    |
| **Updated at**  | TIMESTAMP | Last update timestamp |

## Clean architecture
This project follows clean architecture principles:
- **Domain:** Core business entities and interfaces
- **Repository:** PostgreSQL data access layer
- **Service:** Business logic and use cases
- **CLI:** Presentation layer (user interface)
- **Config:** Environment configuration management
Each layer depends only on inner layers, making the code highly testable and maintainable.
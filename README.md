# To-Do List API

[![Go Tests](https://github.com/TheOtherDavid/to-do-list/actions/workflows/go.yml/badge.svg)](https://github.com/TheOtherDavid/to-do-list/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/TheOtherDavid/to-do-list/graph/badge.svg?token=T8MWLVA0C2)](https://codecov.io/gh/TheOtherDavid/to-do-list)

There are many like it, but this one is mine.

A task management API built with Go. Features task creation, completion tracking, and support for recurring tasks through templates. Uses CSV file storage for simplicity.

## Features

- Task creation and management
- Recurring task templates
- Automatic task generation from templates
- Task completion tracking
- CSV-based persistence
- Comprehensive test coverage

## Project Structure

```
.
├── cmd/
│   ├── api/          # Main server application
│   └── refresh_tasks/ # Task generation service
├── handlers/         # HTTP request handlers
├── jobs/            # Background task processors
├── models/          # Data models
├── routes/          # API route definitions
└── storage/         # Data persistence layer
    └── mocks/       # Test mocks
```

## API Endpoints

### Tasks
- `POST /tasks` - Create a new task
- `GET /tasks/uncompleted` - Retrieve uncompleted tasks
- `GET /tasks/completed` - Retrieve completed tasks
- `POST /tasks/{id}/complete` - Mark a task as completed

### Task Templates
- `POST /templates` - Create a task template
- `GET /templates` - Retrieve all task templates
- `GET /templates/{id}` - Retrieve a specific template
- `PUT /templates/{id}` - Update a template
- `DELETE /templates/{id}` - Delete a template

## Installation

1. Prerequisites:
   - Go 1.21 or higher

2. Clone the repository:
   ```bash
   git clone https://github.com/TheOtherDavid/to-do-list.git
   cd to-do-list
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Build and run:
   ```bash
   go run cmd/api/main.go
   ```

The server will start on `http://localhost:8080`.

## Task Refresh From Templates

There's a utility that creates new tasks from your templates:

```bash
go run cmd/refresh_tasks/main.go
```

You can set it up with cron if you want it to run automatically.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. MIT Licensed - do whatever you want with it!

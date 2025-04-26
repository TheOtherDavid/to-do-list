# To-Do List API

[![Go Tests](https://github.com/TheOtherDavid/to-do-list/actions/workflows/go.yml/badge.svg)](https://github.com/TheOtherDavid/to-do-list/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/TheOtherDavid/to-do-list)](https://goreportcard.com/report/github.com/TheOtherDavid/to-do-list)
[![codecov](https://codecov.io/gh/TheOtherDavid/to-do-list/branch/main/graph/badge.svg)](https://codecov.io/gh/TheOtherDavid/to-do-list)

> There are many like it, but this one is mine.

A task management API built with Go. Tracks tasks, and also has support for recurring tasks (called TaskTemplates). Uses CSV files for storage because we're doing it quick & dirty.

## What's Cool About It

- ✨ Create and track your tasks
- 📋 Task templates for recurring tasks
- 🔄 Auto-refreshes tasks from templates
- 📊 Keeps track of what's done and what's not
- 💾 Saves everything to CSV files
- 🧪 Well-tested

## What's Inside

```
.
├── cmd/
│   ├── api/          # The main server
│   └── refresh_tasks/ # Auto-task creator
├── handlers/         # Handles web stuff
├── jobs/            # Background tasks
├── models/          # Data structures
├── routes/          # API endpoints
└── storage/         # Saves your stuff
    └── mocks/       # Test helpers
```

## API Stuff You Can Do

- `POST /tasks` - Make a new task
- `GET /tasks/uncompleted` - See what's not done
- `GET /tasks/completed` - Check out finished tasks
- `POST /tasks/{id}/complete` - Mark something as done ✅

## Getting Started

1. Make sure you have Go 1.21+ installed

2. Grab the code:
   ```bash
   git clone https://github.com/TheOtherDavid/to-do-list.git
   cd to-do-list
   ```

3. Get the dependencies:
   ```bash
   go mod download
   ```

4. Fire it up:
   ```bash
   go run cmd/api/main.go
   ```

It'll be running at `http://localhost:8080`. Have fun! 🚀

## Task Refresh From Templates

There's a utility that creates new tasks from your templates:

```bash
go run cmd/refresh_tasks/main.go
```

You can set it up with cron if you want it to run automatically. ⏰

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. MIT Licensed - do whatever you want with it! 🎉

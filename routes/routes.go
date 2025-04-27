package routes

import (
	"github.com/TheOtherDavid/to-do-list/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/tasks", handlers.CreateTask).Methods("POST")
	r.HandleFunc("/tasks", handlers.GetUncompletedTasks).Methods("GET")
	r.HandleFunc("/tasks/completed", handlers.GetCompletedTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}/complete", handlers.CompleteTask).Methods("PUT")

	r.HandleFunc("/task-templates", handlers.CreateTaskTemplate).Methods("POST")
	r.HandleFunc("/task-templates", handlers.GetAllTaskTemplates).Methods("GET")
	r.HandleFunc("/task-templates/{id}", handlers.GetTaskTemplate).Methods("GET")
	r.HandleFunc("/task-templates/{id}", handlers.UpdateTaskTemplate).Methods("PUT")
	r.HandleFunc("/task-templates/{id}", handlers.DeleteTaskTemplate).Methods("DELETE")

	return r
}

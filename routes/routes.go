package routes

import (
	"github.com/TheOtherDavid/to-do-list/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/tasks", handlers.CreateTask).Methods("POST")
	r.HandleFunc("/tasks", handlers.ListTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", handlers.DeleteTask).Methods("DELETE")
	r.HandleFunc("/tasks/{id}/complete", handlers.CompleteTask).Methods("PUT")

	return r
}

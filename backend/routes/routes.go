package routes

import (
	"net/http"

	"github.com/TheOtherDavid/to-do-list/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.Use(corsMiddleware)

	r.HandleFunc("/tasks", handlers.CreateTask).Methods("POST", "OPTIONS")
	r.HandleFunc("/tasks", handlers.GetUncompletedTasks).Methods("GET")
	r.HandleFunc("/tasks/completed", handlers.GetCompletedTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}/complete", handlers.CompleteTask).Methods("PUT", "OPTIONS")

	r.HandleFunc("/task-templates", handlers.CreateTaskTemplate).Methods("POST", "OPTIONS")
	r.HandleFunc("/task-templates", handlers.GetAllTaskTemplates).Methods("GET")
	r.HandleFunc("/task-templates/{id}", handlers.GetTaskTemplate).Methods("GET")
	r.HandleFunc("/task-templates/{id}", handlers.UpdateTaskTemplate).Methods("PUT", "OPTIONS")
	r.HandleFunc("/task-templates/{id}", handlers.DeleteTaskTemplate).Methods("DELETE", "OPTIONS")

	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

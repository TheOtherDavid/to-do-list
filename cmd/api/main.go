package main

import (
	"log"
	"net/http"
	"time"

	"github.com/TheOtherDavid/to-do-list/handlers"
	"github.com/TheOtherDavid/to-do-list/routes"
)

func main() {
	handlers.InitStorage("task_instance.csv", "task_template.csv")
	router := routes.SetupRoutes()

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on http://localhost%s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

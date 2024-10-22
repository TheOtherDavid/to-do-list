package main

import (
	"log"
	"net/http"

	"github.com/TheOtherDavid/to-do-list/handlers"
	"github.com/TheOtherDavid/to-do-list/routes"
)

func main() {
	handlers.InitStorage("task_instance.csv", "task_template.csv")
	router := routes.SetupRoutes()

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

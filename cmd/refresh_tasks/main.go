package main

import (
	"log"

	"github.com/TheOtherDavid/to-do-list/jobs"
)

func main() {
	jobs.InitStorage("task_instance.csv", "task_template.csv")
	err := jobs.RefreshTasks()
	if err != nil {
		log.Fatal(err)
	}
}

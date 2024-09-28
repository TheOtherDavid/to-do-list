package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/TheOtherDavid/to-do-list/models"
	"github.com/TheOtherDavid/to-do-list/storage"
)

var csvStorage *storage.CSVStorage

func InitStorage(filename string) {
	csvStorage = storage.NewCSVStorage(filename)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task.ID = uuid.New().String()
	task.CreatedAt = time.Now()

	err = csvStorage.CreateTask(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := csvStorage.ListTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Printf("Deleting task with ID: %s", id)

	err := csvStorage.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func CompleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Printf("Completing task with ID: %s", id)

	tasks, err := csvStorage.ListTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	taskFound := false
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Completed = true
			now := time.Now()
			tasks[i].CompletedAt = &now
			taskFound = true
			break
		}
	}

	if !taskFound {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	err = csvStorage.UpdateTasks(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Task %s completed successfully", id)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Task completed successfully"})
}

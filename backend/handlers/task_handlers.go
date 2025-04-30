package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/TheOtherDavid/to-do-list/models"
	"github.com/TheOtherDavid/to-do-list/storage"
)

var csvStorage storage.Storage

func InitStorage(taskInstanceFile, taskTemplateFile string) {
	csvStorage = storage.NewCSVStorage(taskInstanceFile, taskTemplateFile)
}

func SetStorage(s storage.Storage) {
	csvStorage = s
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.TaskInstance
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task.ID = uuid.New().String()
	task.CreatedAt = time.Now()

	err = csvStorage.SaveTaskInstance(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := csvStorage.GetAllTaskInstances()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetUncompletedTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := csvStorage.GetUncompletedTaskInstances()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetCompletedTasks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	limit, err := strconv.Atoi(vars["limit"])
	if err != nil {
		http.Error(w, "Invalid limit value", http.StatusBadRequest)
		return
	}
	offset, err := strconv.Atoi(vars["offset"])
	if err != nil {
		http.Error(w, "Invalid offset value", http.StatusBadRequest)
		return
	}

	tasks, err := csvStorage.GetCompletedTaskInstances(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func CompleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Printf("Completing task with ID: %s", id)

	tasks, err := csvStorage.GetUncompletedTaskInstances()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	taskFound := false
	for _, task := range tasks {
		if task.ID == id {
			taskFound = true
			break
		}
	}

	if !taskFound {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	err = csvStorage.SetCompleted(id, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Task %s completed successfully", id)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Task completed successfully"}); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

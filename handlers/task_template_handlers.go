package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/TheOtherDavid/to-do-list/models"
)

func CreateTaskTemplate(w http.ResponseWriter, r *http.Request) {
	var template models.TaskTemplate
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	template.ID = uuid.New().String()
	template.CreatedAt = time.Now()

	if err := csvStorage.SaveTaskTemplate(template); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(template)
}

func GetAllTaskTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := csvStorage.GetAllTaskTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(templates)
}

func GetTaskTemplate(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	template, err := csvStorage.GetTaskTemplate(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(template)
}

func UpdateTaskTemplate(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var template models.TaskTemplate
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	template.ID = id
	if err := csvStorage.UpdateTaskTemplate(template); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(template)
}

func DeleteTaskTemplate(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := csvStorage.DeleteTaskTemplate(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

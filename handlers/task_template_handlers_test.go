package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"github.com/TheOtherDavid/to-do-list/models"
	"github.com/TheOtherDavid/to-do-list/storage"
	storageMocks "github.com/TheOtherDavid/to-do-list/storage/mocks"
)

type templateTestHelper struct {
	mockStorage *storageMocks.MockStorage
	recorder    *httptest.ResponseRecorder
	origStorage storage.Storage
}

func setupTemplateTest(t *testing.T) *templateTestHelper {
	th := &templateTestHelper{
		mockStorage: &storageMocks.MockStorage{
			Templates: []models.TaskTemplate{},
		},
		recorder: httptest.NewRecorder(),
	}

	// Save and replace storage
	th.origStorage = csvStorage
	csvStorage = th.mockStorage

	t.Cleanup(func() {
		csvStorage = th.origStorage
	})

	return th
}

func (th *templateTestHelper) makeRequest(method, path string, body interface{}) *http.Request {
	var bodyReader io.Reader

	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		bodyReader = bytes.NewBuffer(jsonBytes)
	}

	req := httptest.NewRequest(method, path, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

func TestCreateTaskTemplate(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		th := setupTemplateTest(t)

		template := models.TaskTemplate{
			Title:          "Test Template",
			Description:    "Test Description",
			RecurrenceRule: "DAILY",
		}

		CreateTaskTemplate(th.recorder, th.makeRequest(http.MethodPost, "/task-templates", template))

		if status := th.recorder.Code; status != http.StatusCreated {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		var responseTemplate models.TaskTemplate
		if err := json.NewDecoder(th.recorder.Body).Decode(&responseTemplate); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Verify the response
		if responseTemplate.Title != template.Title {
			t.Errorf("Expected template title %v, got %v", template.Title, responseTemplate.Title)
		}
		if responseTemplate.Description != template.Description {
			t.Errorf("Expected template description %v, got %v", template.Description, responseTemplate.Description)
		}
		if responseTemplate.RecurrenceRule != template.RecurrenceRule {
			t.Errorf("Expected template recurrence rule %v, got %v", template.RecurrenceRule, responseTemplate.RecurrenceRule)
		}
		if responseTemplate.ID == "" {
			t.Error("Expected template ID to be set")
		}
		if responseTemplate.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}

		// Verify storage
		if th.mockStorage.SavedTemplate == nil {
			t.Fatal("Expected template to be saved to storage")
		}
		if th.mockStorage.SavedTemplate.Title != template.Title {
			t.Errorf("Saved template title doesn't match: expected %v, got %v", template.Title, th.mockStorage.SavedTemplate.Title)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		th := setupTemplateTest(t)

		invalidJSON := []byte(`{"title": "Test Template", "description": "Test Description"`) // missing closing brace
		req := httptest.NewRequest(http.MethodPost, "/task-templates", bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")

		CreateTaskTemplate(th.recorder, req)

		if status := th.recorder.Code; status != http.StatusBadRequest {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		if th.mockStorage.SavedTemplate != nil {
			t.Error("Expected no template to be saved")
		}
	})

	t.Run("storage error", func(t *testing.T) {
		th := setupTemplateTest(t)
		th.mockStorage.TemplateErr = errors.New("storage error")

		template := models.TaskTemplate{
			Title:          "Test Template",
			Description:    "Test Description",
			RecurrenceRule: "DAILY",
		}

		CreateTaskTemplate(th.recorder, th.makeRequest(http.MethodPost, "/task-templates", template))

		if status := th.recorder.Code; status != http.StatusInternalServerError {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}

func TestGetAllTaskTemplates(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		th := setupTemplateTest(t)

		// Setup mock data
		expectedTemplates := []models.TaskTemplate{
			{
				ID:             "1",
				Title:          "Template 1",
				Description:    "Description 1",
				RecurrenceRule: "DAILY",
				CreatedAt:      time.Now(),
			},
			{
				ID:             "2",
				Title:          "Template 2",
				Description:    "Description 2",
				RecurrenceRule: "WEEKLY",
				CreatedAt:      time.Now(),
			},
		}
		th.mockStorage.Templates = expectedTemplates

		GetAllTaskTemplates(th.recorder, th.makeRequest(http.MethodGet, "/task-templates", nil))

		if status := th.recorder.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var responseTemplates []models.TaskTemplate
		if err := json.NewDecoder(th.recorder.Body).Decode(&responseTemplates); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(responseTemplates) != len(expectedTemplates) {
			t.Errorf("Expected %d templates, got %d", len(expectedTemplates), len(responseTemplates))
		}

		for i, template := range responseTemplates {
			if template.ID != expectedTemplates[i].ID {
				t.Errorf("Expected template ID %s, got %s", expectedTemplates[i].ID, template.ID)
			}
			if template.Title != expectedTemplates[i].Title {
				t.Errorf("Expected template title %s, got %s", expectedTemplates[i].Title, template.Title)
			}
		}
	})

	t.Run("storage error", func(t *testing.T) {
		th := setupTemplateTest(t)
		th.mockStorage.TemplateErr = errors.New("storage error")

		GetAllTaskTemplates(th.recorder, th.makeRequest(http.MethodGet, "/task-templates", nil))

		if status := th.recorder.Code; status != http.StatusInternalServerError {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}

func TestGetTaskTemplate(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		th := setupTemplateTest(t)

		// Setup mock data
		templateID := "test-template-id"
		expectedTemplate := models.TaskTemplate{
			ID:             templateID,
			Title:          "Template 1",
			Description:    "Description 1",
			RecurrenceRule: "DAILY",
			CreatedAt:      time.Now(),
		}
		th.mockStorage.Templates = []models.TaskTemplate{expectedTemplate}

		// Create request with template ID
		req := th.makeRequest(http.MethodGet, "/task-templates/"+templateID, nil)
		req = mux.SetURLVars(req, map[string]string{"id": templateID})

		GetTaskTemplate(th.recorder, req)

		if status := th.recorder.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var responseTemplate models.TaskTemplate
		if err := json.NewDecoder(th.recorder.Body).Decode(&responseTemplate); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if responseTemplate.ID != expectedTemplate.ID {
			t.Errorf("Expected template ID %s, got %s", expectedTemplate.ID, responseTemplate.ID)
		}
		if responseTemplate.Title != expectedTemplate.Title {
			t.Errorf("Expected template title %s, got %s", expectedTemplate.Title, responseTemplate.Title)
		}
		if responseTemplate.Description != expectedTemplate.Description {
			t.Errorf("Expected template description %s, got %s", expectedTemplate.Description, responseTemplate.Description)
		}
		if responseTemplate.RecurrenceRule != expectedTemplate.RecurrenceRule {
			t.Errorf("Expected template recurrence rule %s, got %s", expectedTemplate.RecurrenceRule, responseTemplate.RecurrenceRule)
		}
	})

	t.Run("template not found", func(t *testing.T) {
		th := setupTemplateTest(t)
		th.mockStorage.TemplateErr = errors.New("template not found")

		// Create request with non-existent template ID
		templateID := "non-existent-id"
		req := th.makeRequest(http.MethodGet, "/task-templates/"+templateID, nil)
		req = mux.SetURLVars(req, map[string]string{"id": templateID})

		GetTaskTemplate(th.recorder, req)

		if status := th.recorder.Code; status != http.StatusNotFound {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})
}

func TestUpdateTaskTemplate(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		th := setupTemplateTest(t)

		// Setup mock data
		templateID := "test-template-id"
		existingTemplate := models.TaskTemplate{
			ID:             templateID,
			Title:          "Original Title",
			Description:    "Original Description",
			RecurrenceRule: "DAILY",
			CreatedAt:      time.Now().Add(-24 * time.Hour), // Created yesterday
		}
		th.mockStorage.Templates = []models.TaskTemplate{existingTemplate}

		// Create updated template
		updatedTemplate := models.TaskTemplate{
			ID:             templateID,
			Title:          "Updated Title",
			Description:    "Updated Description",
			RecurrenceRule: "WEEKLY",
			CreatedAt:      existingTemplate.CreatedAt, // Should remain unchanged
		}

		// Create request with template ID and updated data
		req := th.makeRequest(http.MethodPut, "/task-templates/"+templateID, updatedTemplate)
		req = mux.SetURLVars(req, map[string]string{"id": templateID})

		UpdateTaskTemplate(th.recorder, req)

		if status := th.recorder.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Verify the template was updated in storage
		if th.mockStorage.LastUpdatedTemplate == nil {
			t.Fatal("Expected template to be updated in storage")
		}
		if th.mockStorage.LastUpdatedTemplate.Title != updatedTemplate.Title {
			t.Errorf("Updated template title doesn't match: expected %v, got %v", 
				updatedTemplate.Title, th.mockStorage.LastUpdatedTemplate.Title)
		}
		if th.mockStorage.LastUpdatedTemplate.Description != updatedTemplate.Description {
			t.Errorf("Updated template description doesn't match: expected %v, got %v", 
				updatedTemplate.Description, th.mockStorage.LastUpdatedTemplate.Description)
		}
		if th.mockStorage.LastUpdatedTemplate.RecurrenceRule != updatedTemplate.RecurrenceRule {
			t.Errorf("Updated template recurrence rule doesn't match: expected %v, got %v", 
				updatedTemplate.RecurrenceRule, th.mockStorage.LastUpdatedTemplate.RecurrenceRule)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		th := setupTemplateTest(t)

		templateID := "test-template-id"
		invalidJSON := []byte(`{"title": "Updated Title", "description": "Updated Description"`) // missing closing brace
		req := httptest.NewRequest(http.MethodPut, "/task-templates/"+templateID, bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		req = mux.SetURLVars(req, map[string]string{"id": templateID})

		UpdateTaskTemplate(th.recorder, req)

		if status := th.recorder.Code; status != http.StatusBadRequest {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		if th.mockStorage.LastUpdatedTemplate != nil {
			t.Error("Expected no template to be updated")
		}
	})

	t.Run("template not found", func(t *testing.T) {
		th := setupTemplateTest(t)
		th.mockStorage.TemplateErr = errors.New("template not found")

		templateID := "non-existent-id"
		updatedTemplate := models.TaskTemplate{
			ID:             templateID,
			Title:          "Updated Title",
			Description:    "Updated Description",
			RecurrenceRule: "WEEKLY",
		}

		req := th.makeRequest(http.MethodPut, "/task-templates/"+templateID, updatedTemplate)
		req = mux.SetURLVars(req, map[string]string{"id": templateID})

		UpdateTaskTemplate(th.recorder, req)

		if status := th.recorder.Code; status != http.StatusInternalServerError {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}

func TestDeleteTaskTemplate(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		th := setupTemplateTest(t)

		// Setup mock data
		templateID := "test-template-id"
		template := models.TaskTemplate{
			ID:             templateID,
			Title:          "Template to Delete",
			Description:    "This template will be deleted",
			RecurrenceRule: "DAILY",
			CreatedAt:      time.Now(),
		}
		th.mockStorage.Templates = []models.TaskTemplate{template}

		// Create request with template ID
		req := th.makeRequest(http.MethodDelete, "/task-templates/"+templateID, nil)
		req = mux.SetURLVars(req, map[string]string{"id": templateID})

		DeleteTaskTemplate(th.recorder, req)

		if status := th.recorder.Code; status != http.StatusNoContent {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
		}

		// Verify the template was deleted in storage
		if th.mockStorage.LastDeletedID != templateID {
			t.Errorf("Expected template %s to be deleted, got %s", templateID, th.mockStorage.LastDeletedID)
		}
	})

	t.Run("template not found", func(t *testing.T) {
		th := setupTemplateTest(t)
		th.mockStorage.TemplateErr = errors.New("template not found")

		// Create request with non-existent template ID
		templateID := "non-existent-id"
		req := th.makeRequest(http.MethodDelete, "/task-templates/"+templateID, nil)
		req = mux.SetURLVars(req, map[string]string{"id": templateID})

		DeleteTaskTemplate(th.recorder, req)

		if status := th.recorder.Code; status != http.StatusInternalServerError {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}

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

type testHelper struct {
	mockStorage *storageMocks.MockStorage
	recorder    *httptest.ResponseRecorder
	origStorage storage.Storage
}

func setupTest(t *testing.T) *testHelper {
	th := &testHelper{
		mockStorage: &storageMocks.MockStorage{
			UncompletedTasks: []models.TaskInstance{},
			CompletedTasks:   []models.TaskInstance{},
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

func (th *testHelper) makeRequest(method, path string, body interface{}) *http.Request {
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

func TestCreateTask(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		th := setupTest(t)

		task := models.TaskInstance{
			Title:       "Test Task",
			Description: "Test Description",
		}

		CreateTask(th.recorder, th.makeRequest(http.MethodPost, "/tasks", task))

		if status := th.recorder.Code; status != http.StatusCreated {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		var responseTask models.TaskInstance
		if err := json.NewDecoder(th.recorder.Body).Decode(&responseTask); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Verify the response
		if responseTask.Title != task.Title {
			t.Errorf("Expected task title %v, got %v", task.Title, responseTask.Title)
		}
		if responseTask.Description != task.Description {
			t.Errorf("Expected task description %v, got %v", task.Description, responseTask.Description)
		}
		if responseTask.ID == "" {
			t.Error("Expected task ID to be set")
		}
		if responseTask.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}

		// Verify storage
		if th.mockStorage.SavedTask == nil {
			t.Fatal("Expected task to be saved to storage")
		}
		if th.mockStorage.SavedTask.Title != task.Title {
			t.Errorf("Saved task title doesn't match: expected %v, got %v", task.Title, th.mockStorage.SavedTask.Title)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		th := setupTest(t)

		invalidJSON := []byte(`{"title": "Test Task", "description": "Test Description"`) // missing closing brace
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")

		CreateTask(th.recorder, req)

		if status := th.recorder.Code; status != http.StatusBadRequest {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		if th.mockStorage.SavedTask != nil {
			t.Error("Expected no task to be saved")
		}
	})

	t.Run("storage error", func(t *testing.T) {
		th := setupTest(t)
		th.mockStorage.Err = errors.New("storage error")

		task := models.TaskInstance{
			Title:       "Test Task",
			Description: "Test Description",
		}

		CreateTask(th.recorder, th.makeRequest(http.MethodPost, "/tasks", task))

		if status := th.recorder.Code; status != http.StatusInternalServerError {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}

func TestGetAllTasks(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		th := setupTest(t)

		// Setup mock data
		expectedTasks := []models.TaskInstance{
			{
				ID:          "1",
				Title:       "Task 1",
				Description: "Description 1",
				CreatedAt:   time.Now(),
			},
			{
				ID:          "2",
				Title:       "Task 2",
				Description: "Description 2",
				CreatedAt:   time.Now(),
			},
		}
		th.mockStorage.AllTaskInstances = expectedTasks

		// Make request
		GetAllTasks(th.recorder, th.makeRequest(http.MethodGet, "/tasks", nil))

		// Check status code
		if status := th.recorder.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Check response body
		var responseTasks []models.TaskInstance
		if err := json.NewDecoder(th.recorder.Body).Decode(&responseTasks); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(responseTasks) != len(expectedTasks) {
			t.Errorf("Expected %d tasks, got %d", len(expectedTasks), len(responseTasks))
		}

		for i, task := range expectedTasks {
			if responseTasks[i].ID != task.ID {
				t.Errorf("Task %d: expected ID %s, got %s", i, task.ID, responseTasks[i].ID)
			}
			if responseTasks[i].Title != task.Title {
				t.Errorf("Task %d: expected title %s, got %s", i, task.Title, responseTasks[i].Title)
			}
		}
	})

	t.Run("storage error", func(t *testing.T) {
		th := setupTest(t)
		th.mockStorage.Err = errors.New("storage error")

		GetAllTasks(th.recorder, th.makeRequest(http.MethodGet, "/tasks", nil))

		if status := th.recorder.Code; status != http.StatusInternalServerError {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}

func TestGetUncompletedTasks(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		th := setupTest(t)

		// Setup mock data
		expectedTasks := []models.TaskInstance{
			{
				ID:          "1",
				Title:       "Task 1",
				Description: "Description 1",
				CreatedAt:   time.Now(),
			},
			{
				ID:          "2",
				Title:       "Task 2",
				Description: "Description 2",
				CreatedAt:   time.Now(),
			},
		}
		th.mockStorage.UncompletedTasks = expectedTasks

		// Make request
		GetUncompletedTasks(th.recorder, th.makeRequest(http.MethodGet, "/tasks/uncompleted", nil))

		// Check status code
		if status := th.recorder.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Check response body
		var responseTasks []models.TaskInstance
		if err := json.NewDecoder(th.recorder.Body).Decode(&responseTasks); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(responseTasks) != len(expectedTasks) {
			t.Errorf("Expected %d tasks, got %d", len(expectedTasks), len(responseTasks))
		}

		for i, task := range expectedTasks {
			if responseTasks[i].ID != task.ID {
				t.Errorf("Task %d: expected ID %s, got %s", i, task.ID, responseTasks[i].ID)
			}
			if responseTasks[i].Title != task.Title {
				t.Errorf("Task %d: expected title %s, got %s", i, task.Title, responseTasks[i].Title)
			}
		}
	})

	t.Run("storage error", func(t *testing.T) {
		th := setupTest(t)
		th.mockStorage.Err = errors.New("storage error")

		GetUncompletedTasks(th.recorder, th.makeRequest(http.MethodGet, "/tasks/uncompleted", nil))

		if status := th.recorder.Code; status != http.StatusInternalServerError {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}

func TestGetCompletedTasks(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		th := setupTest(t)

		// Setup mock data
		expectedTasks := []models.TaskInstance{
			{
				ID:          "1",
				Title:       "Task 1",
				Description: "Description 1",
				CreatedAt:   time.Now(),
				Completed:   true,
				CompletedAt: &time.Time{},
			},
		}
		th.mockStorage.CompletedTasks = expectedTasks

		// Create request with mux vars
		req := th.makeRequest(http.MethodGet, "/tasks/completed", nil)
		req = mux.SetURLVars(req, map[string]string{
			"limit":  "10",
			"offset": "0",
		})

		GetCompletedTasks(th.recorder, req)

		// Check status code
		if status := th.recorder.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Check response body
		var responseTasks []models.TaskInstance
		if err := json.NewDecoder(th.recorder.Body).Decode(&responseTasks); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(responseTasks) != len(expectedTasks) {
			t.Errorf("Expected %d tasks, got %d", len(expectedTasks), len(responseTasks))
		}

		for i, task := range expectedTasks {
			if responseTasks[i].ID != task.ID {
				t.Errorf("Task %d: expected ID %s, got %s", i, task.ID, responseTasks[i].ID)
			}
			if responseTasks[i].Title != task.Title {
				t.Errorf("Task %d: expected title %s, got %s", i, task.Title, responseTasks[i].Title)
			}
			if !responseTasks[i].Completed {
				t.Errorf("Task %d: expected completed to be true", i)
			}
		}
	})

	t.Run("invalid limit", func(t *testing.T) {
		th := setupTest(t)

		req := th.makeRequest(http.MethodGet, "/tasks/completed", nil)
		req = mux.SetURLVars(req, map[string]string{
			"limit":  "invalid",
			"offset": "0",
		})

		GetCompletedTasks(th.recorder, req)

		if status := th.recorder.Code; status != http.StatusBadRequest {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("invalid offset", func(t *testing.T) {
		th := setupTest(t)

		req := th.makeRequest(http.MethodGet, "/tasks/completed", nil)
		req = mux.SetURLVars(req, map[string]string{
			"limit":  "10",
			"offset": "invalid",
		})

		GetCompletedTasks(th.recorder, req)

		if status := th.recorder.Code; status != http.StatusBadRequest {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("storage error", func(t *testing.T) {
		th := setupTest(t)
		th.mockStorage.Err = errors.New("storage error")

		req := th.makeRequest(http.MethodGet, "/tasks/completed", nil)
		req = mux.SetURLVars(req, map[string]string{
			"limit":  "10",
			"offset": "0",
		})

		GetCompletedTasks(th.recorder, req)

		if status := th.recorder.Code; status != http.StatusInternalServerError {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}

func TestCompleteTask(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		th := setupTest(t)

		// Setup mock data
		taskID := "test-task-id"
		uncompletedTasks := []models.TaskInstance{
			{
				ID:          taskID,
				Title:       "Task 1",
				Description: "Description 1",
				CreatedAt:   time.Now(),
			},
		}
		th.mockStorage.UncompletedTasks = uncompletedTasks

		// Create request with task ID
		req := th.makeRequest(http.MethodPost, "/tasks/"+taskID+"/complete", nil)
		req = mux.SetURLVars(req, map[string]string{"id": taskID})

		CompleteTask(th.recorder, req)

		// Check status code
		if status := th.recorder.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Verify the task was marked as completed in storage
		if th.mockStorage.LastCompletedID != taskID {
			t.Errorf("Expected task %s to be marked as completed", taskID)
		}
	})

	t.Run("task not found", func(t *testing.T) {
		th := setupTest(t)

		// Setup mock data with no tasks
		th.mockStorage.UncompletedTasks = []models.TaskInstance{}

		// Create request with non-existent task ID
		taskID := "non-existent-id"
		req := th.makeRequest(http.MethodPost, "/tasks/"+taskID+"/complete", nil)
		req = mux.SetURLVars(req, map[string]string{"id": taskID})

		CompleteTask(th.recorder, req)

		// Check status code
		if status := th.recorder.Code; status != http.StatusNotFound {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	t.Run("get tasks error", func(t *testing.T) {
		th := setupTest(t)
		th.mockStorage.Err = errors.New("storage error")

		taskID := "test-task-id"
		req := th.makeRequest(http.MethodPost, "/tasks/"+taskID+"/complete", nil)
		req = mux.SetURLVars(req, map[string]string{"id": taskID})

		CompleteTask(th.recorder, req)

		if status := th.recorder.Code; status != http.StatusInternalServerError {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})

	t.Run("set completed error", func(t *testing.T) {
		th := setupTest(t)

		// Setup mock data
		taskID := "test-task-id"
		uncompletedTasks := []models.TaskInstance{
			{
				ID:          taskID,
				Title:       "Task 1",
				Description: "Description 1",
				CreatedAt:   time.Now(),
			},
		}
		th.mockStorage.UncompletedTasks = uncompletedTasks

		// Set error for SetCompleted operation
		th.mockStorage.SetCompletedErr = errors.New("failed to set completed")

		req := th.makeRequest(http.MethodPost, "/tasks/"+taskID+"/complete", nil)
		req = mux.SetURLVars(req, map[string]string{"id": taskID})

		CompleteTask(th.recorder, req)

		if status := th.recorder.Code; status != http.StatusInternalServerError {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}

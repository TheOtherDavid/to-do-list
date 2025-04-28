package routes_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TheOtherDavid/to-do-list/handlers"
	"github.com/TheOtherDavid/to-do-list/models"
	"github.com/TheOtherDavid/to-do-list/routes"
	"github.com/TheOtherDavid/to-do-list/storage/mocks"
)

func setupTest() {
	mockStorage := mocks.NewMockStorage()
	mockStorage.UncompletedTasks = []models.TaskInstance{}
	mockStorage.CompletedTasks = []models.TaskInstance{}

	handlers.SetStorage(mockStorage)
}

func setupTestWithTask(id string) {
	mockStorage := mocks.NewMockStorage()
	mockStorage.UncompletedTasks = []models.TaskInstance{
		{
			ID:          id,
			Title:       "Test Task",
			Description: "Test Description",
			Completed:   false,
		},
	}
	mockStorage.CompletedTasks = []models.TaskInstance{}

	handlers.SetStorage(mockStorage)
}

func TestCreateTaskRoute(t *testing.T) {
	setupTest()
	router := routes.SetupRoutes()

	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"name":"Test Task"}`))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK && rr.Code != http.StatusCreated {
		t.Errorf("expected status 200 or 201, got %d", rr.Code)
	}
}

func TestGetUncompletedTasksRoute(t *testing.T) {
	setupTest()
	router := routes.SetupRoutes()

	req := httptest.NewRequest("GET", "/tasks", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestCompleteTaskRoute(t *testing.T) {
	setupTestWithTask("123")
	router := routes.SetupRoutes()

	req := httptest.NewRequest("PUT", "/tasks/123/complete", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestCORSOptions(t *testing.T) {
	setupTest()
	router := routes.SetupRoutes()

	req := httptest.NewRequest("OPTIONS", "/tasks", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected 204 for OPTIONS, got %d", rr.Code)
	}
}

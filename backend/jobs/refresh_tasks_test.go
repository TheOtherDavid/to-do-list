package jobs_test

import (
	"testing"
	"time"

	"github.com/TheOtherDavid/to-do-list/jobs"
	"github.com/TheOtherDavid/to-do-list/models"
	"github.com/TheOtherDavid/to-do-list/storage/mocks"
)

func TestRefreshTasks_NoTemplates(t *testing.T) {
	mockStorage := mocks.NewMockStorage()
	mockStorage.Templates = []models.TaskTemplate{}

	jobs.SetStorage(mockStorage)
	err := jobs.RefreshTasks()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mockStorage.SavedTask != nil {
		t.Errorf("Expected no tasks to be created, but one was created")
	}
}

func TestRefreshTasks_TemplateWithNoInstances(t *testing.T) {
	mockStorage := mocks.NewMockStorage()
	mockStorage.Templates = []models.TaskTemplate{
		{
			ID:             "template1",
			Title:          "Test Template",
			Description:    "Test Description",
			RecurrenceRule: "DAILY",
		},
	}

	jobs.SetStorage(mockStorage)
	err := jobs.RefreshTasks()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mockStorage.SavedTask == nil {
		t.Errorf("Expected a task to be created, but none was created")
	}
	if mockStorage.SavedTask.TemplateID != "template1" {
		t.Errorf("Expected template ID to be 'template1', got '%s'", mockStorage.SavedTask.TemplateID)
	}
}

func TestRefreshTasks_TemplateWithIncompleteInstance(t *testing.T) {
	mockStorage := mocks.NewMockStorage()
	mockStorage.Templates = []models.TaskTemplate{
		{
			ID:             "template1",
			Title:          "Test Template",
			Description:    "Test Description",
			RecurrenceRule: "DAILY",
		},
	}

	incompleteInstance := &models.TaskInstance{
		ID:          "instance1",
		TemplateID:  "template1",
		Title:       "Test Task",
		Description: "Test Description",
		Completed:   false,
		CreatedAt:   time.Now().Add(-24 * time.Hour),
	}
	mockStorage.LastInstanceForTemplate = incompleteInstance

	jobs.SetStorage(mockStorage)
	err := jobs.RefreshTasks()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mockStorage.SavedTask != nil {
		t.Errorf("Expected no task to be created, but one was created")
	}
}

func TestRefreshTasks_TemplateWithCompletedInstanceNotEnoughTimePassed(t *testing.T) {
	mockStorage := mocks.NewMockStorage()
	mockStorage.Templates = []models.TaskTemplate{
		{
			ID:             "template1",
			Title:          "Test Template",
			Description:    "Test Description",
			RecurrenceRule: "DAILY",
		},
	}

	// Set up the mock to return a completed instance when GetLastInstanceForTemplate is called
	// but not enough time has passed (completed 12 hours ago, but rule is DAILY)
	completedTime := time.Now().Add(-12 * time.Hour)
	completedInstance := &models.TaskInstance{
		ID:          "instance1",
		TemplateID:  "template1",
		Title:       "Test Task",
		Description: "Test Description",
		Completed:   true,
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		CompletedAt: &completedTime,
	}
	mockStorage.LastInstanceForTemplate = completedInstance

	jobs.SetStorage(mockStorage)
	err := jobs.RefreshTasks()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mockStorage.SavedTask != nil {
		t.Errorf("Expected no task to be created, but one was created")
	}
}

func TestRefreshTasks_TemplateWithCompletedInstanceEnoughTimePassed(t *testing.T) {
	mockStorage := mocks.NewMockStorage()
	mockStorage.Templates = []models.TaskTemplate{
		{
			ID:             "template1",
			Title:          "Test Template",
			Description:    "Test Description",
			RecurrenceRule: "DAILY",
		},
	}

	// Set up the mock to return a completed instance when GetLastInstanceForTemplate is called
	// and enough time has passed (completed 25 hours ago, and rule is DAILY)
	completedTime := time.Now().Add(-25 * time.Hour)
	completedInstance := &models.TaskInstance{
		ID:          "instance1",
		TemplateID:  "template1",
		Title:       "Test Task",
		Description: "Test Description",
		Completed:   true,
		CreatedAt:   time.Now().Add(-48 * time.Hour),
		CompletedAt: &completedTime,
	}
	mockStorage.LastInstanceForTemplate = completedInstance

	jobs.SetStorage(mockStorage)
	err := jobs.RefreshTasks()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mockStorage.SavedTask == nil {
		t.Errorf("Expected a task to be created, but none was created")
	}

	if mockStorage.SavedTask.TemplateID != "template1" {
		t.Errorf("Expected template ID to be 'template1', got '%s'", mockStorage.SavedTask.TemplateID)
	}
}

func TestRefreshTasks_MultipleTemplates(t *testing.T) {
	mockStorage := mocks.NewMockStorage()
	mockStorage.Templates = []models.TaskTemplate{
		{
			ID:             "template1",
			Title:          "Daily Template",
			Description:    "Daily Task",
			RecurrenceRule: "DAILY",
		},
		{
			ID:             "template2",
			Title:          "Weekly Template",
			Description:    "Weekly Task",
			RecurrenceRule: "WEEKLY",
		},
	}

	// For this test, we need to customize the GetLastInstanceForTemplate method
	// to return different results based on the template ID
	// We'll use a custom implementation of the mock storage for this
	customMockStorage := &mocks.MockStorage{
		Templates: mockStorage.Templates,
		GetLastInstanceForTemplateFunc: func(templateID string) (*models.TaskInstance, error) {
			if templateID == "template1" {
				// For template1, return a completed instance with enough time passed
				completedTime := time.Now().Add(-25 * time.Hour)
				return &models.TaskInstance{
					ID:          "instance1",
					TemplateID:  "template1",
					Completed:   true,
					CreatedAt:   time.Now().Add(-48 * time.Hour),
					CompletedAt: &completedTime,
				}, nil
			} else if templateID == "template2" {
				// For template2, return an incomplete instance
				return &models.TaskInstance{
					ID:         "instance2",
					TemplateID: "template2",
					Completed:  false,
					CreatedAt:  time.Now().Add(-24 * time.Hour),
				}, nil
			}
			return nil, nil
		},
	}

	jobs.SetStorage(customMockStorage)
	err := jobs.RefreshTasks()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if customMockStorage.SavedTask == nil {
		t.Errorf("Expected a task to be created, but none was created")
	}

	if customMockStorage.SavedTask.TemplateID != "template1" {
		t.Errorf("Expected template ID to be 'template1', got '%s'", customMockStorage.SavedTask.TemplateID)
	}
}

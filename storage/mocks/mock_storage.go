package mocks

import (
	"github.com/TheOtherDavid/to-do-list/models"
)

// MockStorage implements storage.Storage interface
type MockStorage struct {
	SavedTask        *models.TaskInstance
	UncompletedTasks []models.TaskInstance
	CompletedTasks   []models.TaskInstance
	LastCompletedID  string
	Err             error
	SetCompletedErr error
}

func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

// Task Instance methods
func (m *MockStorage) SaveTaskInstance(task models.TaskInstance) error {
	if m.Err != nil {
		return m.Err
	}
	m.SavedTask = &task
	return nil
}

func (m *MockStorage) GetLastInstanceForTemplate(templateID string) (models.TaskInstance, error) {
	return models.TaskInstance{}, nil
}

func (m *MockStorage) GetUncompletedTaskInstances() ([]models.TaskInstance, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.UncompletedTasks, nil
}

func (m *MockStorage) GetCompletedTaskInstances(limit, offset int) ([]models.TaskInstance, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.CompletedTasks, nil
}

func (m *MockStorage) SetCompleted(id string, completed bool) error {
	if m.SetCompletedErr != nil {
		return m.SetCompletedErr
	}
	if m.Err != nil {
		return m.Err
	}
	m.LastCompletedID = id
	return nil
}

// Task Template methods
func (m *MockStorage) SaveTaskTemplate(template models.TaskTemplate) error {
	return nil
}

func (m *MockStorage) GetAllTaskTemplates() ([]models.TaskTemplate, error) {
	return nil, nil
}

func (m *MockStorage) GetTaskTemplate(id string) (models.TaskTemplate, error) {
	return models.TaskTemplate{}, nil
}

func (m *MockStorage) UpdateTaskTemplate(template models.TaskTemplate) error {
	return nil
}

func (m *MockStorage) DeleteTaskTemplate(id string) error {
	return nil
}

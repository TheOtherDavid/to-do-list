package mocks

import (
	"github.com/TheOtherDavid/to-do-list/models"
)

// MockStorage implements storage.Storage interface
type MockStorage struct {
	// Task Instance fields
	SavedTask               *models.TaskInstance
	UncompletedTasks        []models.TaskInstance
	CompletedTasks          []models.TaskInstance
	LastCompletedID         string
	LastInstanceForTemplate *models.TaskInstance
	AllTaskInstances        []models.TaskInstance

	// Task Template fields
	SavedTemplate       *models.TaskTemplate
	Templates           []models.TaskTemplate
	LastUpdatedTemplate *models.TaskTemplate
	LastDeletedID       string

	// Error fields
	Err             error
	SetCompletedErr error
	TemplateErr     error

	// Custom function fields for advanced mocking
	GetLastInstanceForTemplateFunc func(templateID string) (*models.TaskInstance, error)
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

func (m *MockStorage) GetLastInstanceForTemplate(templateID string) (*models.TaskInstance, error) {
	// If a custom function is provided, use it
	if m.GetLastInstanceForTemplateFunc != nil {
		return m.GetLastInstanceForTemplateFunc(templateID)
	}

	// Otherwise, return the default mock value
	if m.Err != nil {
		return nil, m.Err
	}
	return m.LastInstanceForTemplate, nil
}

func (m *MockStorage) GetAllTaskInstances() ([]models.TaskInstance, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.AllTaskInstances, nil
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
	if m.TemplateErr != nil {
		return m.TemplateErr
	}
	m.SavedTemplate = &template
	return nil
}

func (m *MockStorage) GetAllTaskTemplates() ([]models.TaskTemplate, error) {
	if m.TemplateErr != nil {
		return nil, m.TemplateErr
	}
	return m.Templates, nil
}

func (m *MockStorage) GetTaskTemplate(id string) (models.TaskTemplate, error) {
	if m.TemplateErr != nil {
		return models.TaskTemplate{}, m.TemplateErr
	}

	for _, template := range m.Templates {
		if template.ID == id {
			return template, nil
		}
	}

	return models.TaskTemplate{}, m.TemplateErr
}

func (m *MockStorage) UpdateTaskTemplate(template models.TaskTemplate) error {
	if m.TemplateErr != nil {
		return m.TemplateErr
	}
	m.LastUpdatedTemplate = &template
	return nil
}

func (m *MockStorage) DeleteTaskTemplate(id string) error {
	if m.TemplateErr != nil {
		return m.TemplateErr
	}
	m.LastDeletedID = id
	return nil
}

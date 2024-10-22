package storage

import (
	"github.com/TheOtherDavid/to-do-list/models"
)

type Storage interface {
	SaveTaskInstance(instance models.TaskInstance) error
	GetLastInstanceForTemplate(templateID string) (models.TaskInstance, error)
	GetUncompletedTaskInstances() ([]models.TaskInstance, error)
	GetCompletedTaskInstances(limit, offset int) ([]models.TaskInstance, error)
	SetCompleted(id string, completed bool) error

	SaveTaskTemplate(template models.TaskTemplate) error
	GetAllTaskTemplates() ([]models.TaskTemplate, error)
	GetTaskTemplate(id string) (models.TaskTemplate, error)
	UpdateTaskTemplate(template models.TaskTemplate) error
	DeleteTaskTemplate(id string) error
}

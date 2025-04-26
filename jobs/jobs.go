package jobs

import (
	"time"

	"github.com/TheOtherDavid/to-do-list/models"
	"github.com/TheOtherDavid/to-do-list/storage"
	"github.com/google/uuid"
)

var csvStorage *storage.CSVStorage

func InitStorage(taskInstanceFile, taskTemplateFile string) {
	csvStorage = storage.NewCSVStorage(taskInstanceFile, taskTemplateFile)
}

func RefreshTasks() error {
	templates, err := csvStorage.GetAllTaskTemplates()
	if err != nil {
		return err
	}

	for _, template := range templates {
		lastInstance, err := csvStorage.GetLastInstanceForTemplate(template.ID)
		if err != nil {
			return err
		}

		if needsNewInstance(template, lastInstance) {
			newInstance := createNewInstance(template)
			err = csvStorage.SaveTaskInstance(newInstance)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func needsNewInstance(template models.TaskTemplate, lastInstance models.TaskInstance) bool {
	if lastInstance.ID == "" {
		return true // No instances yet, create first one
	}

	timeSinceLastInstance := time.Since(lastInstance.CreatedAt)

	switch template.RecurrenceRule {
	case "DAILY":
		return timeSinceLastInstance >= 24*time.Hour
	case "WEEKLY":
		return timeSinceLastInstance >= 7*24*time.Hour
	case "FORTNIGHTLY":
		return timeSinceLastInstance >= 14*24*time.Hour
	case "MONTHLY":
		return timeSinceLastInstance >= 30*24*time.Hour // I don't really care about varying month lengths
	default:
		return false
	}
}

func createNewInstance(template models.TaskTemplate) models.TaskInstance {
	return models.TaskInstance{
		ID:          uuid.New().String(),
		TemplateID:  template.ID,
		Title:       template.Title,
		Description: template.Description,
		Completed:   false,
		CreatedAt:   time.Now(),
	}
}

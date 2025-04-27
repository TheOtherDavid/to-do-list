package jobs

import (
	"log"
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
	log.Println("Starting task refresh process")
	templates, err := csvStorage.GetAllTaskTemplates()
	if err != nil {
		log.Printf("Error getting templates: %v", err)
		return err
	}

	log.Printf("Found %d templates to process", len(templates))
	for _, template := range templates {
		log.Printf("Processing template: %s - %s (Rule: %s)", template.ID, template.Title, template.RecurrenceRule)
		lastInstance, err := csvStorage.GetLastInstanceForTemplate(template.ID)
		if err != nil {
			log.Printf("Error getting last instance for template %s: %v", template.ID, err)
			return err
		}

		if lastInstance != nil {
			completionStatus := "incomplete"
			if lastInstance.Completed {
				completionStatus = "completed"
			}
			log.Printf("Last instance: %s (%s)", lastInstance.ID, completionStatus)
		}

		createNew := needsNewInstance(template, lastInstance)
		if createNew {
			newInstance := createNewInstance(template)
			log.Printf("Creating new instance %s for template %s", newInstance.ID, template.ID)
			err = csvStorage.SaveTaskInstance(newInstance)
			if err != nil {
				log.Printf("Error saving new instance: %v", err)
				return err
			}
			log.Printf("Successfully created new instance")
		} else {
			log.Printf("Skipping template %s - no new instance needed", template.ID)
		}
	}

	log.Println("Task refresh completed")
	return nil
}

func needsNewInstance(template models.TaskTemplate, lastInstance *models.TaskInstance) bool {
	if lastInstance == nil {
		log.Printf("No previous instances exist - creating first instance")
		return true
	}

	if !lastInstance.Completed {
		log.Printf("Last instance is not yet completed - skipping")
		return false
	}

	timeSinceLastInstance := time.Since(*lastInstance.CompletedAt)
	var requiredTime time.Duration
	var createNew bool

	switch template.RecurrenceRule {
	case "DAILY":
		requiredTime = 24 * time.Hour
		createNew = timeSinceLastInstance >= requiredTime
	case "WEEKLY":
		requiredTime = 7 * 24 * time.Hour
		createNew = timeSinceLastInstance >= requiredTime
	case "FORTNIGHTLY":
		requiredTime = 14 * 24 * time.Hour
		createNew = timeSinceLastInstance >= requiredTime
	case "MONTHLY":
		requiredTime = 30 * 24 * time.Hour
		createNew = timeSinceLastInstance >= requiredTime
	default:
		log.Printf("Unknown recurrence rule: %s", template.RecurrenceRule)
		return false
	}

	if createNew {
		log.Printf("Time since last completion: %v (required: %v) - creating new instance",
			timeSinceLastInstance.Round(time.Hour), requiredTime)
	} else {
		log.Printf("Time since last completion: %v (required: %v) - not enough time has passed",
			timeSinceLastInstance.Round(time.Hour), requiredTime)
	}

	return createNew
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

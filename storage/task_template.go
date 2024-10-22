package storage

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/TheOtherDavid/to-do-list/models"
)

func (s *CSVStorage) SaveTaskTemplate(template models.TaskTemplate) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	file, err := os.OpenFile(s.templateFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.Write([]string{
		template.ID,
		template.Title,
		template.Description,
		template.RecurrenceRule,
		template.CreatedAt.Format(time.RFC3339),
	})
}

func (s *CSVStorage) GetAllTaskTemplates() ([]models.TaskTemplate, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	file, err := os.Open(s.templateFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var templates []models.TaskTemplate
	for _, record := range records {
		createdAt, err := time.Parse(time.RFC3339, record[4])
		if err != nil {
			return nil, fmt.Errorf("error parsing CreatedAt time for template %s: %v", record[0], err)
		}

		template := models.TaskTemplate{
			ID:             record[0],
			Title:          record[1],
			Description:    record[2],
			RecurrenceRule: record[3],
			CreatedAt:      createdAt,
		}
		templates = append(templates, template)
	}

	return templates, nil
}

func (s *CSVStorage) GetTaskTemplate(id string) (models.TaskTemplate, error) {
	templates, err := s.GetAllTaskTemplates()
	if err != nil {
		return models.TaskTemplate{}, err
	}

	for _, template := range templates {
		if template.ID == id {
			return template, nil
		}
	}

	return models.TaskTemplate{}, fmt.Errorf("task template with ID %s not found", id)
}

func (s *CSVStorage) UpdateTaskTemplate(template models.TaskTemplate) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	templates, err := s.GetAllTaskTemplates()
	if err != nil {
		return err
	}

	found := false
	for i, t := range templates {
		if t.ID == template.ID {
			templates[i] = template
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task template with ID %s not found", template.ID)
	}

	// Rewrite the entire file with the updated data
	file, err := os.Create(s.templateFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, t := range templates {
		err := writer.Write([]string{
			t.ID,
			t.Title,
			t.Description,
			t.RecurrenceRule,
			t.CreatedAt.Format(time.RFC3339),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *CSVStorage) DeleteTaskTemplate(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	templates, err := s.GetAllTaskTemplates()
	if err != nil {
		return err
	}

	found := false
	var updatedTemplates []models.TaskTemplate
	for _, t := range templates {
		if t.ID != id {
			updatedTemplates = append(updatedTemplates, t)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("task template with ID %s not found", id)
	}

	// Rewrite the entire file with the updated data
	file, err := os.Create(s.templateFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, t := range updatedTemplates {
		err := writer.Write([]string{
			t.ID,
			t.Title,
			t.Description,
			t.RecurrenceRule,
			t.CreatedAt.Format(time.RFC3339),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

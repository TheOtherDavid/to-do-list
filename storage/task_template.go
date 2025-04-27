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

	file, err := os.OpenFile(s.templateFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		template.ID,
		template.Title,
		template.Description,
		template.RecurrenceRule,
		template.CreatedAt.Format(time.RFC3339),
	}

	if err := writer.Write(record); err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}

	return nil
}

func (s *CSVStorage) GetAllTaskTemplates() ([]models.TaskTemplate, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	file, err := os.Open(s.templateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read records: %w", err)
	}

	templates := make([]models.TaskTemplate, 0, len(records))
	for _, record := range records {
		createdAt, err := time.Parse(time.RFC3339, record[4])
		if err != nil {
			return nil, fmt.Errorf("failed to parse created at: %w", err)
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
	s.mutex.Lock()
	defer s.mutex.Unlock()

	file, err := os.Open(s.templateFile)
	if err != nil {
		return models.TaskTemplate{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return models.TaskTemplate{}, fmt.Errorf("failed to read records: %w", err)
	}

	for _, record := range records {
		if record[0] == id {
			createdAt, err := time.Parse(time.RFC3339, record[4])
			if err != nil {
				return models.TaskTemplate{}, fmt.Errorf("failed to parse created at: %w", err)
			}

			return models.TaskTemplate{
				ID:             record[0],
				Title:          record[1],
				Description:    record[2],
				RecurrenceRule: record[3],
				CreatedAt:      createdAt,
			}, nil
		}
	}

	return models.TaskTemplate{}, fmt.Errorf("template with ID %s not found", id)
}

func (s *CSVStorage) UpdateTaskTemplate(template models.TaskTemplate) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	file, err := os.Open(s.templateFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read records: %w", err)
	}

	found := false
	for i, record := range records {
		if record[0] == template.ID {
			found = true
			records[i] = []string{
				template.ID,
				template.Title,
				template.Description,
				template.RecurrenceRule,
				template.CreatedAt.Format(time.RFC3339),
			}
			break
		}
	}

	if !found {
		return fmt.Errorf("template with ID %s not found", template.ID)
	}

	file, err = os.OpenFile(s.templateFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf("failed to write records: %w", err)
	}

	return nil
}

func (s *CSVStorage) DeleteTaskTemplate(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	file, err := os.Open(s.templateFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read records: %w", err)
	}

	found := false
	var updatedRecords [][]string
	for _, record := range records {
		if record[0] != id {
			updatedRecords = append(updatedRecords, record)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("template with ID %s not found", id)
	}

	file, err = os.OpenFile(s.templateFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.WriteAll(updatedRecords); err != nil {
		return fmt.Errorf("failed to write records: %w", err)
	}

	return nil
}

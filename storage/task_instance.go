package storage

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/TheOtherDavid/to-do-list/models"
)

func (s *CSVStorage) SaveTaskInstance(instance models.TaskInstance) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	file, err := os.OpenFile(s.instanceFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
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

	completedAt := ""
	if instance.CompletedAt != nil {
		completedAt = instance.CompletedAt.Format(time.RFC3339)
	}

	record := []string{
		instance.ID,
		instance.TemplateID,
		instance.Title,
		instance.Description,
		strconv.FormatBool(instance.Completed),
		instance.CreatedAt.Format(time.RFC3339),
		completedAt,
	}

	if err := writer.Write(record); err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}

	return err
}

func (s *CSVStorage) GetLastInstanceForTemplate(templateID string) (*models.TaskInstance, error) {
	instances, err := s.getAllTaskInstances()
	if err != nil {
		return nil, fmt.Errorf("failed to get all task instances: %w", err)
	}

	var lastInstance models.TaskInstance
	for _, instance := range instances {
		if instance.TemplateID == templateID {
			if lastInstance.CreatedAt.Before(instance.CreatedAt) {
				lastInstance = instance
			}
		}
	}

	if lastInstance.ID == "" {
		return nil, nil
	}

	return &lastInstance, nil
}

func (s *CSVStorage) GetUncompletedTaskInstances() ([]models.TaskInstance, error) {
	instances, err := s.getAllTaskInstances()
	if err != nil {
		return nil, fmt.Errorf("failed to get all task instances: %w", err)
	}

	var uncompletedInstances []models.TaskInstance
	for _, instance := range instances {
		if !instance.Completed {
			uncompletedInstances = append(uncompletedInstances, instance)
		}
	}

	return uncompletedInstances, nil
}

func (s *CSVStorage) GetCompletedTaskInstances(limit, offset int) ([]models.TaskInstance, error) {
	instances, err := s.getAllTaskInstances()
	if err != nil {
		return nil, fmt.Errorf("failed to get all task instances: %w", err)
	}

	var completedInstances []models.TaskInstance
	for _, instance := range instances {
		if instance.Completed {
			completedInstances = append(completedInstances, instance)
		}
	}

	sort.Slice(completedInstances, func(i, j int) bool {
		return completedInstances[i].CompletedAt.After(*completedInstances[j].CompletedAt)
	})

	// Apply pagination
	start := offset
	end := offset + limit
	if start > len(completedInstances) {
		return []models.TaskInstance{}, nil
	}
	if end > len(completedInstances) {
		end = len(completedInstances)
	}

	return completedInstances[start:end], nil
}

func (s *CSVStorage) getAllTaskInstances() ([]models.TaskInstance, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	file, err := os.Open(s.instanceFile)
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

	var instances []models.TaskInstance
	for _, record := range records {
		completed, _ := strconv.ParseBool(record[4])
		createdAt, _ := time.Parse(time.RFC3339, record[5])

		var completedAt *time.Time
		if record[6] != "" {
			t, _ := time.Parse(time.RFC3339, record[6])
			completedAt = &t
		}

		instance := models.TaskInstance{
			ID:          record[0],
			TemplateID:  record[1],
			Title:       record[2],
			Description: record[3],
			Completed:   completed,
			CreatedAt:   createdAt,
			CompletedAt: completedAt,
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

func (s *CSVStorage) SetCompleted(id string, completed bool) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	instances, err := s.getAllTaskInstances()
	if err != nil {
		return fmt.Errorf("failed to get all task instances: %w", err)
	}

	found := false
	for i, instance := range instances {
		if instance.ID == id {
			instances[i].Completed = completed
			if completed {
				now := time.Now()
				instances[i].CompletedAt = &now
			} else {
				instances[i].CompletedAt = nil
			}
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task instance with ID %s not found", id)
	}

	file, err := os.Create(s.instanceFile)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, instance := range instances {
		completedAt := ""
		if instance.CompletedAt != nil {
			completedAt = instance.CompletedAt.Format(time.RFC3339)
		}

		record := []string{
			instance.ID,
			instance.TemplateID,
			instance.Title,
			instance.Description,
			strconv.FormatBool(instance.Completed),
			instance.CreatedAt.Format(time.RFC3339),
			completedAt,
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return err
}

package storage

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/TheOtherDavid/to-do-list/models"
)

type CSVStorage struct {
	filename string
	mutex    sync.RWMutex
}

func NewCSVStorage(filename string) *CSVStorage {
	return &CSVStorage{
		filename: filename,
	}
}

func (s *CSVStorage) CreateTask(task models.Task) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	file, err := os.OpenFile(s.filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	completedAt := ""
	if task.CompletedAt != nil {
		completedAt = task.CompletedAt.Format(time.RFC3339)
	}

	return writer.Write([]string{
		task.ID,
		task.Title,
		task.Description,
		fmt.Sprintf("%v", task.Completed),
		task.CreatedAt.Format(time.RFC3339),
		completedAt,
	})
}

func (s *CSVStorage) ListTasks() ([]models.Task, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	file, err := os.Open(s.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var tasks []models.Task
	for _, record := range records {
		task := models.Task{
			ID:          record[0],
			Title:       record[1],
			Description: record[2],
			Completed:   record[3] == "true",
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *CSVStorage) DeleteTask(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	tasks, err := s.ListTasks()
	if err != nil {
		return err
	}

	var newTasks []models.Task
	for _, task := range tasks {
		if task.ID != id {
			newTasks = append(newTasks, task)
		}
	}

	file, err := os.Create(s.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, task := range newTasks {
		err := writer.Write([]string{task.ID, task.Title, task.Description, fmt.Sprintf("%v", task.Completed)})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *CSVStorage) UpdateTasks(tasks []models.Task) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	file, err := os.Create(s.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, task := range tasks {
		err := writer.Write([]string{task.ID, task.Title, task.Description, fmt.Sprintf("%v", task.Completed)})
		if err != nil {
			return err
		}
	}

	return nil
}

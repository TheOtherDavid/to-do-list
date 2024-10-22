package storage

import (
	"sync"
)

type CSVStorage struct {
	templateFile string
	instanceFile string
	mutex        sync.RWMutex
}

func NewCSVStorage(templateFile, instanceFile string) *CSVStorage {
	return &CSVStorage{
		templateFile: templateFile,
		instanceFile: instanceFile,
	}
}

// Helper methods for reading/writing CSV files

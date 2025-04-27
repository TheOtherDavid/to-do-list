package storage

import (
	"sync"
)

type CSVStorage struct {
	templateFile string
	instanceFile string
	mutex        sync.RWMutex
}

func NewCSVStorage(instanceFile, templateFile string) *CSVStorage {
	return &CSVStorage{
		instanceFile: instanceFile,
		templateFile: templateFile,
	}
}

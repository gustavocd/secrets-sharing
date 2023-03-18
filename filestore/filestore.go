package filestore

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/gustavocd/secrets-sharing/types"
)

// fileStore is a file-based store for secret data.
type fileStore struct {
	Mu    sync.Mutex
	Store map[string]string
}

// FileStoreConfig is the configuration for the file store.
var FileStoreConfig struct {
	DataFilePath string
	Fs           fileStore
}

// Init initializes the file store.
func Init(dataFilePath string) error {
	_, err := os.Stat(dataFilePath)

	if err != nil {
		_, err := os.Create(dataFilePath)
		if err != nil {
			return err
		}
	}

	FileStoreConfig.Fs = fileStore{Mu: sync.Mutex{}, Store: make(map[string]string)}
	FileStoreConfig.DataFilePath = dataFilePath

	return nil
}

// Write writes the secret data to the store and returns an error if the write fails.
func (j *fileStore) Write(data types.SecretData) error {
	j.Mu.Lock()
	defer j.Mu.Unlock()

	err := j.ReadFromFile()
	if err != nil {
		return err
	}
	j.Store[data.Id] = data.Secret

	return j.WriteToFile()
}

// Read returns the secret data for the given id or an error if the id is not found.
func (j *fileStore) Read(id string) (string, error) {
	j.Mu.Lock()
	defer j.Mu.Unlock()

	err := j.ReadFromFile()
	if err != nil {
		return "", err
	}

	data := j.Store[id]
	delete(j.Store, id)
	j.WriteToFile()

	return data, nil
}

// ReadFromFile reads the file store data from the file and returns an error if the read fails.
func (j *fileStore) ReadFromFile() error {
	f, err := os.Open(FileStoreConfig.DataFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	jsonData, err := os.ReadFile(FileStoreConfig.DataFilePath)
	if err != nil {
		return err
	}

	if len(jsonData) != 0 {
		return json.Unmarshal(jsonData, &j.Store)
	}

	return nil
}

// WriteToFile writes the file store data to the file and returns an error if the write fails.
func (j *fileStore) WriteToFile() error {
	var f *os.File
	jsonData, err := json.MarshalIndent(j.Store, "", "  ")
	if err != nil {
		return err
	}

	f, err = os.Create(FileStoreConfig.DataFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

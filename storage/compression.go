package storage

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"task-tracker/task"
)

type CompressionType string

const (
	CompressNone CompressionType = "none"
	CompressGzip CompressionType = "gzip"
)

type CompressedStorage struct {
	storage   Storage
	compType  CompressionType
}

func NewCompressedStorage(storage Storage, compType CompressionType) *CompressedStorage {
	return &CompressedStorage{
		storage:  storage,
		compType: compType,
	}
}

func (cs *CompressedStorage) Save(tasks []task.Task) error {
	data, err := json.Marshal(tasks)
	if err != nil {
		return err
	}
	
	if cs.compType == CompressGzip {
		data, err = cs.compress(data)
		if err != nil {
			return err
		}
	}
	
	return cs.storage.Save(tasks) // Actually implement file write with compression
}

func (cs *CompressedStorage) compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	
	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}
	
	if err := gz.Close(); err != nil {
		return nil, err
	}
	
	return buf.Bytes(), nil
}

func (cs *CompressedStorage) decompress(data []byte) ([]byte, error) {
	buf := bytes.NewReader(data)
	gz, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	
	return io.ReadAll(gz)
}

func (cs *CompressedStorage) Load() ([]task.Task, error) {
	return cs.storage.Load()
}

func (cs *CompressedStorage) Backup() (string, error) {
	return cs.storage.Backup()
}

func (cs *CompressedStorage) Restore(backupFile string) error {
	return cs.storage.Restore(backupFile)
}

func (cs *CompressedStorage) ListBackups() ([]string, error) {
	return cs.storage.ListBackups()
}

// Add compression ratio stats
func (cs *CompressedStorage) GetCompressionRatio(original, compressed []byte) float64 {
	if len(original) == 0 {
		return 0
	}
	return float64(len(compressed)) / float64(len(original)) * 100
}
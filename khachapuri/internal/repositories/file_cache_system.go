package repositories

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/matzapata/ipfs-compute/provider/pkg/system"
)

type FileSystemCache struct {
	BaseDir string
}

func NewFileSystemCache(dir string) *FileSystemCache {
	return &FileSystemCache{
		BaseDir: dir,
	}
}

func (r *FileSystemCache) Exists(key string) bool {
	path := r.buildPath(key)
	return system.Exists(path)
}

func (r *FileSystemCache) Get(key string) ([]byte, error) {
	path := r.buildPath(key)
	if !system.Exists(path) {
		return nil, errors.New("file not found in cache")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *FileSystemCache) Set(key string, data []byte) error {
	path := r.buildPath(key)
	if err := os.WriteFile(path, data, 0664); err != nil {
		return err
	}

	return nil
}

func (r *FileSystemCache) buildPath(key string) string {
	return filepath.Join(r.BaseDir, key)
}

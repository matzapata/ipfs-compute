package system

import (
	"os"
	"path/filepath"
)

func BuildCwdPath(path string) string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, path)
}

func BuildLocalPath(path string) string {
	exe := os.Args[0]
	exeDir := filepath.Dir(exe)
	return filepath.Join(exeDir, path)
}

func EnsureDirExists(path string, empty bool) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}

	if empty {
		return EnsureDirEmpty(path)
	}

	return nil
}

func EnsureDirEmpty(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	return os.RemoveAll(path)
}

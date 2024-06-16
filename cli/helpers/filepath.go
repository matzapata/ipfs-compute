package helpers

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

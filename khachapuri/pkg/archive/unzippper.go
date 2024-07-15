package archive

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type IUnzipper interface {
	UnzipBytes(data []byte, dest string) error
	UnzipFilepath(path string, dest string) error
}

type Unzipper struct {
}

func NewUnzipper() *Unzipper {
	return &Unzipper{}
}

func (z *Unzipper) UnzipBytes(data []byte, dest string) error {
	reader := bytes.NewReader(data)
	return unzip(reader, int64(len(data)), dest)
}

func (z *Unzipper) UnzipFilepath(zipPath string, dest string) error {
	file, err := os.Open(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	return unzip(file, fileInfo.Size(), dest)
}

func unzip(reader io.ReaderAt, size int64, dest string) error {
	r, err := zip.NewReader(reader, size)
	if err != nil {
		return fmt.Errorf("failed to create zip reader: %w", err)
	}

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		// Create directory structure if necessary
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Create a file in the destination directory
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)
		// Close the file and the reader
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

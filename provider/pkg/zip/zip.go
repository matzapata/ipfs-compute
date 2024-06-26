package zip_service

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type IZipService interface {
	Unzip(src string) (string, error)
	ZipFolder(srcFolder string, destZip string) error
}

type ZipService struct {
}

func NewZipService() *ZipService {
	return &ZipService{}
}

func (z *ZipService) Unzip(src string) (string, error) {
	dest, err := os.CreateTemp("", "*")
	if err != nil {
		return "", err
	}

	r, err := zip.OpenReader(src)
	if err != nil {
		return "", fmt.Errorf("failed to open zip file: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest.Name(), f.Name)

		// Create directory structure if necessary
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Create a file in the destination directory
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return "", err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", err
		}

		rc, err := f.Open()
		if err != nil {
			return "", err
		}

		_, err = io.Copy(outFile, rc)
		// Close the file and the reader
		outFile.Close()
		rc.Close()

		if err != nil {
			return "", err
		}
	}
	return dest.Name(), nil
}

// ZipFolder recursively zips a folder and its contents
func (z *ZipService) ZipFolder(srcFolder string, destZip string) error {
	// Create the zip file
	zipFile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Traverse the folder
	err = filepath.Walk(srcFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the source folder itself
		if path == srcFolder {
			return nil
		}

		// If it's a directory, create it in the zip file
		if info.IsDir() {
			relPath, err := filepath.Rel(srcFolder, path)
			if err != nil {
				return err
			}
			_, err = zipWriter.Create(relPath + "/")
			return err
		}

		// If it's a file, add it to the zip file
		return AddFileToZip(zipWriter, path, srcFolder)
	})

	return err
}

// AddFileToZip adds a file to the zip archive
func AddFileToZip(zipWriter *zip.Writer, filePath string, basePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get the file's relative path
	relPath, err := filepath.Rel(basePath, filePath)
	if err != nil {
		return err
	}

	// Create a zip file header
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	zipFileHeader, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}
	zipFileHeader.Name = relPath

	// Create a writer for the zip file entry
	zipWriterEntry, err := zipWriter.CreateHeader(zipFileHeader)
	if err != nil {
		return err
	}

	// Copy the file content to the zip file
	_, err = io.Copy(zipWriterEntry, file)
	return err
}

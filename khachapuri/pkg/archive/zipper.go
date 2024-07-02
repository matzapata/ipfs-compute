package archive

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

type IZipper interface {
	ZipFolder(srcFolder string, dest string) error
}

type Zipper struct {
}

func NewZipper() *Zipper {
	return &Zipper{}
}

// ZipFolder recursively zips a folder and its contents
func (z *Zipper) ZipFolder(srcFolder string, destZip string) error {
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
		return addFileToZip(zipWriter, path, srcFolder)
	})

	return err
}

// AddFileToZip adds a file to the zip archive
func addFileToZip(zipWriter *zip.Writer, filePath string, basePath string) error {
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

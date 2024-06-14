package helpers

import (
	"fmt"
	"os"
)

func WriteFile(data []byte, filePath string) error {
	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Write data to the file
	_, err = out.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %w", err)
	}

	return nil
}

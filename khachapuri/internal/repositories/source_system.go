package repositories

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/matzapata/ipfs-compute/provider/internal/domain"
	"github.com/matzapata/ipfs-compute/provider/pkg/system"
	"gopkg.in/yaml.v2"
)

type SystemSourceRepository struct {
}

func NewSystemSourceRepository() *SystemSourceRepository {
	return &SystemSourceRepository{}
}

// TODO: remove this one
func (r *SystemSourceRepository) GetSource() (*domain.Source, error) {
	return &domain.Source{
		SourcePath: system.BuildCwdPath("."),
		ConfigPath: system.BuildCwdPath("khachapuri.yaml"),
	}, nil
}

func (r *SystemSourceRepository) GetSourceSpecification() (*domain.SourceSpecification, error) {
	source, err := r.GetSource()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(source.ConfigPath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the YAML data into the struct
	var spec domain.SourceSpecification
	err = yaml.Unmarshal(data, &spec)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling specification YAML: %w", err)
	}

	// load EnvFile	and add it to Env
	for serviceName, service := range spec.Services {
		for _, envFile := range service.EnvFile {
			envVars, err := parseEnvFile(envFile)
			if err != nil {
				return nil, err
			}

			service := spec.Services[serviceName]
			service.Env = append(service.Env, envVars...)
			spec.Services[serviceName] = service
		}
	}

	return &spec, nil
}

func parseEnvFile(filePath string) ([]domain.Env, error) {
	file, err := os.Open(system.BuildCwdPath(filePath))
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Slice to hold the results
	var envVars []domain.Env

	// Iterate over each line
	for scanner.Scan() {
		// Split each line at "="
		line := scanner.Text()
		parts := strings.Split(line, "=")

		// Ensure the line is in name=value format
		if len(parts) == 2 {
			// Create a new EnvVar struct and append to envVars slice
			envVar := domain.Env{
				Name:  parts[0],
				Value: parts[1],
			}
			envVars = append(envVars, envVar)
		}
	}

	return envVars, nil

}

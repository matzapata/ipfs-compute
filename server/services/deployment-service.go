package services

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/matzapata/ipfs-compute/helpers"
	"github.com/matzapata/ipfs-compute/repositories"
)

type DeploymentService struct {
	repo repositories.DeploymentsRepository
}

type DeploymentMetadataJson struct {
	Owner string `json:"owner"` // owner address encrypted with our public key
	Env   string `json:"env"`   // environment variables in format "KEY=VALUE;KEY2=VALUE2" encrypted with our public key
}

type DeploymentMetadata struct {
	Owner string
	Env   []string
}

func NewDeploymentService(repo *repositories.DeploymentsRepository) *DeploymentService {
	return &DeploymentService{
		repo: *repo,
	}
}

const (
	MB                    = 1 << 20 // 1 MB in bytes (1 << 20 is 2^20)
	MAX_ZIPPED_DEPLOYMENT = 50 * MB
	DEPLOYMENT_SPEC_FILE  = "deployment.json"
	DEPLOYMENT_BIN_FILE   = "main"
)

func (d *DeploymentService) GetDeployment(cid string, dstDir string) (*DeploymentMetadata, error) {
	// download zipped file
	data, err := d.repo.GetZippedDeployment(cid, MAX_ZIPPED_DEPLOYMENT)
	if err != nil {
		return nil, err
	}

	// write file
	zipFilePath := filepath.Join(dstDir, "deployment.zip")
	err = helpers.WriteFile(data, zipFilePath)
	if err != nil {
		return nil, err
	}

	// unzip
	err = helpers.Unzip(zipFilePath, dstDir)
	if err != nil {
		return nil, err
	}

	// read spec file
	specFilePath := filepath.Join(dstDir, DEPLOYMENT_SPEC_FILE)
	metadata, err := readDeploymentSpecFile(specFilePath)
	if err != nil {
		return nil, err
	}

	// decrypt owner address and env vars

	return &DeploymentMetadata{
		Owner: metadata.Owner,
		Env:   metadata.Env,
	}, nil
}

func readDeploymentSpecFile(deploymentJsonFilePath string) (*DeploymentMetadata, error) {
	// open the JSON file
	file, err := os.Open(deploymentJsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// read the file contents
	jsonData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	// unmarshal the JSON data
	var metadata DeploymentMetadata
	err = json.Unmarshal(jsonData, &metadata)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %v", err)
	}

	return &metadata, nil
}

package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/matzapata/ipfs-compute/helpers"
	"github.com/matzapata/ipfs-compute/repositories"
)

type DeploymentService struct {
	repo repositories.DeploymentsRepository
}

type DeploymentMetadata struct {
	Owner          string   `json:"owner"`           // owner address. If no payer specified we default to this one
	OwnerSignature string   `json:"owner_signature"` // owner signature of the terms and conditions
	Env            []string `json:"env"`             // environment variables in format "KEY=VALUE;KEY2=VALUE2"
	DeploymentCid  string   `json:"deployment_cid"`  // where the deployment zip is stored
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
)

func (d *DeploymentService) GetDeployment(cid string, dstDir string) error {
	// download zipped file
	data, err := d.repo.GetZippedDeployment(cid, MAX_ZIPPED_DEPLOYMENT)
	if err != nil {
		return err
	}

	// write file
	zipFilePath := filepath.Join(dstDir, "deployment.zip")
	err = helpers.WriteFile(data, zipFilePath)
	if err != nil {
		return err
	}

	// unzip
	err = helpers.Unzip(zipFilePath, dstDir)
	if err != nil {
		return err
	}

	return nil
}

func (d *DeploymentService) GetDeploymentMetadata(cid string) (*DeploymentMetadata, error) {
	// open the JSON file
	data, err := d.repo.GetDeploymentSpecFile(cid)
	if err != nil {
		return nil, err
	}

	// decrypt the JSON data
	privateKey, err := helpers.LoadPrivateKeyFromString(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return nil, fmt.Errorf("error loading private key: %v", err)
	}
	decDeploymentMetadataStr, err := helpers.DecryptBytes(privateKey, data)
	if err != nil {
		return nil, fmt.Errorf("error decrypting JSON data: %v", err)
	}

	// unmarshal the JSON data
	var metadata DeploymentMetadata
	err = json.Unmarshal(decDeploymentMetadataStr, &metadata)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling metadata json: %v", err)
	}

	return &metadata, nil
}

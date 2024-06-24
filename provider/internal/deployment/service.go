package deployment

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/matzapata/ipfs-compute/provider/internal/config"
	deployment_repository "github.com/matzapata/ipfs-compute/provider/internal/deployment/repository"
	ecdsa_helpers "github.com/matzapata/ipfs-compute/provider/pkg/helpers/ecdsa"
	files_helpers "github.com/matzapata/ipfs-compute/provider/pkg/helpers/files"
	rsa_helpers "github.com/matzapata/ipfs-compute/provider/pkg/helpers/rsa"

	cp "github.com/otiai10/copy"
)

type DeploymentService struct {
	repo deployment_repository.DeploymentsRepository
}

type DeploymentSpecification struct {
	Env []string `json:"env"`
}

type Deployment struct {
	DeploymentSpecification
	Owner          string `json:"owner"`
	OwnerSignature string `json:"owner_signature"`
	DeploymentCid  string `json:"deployment_cid"`
}

func NewDeploymentService(repo *deployment_repository.DeploymentsRepository) *DeploymentService {
	return &DeploymentService{
		repo: *repo,
	}
}

// TODO: refactor, return something/rename/get metadata
func (d *DeploymentService) GetDeployment(cid string, dstDir string) error {
	// download zipped file
	data, err := d.repo.GetZippedDeployment(cid, config.MAX_ZIPPED_DEPLOYMENT)
	if err != nil {
		return err
	}

	// TODO: this should be out of here already
	// write file
	zipFilePath := filepath.Join(dstDir, "deployment.zip")
	err = files_helpers.WriteFile(data, zipFilePath)
	if err != nil {
		return err
	}

	// unzip
	err = files_helpers.Unzip(zipFilePath, dstDir)
	if err != nil {
		return err
	}

	return nil
}

func (d *DeploymentService) GetDeploymentMetadata(cid string) (*Deployment, error) {
	// open the JSON file
	data, err := d.repo.GetDeploymentSpecFile(cid)
	if err != nil {
		return nil, err
	}

	// decrypt the JSON data
	// TODO: get private key from settings
	privateKey, err := rsa_helpers.LoadPrivateKeyFromString(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return nil, fmt.Errorf("error loading private key: %v", err)
	}
	decDeploymentMetadataStr, err := rsa_helpers.DecryptBytes(privateKey, data)
	if err != nil {
		return nil, fmt.Errorf("error decrypting JSON data: %v", err)
	}

	// unmarshal the JSON data
	var metadata Deployment
	err = json.Unmarshal(decDeploymentMetadataStr, &metadata)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling metadata json: %v", err)
	}

	return &metadata, nil
}

func (*DeploymentService) BuildDeploymentZip() error {
	err := os.Mkdir(config.DIST_DIR, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// copy binary and public folder to dist folder
	err = os.Mkdir(config.DIST_DEPLOYMENT_DIR, 0755)
	if err != nil {
		return err
	}
	err = cp.Copy(config.SRC_BIN_FILE, config.DIST_BIN_FILE)
	if err != nil {
		return err
	}
	err = cp.Copy(config.SRC_PUBLIC_DIR, config.DIST_PUBLIC_DIR)
	if err != nil {
		return err
	}

	// zip dist folder
	err = files_helpers.ZipFolder(config.DIST_DEPLOYMENT_DIR, config.DIST_ZIP_FILE)
	if err != nil {
		return err
	}

	return nil
}

func (*DeploymentService) BuildDeploymentSpecification(deploymentZipCid string, signature *ecdsa_helpers.Signature, providerPublicKey string) error {
	// load public key
	ipfsComputePublicKey, err := rsa_helpers.LoadPublicKeyFromString(providerPublicKey)
	if err != nil {
		log.Fatal(err)
	}

	// read spec json file
	deploymentSpecJson, err := os.ReadFile(config.SRC_SPEC_FILE)
	if err != nil {
		log.Fatal(err)
	}
	var deploymentSpec DeploymentSpecification
	err = json.Unmarshal(deploymentSpecJson, &deploymentSpec)
	if err != nil {
		log.Fatal(err)
	}

	// add signature to json
	deployment := Deployment{
		DeploymentSpecification: deploymentSpec,
		Owner:                   signature.Address,
		OwnerSignature:          signature.Signature,
		DeploymentCid:           deploymentZipCid,
	}

	// encrypt it with public key
	deploymentJson, err := json.Marshal(deployment)
	if err != nil {
		return err
	}
	encDeploymentJson, err := rsa_helpers.EncryptBytes(ipfsComputePublicKey, deploymentJson)
	if err != nil {
		return err
	}

	// save it to a file
	err = os.WriteFile(config.DIST_SPEC_FILE, encDeploymentJson, 0644)
	if err != nil {
		return err
	}

	return nil
}

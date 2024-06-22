package services

import (
	"encoding/json"
	"log"
	"os"

	"github.com/matzapata/ipfs-compute/cli/config"
	"github.com/matzapata/ipfs-compute/cli/helpers"
	"github.com/matzapata/ipfs-compute/shared/cryptoecdsa"
	"github.com/matzapata/ipfs-compute/shared/cryptorsa"
	cp "github.com/otiai10/copy"
)

type DeploymentService struct {
}

func NewDeploymentService() *DeploymentService {
	return &DeploymentService{}
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

const TERMS_AND_CONDITIONS = "Deploy to IPFS Compute. Code will be runnable by everyone with the CID."

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
	err = helpers.ZipFolder(config.DIST_DEPLOYMENT_DIR, config.DIST_ZIP_FILE)
	if err != nil {
		return err
	}

	return nil
}

func (*DeploymentService) BuildDeploymentSpecification(deploymentZipCid string, signature *cryptoecdsa.Signature, providerPublicKey string) error {
	// load public key
	ipfsComputePublicKey, err := cryptorsa.LoadPublicKeyFromString(providerPublicKey)
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
	encDeploymentJson, err := cryptorsa.EncryptBytes(ipfsComputePublicKey, deploymentJson)
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

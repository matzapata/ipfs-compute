package artifact

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"os"

	artifact_repository "github.com/matzapata/ipfs-compute/provider/internal/artifact/repository"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	files_helpers "github.com/matzapata/ipfs-compute/provider/pkg/helpers/files"
	rsa_helpers "github.com/matzapata/ipfs-compute/provider/pkg/helpers/rsa"
)

type ArtifactService struct {
	ArtifactRepository artifact_repository.ArtifactRepository
}

type ArtifactSpecification struct {
	Env []string `json:"env"`
}

type Artifact struct {
	ArtifactSpecification
	Owner          string `json:"owner"`
	OwnerSignature string `json:"owner_signature"`
	DeploymentCid  string `json:"deployment_cid"`
}

func NewArtifactService(artifactRepository artifact_repository.ArtifactRepository) *ArtifactService {
	return &ArtifactService{
		ArtifactRepository: artifactRepository,
	}
}

func (d *ArtifactService) GetArtifactExecutable(cid string) (executablePath string, err error) {
	zippedExecutablePath, err := d.ArtifactRepository.GetZippedExecutable(cid, config.MAX_ZIPPED_DEPLOYMENT)
	if err != nil {
		return
	}
	defer os.Remove(zippedExecutablePath)

	// unzip to destination directory
	unzippedExecutable, err := os.CreateTemp("", "executable-*")
	if err != nil {
		return
	}
	err = files_helpers.Unzip(zippedExecutablePath, unzippedExecutable.Name())
	if err != nil {
		return
	}

	return unzippedExecutable.Name(), nil
}

func (d *ArtifactService) GetArtifactSpecification(cid string, providerRsaPrivateKey *rsa.PrivateKey) (*Artifact, error) {
	encSpecPath, err := d.ArtifactRepository.GetSpecificationFile(cid)
	if err != nil {
		return nil, err
	}

	// TODO: it may or may not be encrypted

	// read the JSON encSpecData
	encSpecData, err := os.ReadFile(encSpecPath)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON data: %v", err)
	}
	artifactSpecification, err := rsa_helpers.DecryptBytes(providerRsaPrivateKey, encSpecData)
	if err != nil {
		return nil, fmt.Errorf("error decrypting JSON data: %v", err)
	}

	// unmarshal the JSON data
	var artifact Artifact
	err = json.Unmarshal(artifactSpecification, &artifact)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling metadata json: %v", err)
	}

	return &artifact, nil
}

package services

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/domain"
	"github.com/matzapata/ipfs-compute/provider/pkg/archive"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
)

type ArtifactService struct {
	ArtifactRepository domain.IArtifactRepository
	Unzipper           archive.IUnzipper
}

func NewArtifactService(
	artifactRepository domain.IArtifactRepository,
	unzipper archive.IUnzipper,
) *ArtifactService {
	return &ArtifactService{
		ArtifactRepository: artifactRepository,
		Unzipper:           unzipper,
	}
}

func (d *ArtifactService) GetArtifactExecutable(cid string) (executablePath string, err error) {
	zippedExecutablePath, err := d.ArtifactRepository.GetZippedExecutable(cid, config.MAX_ZIPPED_DEPLOYMENT)
	if err != nil {
		return "", err
	}
	defer os.Remove(zippedExecutablePath)

	return d.Unzipper.Unzip(zippedExecutablePath)
}

func (d *ArtifactService) GetArtifactSpecification(cid string, providerRsaPrivateKey *crypto.RsaPrivateKey) (*domain.Artifact, error) {
	encSpecPath, err := d.ArtifactRepository.GetSpecificationFile(cid)
	if err != nil {
		return nil, err
	}

	// read the JSON encSpecData
	encSpecData, err := os.ReadFile(encSpecPath)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON data: %v", err)
	}
	artifactSpecification, err := crypto.RsaDecryptBytes(providerRsaPrivateKey, encSpecData)
	if err != nil {
		return nil, fmt.Errorf("error decrypting JSON data: %v", err)
	}

	// unmarshal the JSON data
	var artifact domain.Artifact
	err = json.Unmarshal(artifactSpecification, &artifact)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling metadata json: %v", err)
	}

	return &artifact, nil
}

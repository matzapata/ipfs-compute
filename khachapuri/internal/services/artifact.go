package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/domain"
	"github.com/matzapata/ipfs-compute/provider/pkg/archive"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
)

type ArtifactService struct {
	Config             *config.Config
	ArtifactRepository domain.IArtifactRepository
	Unzipper           archive.IUnzipper
}

func NewArtifactService(
	cfg *config.Config,
	artifactRepository domain.IArtifactRepository,
) *ArtifactService {
	return &ArtifactService{
		Config:             cfg,
		ArtifactRepository: artifactRepository,
		Unzipper:           archive.NewUnzipper(),
	}
}

func (d *ArtifactService) GetArtifact(cid string) (string, error) {
	zippedExecutablePath, err := d.ArtifactRepository.GetZippedExecutable(cid, d.Config.ArtifactMaxSize)
	if err != nil {
		return "", err
	}
	defer os.Remove(zippedExecutablePath)

	dir, err := os.MkdirTemp("", "khachapuri")
	if err != nil {
		return "", err
	}
	if err = d.Unzipper.Unzip(zippedExecutablePath, dir); err != nil {
		return "", err
	}

	return dir, nil
}

func (d *ArtifactService) GetArtifactSpecification(cid string, providerRsaPrivateKey *crypto.RsaPrivateKey) (*domain.ArtifactSpec, error) {
	specPath, err := d.ArtifactRepository.GetSpecificationFile(cid)
	if err != nil {
		return nil, err
	}

	// read the JSON encSpecData
	specData, err := os.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON data: %v", err)
	}

	// unmarshal the JSON data
	var artifact domain.ArtifactSpec
	if err = json.Unmarshal(specData, &artifact); err != nil {
		return nil, fmt.Errorf("error unmarshalling metadata json: %v", err)
	}

	// decrypt env vars
	encBytes, err := base64.StdEncoding.DecodeString(artifact.Env)
	if err != nil {
		return nil, err
	}
	decEnvVars, err := crypto.RsaDecryptBytes(providerRsaPrivateKey, encBytes)
	if err != nil {
		return nil, err
	}
	artifact.Env = string(decEnvVars)

	return &artifact, nil
}

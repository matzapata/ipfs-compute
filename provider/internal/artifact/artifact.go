package artifact

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"os"

	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/repositories"
	crypto_service "github.com/matzapata/ipfs-compute/provider/pkg/crypto"
	zip_service "github.com/matzapata/ipfs-compute/provider/pkg/zip"
)

type IArtifactService interface {
	GetArtifactExecutable(cid string) (executablePath string, err error)
	GetArtifactSpecification(cid string, providerRsaPrivateKey *rsa.PrivateKey) (*Artifact, error)
}

type ArtifactService struct {
	ArtifactRepository repositories.ArtifactRepository
	CryptoRsaService   crypto_service.ICryptoRsaService
	ZipService         zip_service.IZipService
}

type Artifact struct {
	Env            []string `json:"env"`
	Owner          string   `json:"owner"`
	OwnerSignature string   `json:"owner_signature"`
	DeploymentCid  string   `json:"deployment_cid"`
}

func NewArtifactService(
	artifactRepository repositories.ArtifactRepository,
	cryptoRsaService crypto_service.ICryptoRsaService,
	zipService zip_service.IZipService,
) *ArtifactService {
	return &ArtifactService{
		ArtifactRepository: artifactRepository,
		CryptoRsaService:   cryptoRsaService,
		ZipService:         zipService,
	}
}

func (d *ArtifactService) GetArtifactExecutable(cid string) (executablePath string, err error) {
	zippedExecutablePath, err := d.ArtifactRepository.GetZippedExecutable(cid, config.MAX_ZIPPED_DEPLOYMENT)
	if err != nil {
		return
	}
	defer os.Remove(zippedExecutablePath)

	return d.ZipService.Unzip(zippedExecutablePath)
}

func (d *ArtifactService) GetArtifactSpecification(cid string, providerRsaPrivateKey *rsa.PrivateKey) (*Artifact, error) {
	encSpecPath, err := d.ArtifactRepository.GetSpecificationFile(cid)
	if err != nil {
		return nil, err
	}

	// read the JSON encSpecData
	encSpecData, err := os.ReadFile(encSpecPath)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON data: %v", err)
	}
	artifactSpecification, err := d.CryptoRsaService.DecryptBytes(providerRsaPrivateKey, encSpecData)
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

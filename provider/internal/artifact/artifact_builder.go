package artifact

import (
	"crypto/rsa"
	"encoding/json"
	"os"

	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/repositories"
	"github.com/matzapata/ipfs-compute/provider/internal/source"

	crypto_service "github.com/matzapata/ipfs-compute/provider/pkg/crypto"
	zip_service "github.com/matzapata/ipfs-compute/provider/pkg/zip"
	cp "github.com/otiai10/copy"
)

type ArtifactBuilderService struct {
	SourceService    *source.SourceService
	ArtifactRepo     repositories.ArtifactRepository
	CryptoRsaService *crypto_service.CryptoRsaService
	ZipService       *zip_service.ZipService
}

func NewArtifactBuilderService(
	sourceService *source.SourceService,
	artifactRepository repositories.ArtifactRepository,
	cryptoRsaService *crypto_service.CryptoRsaService,
	zipService *zip_service.ZipService,
) *ArtifactBuilderService {
	return &ArtifactBuilderService{
		SourceService:    sourceService,
		ArtifactRepo:     artifactRepository,
		CryptoRsaService: cryptoRsaService,
		ZipService:       zipService,
	}
}

func (as *ArtifactBuilderService) BuildDeploymentZip() (cid string, err error) {
	// make dist folder
	err = os.Mkdir(config.DIST_DIR, 0755)
	if err != nil {
		return
	}
	defer os.RemoveAll(config.DIST_DIR)

	// get source files
	sources, err := as.SourceService.GetSource()
	if err != nil {
		return
	}

	// copy binary and public folder to dist folder
	err = os.Mkdir(config.DIST_DEPLOYMENT_DIR, 0755)
	if err != nil {
		return
	}
	err = cp.Copy(sources.ExecutablePath, config.DIST_BIN_FILE)
	if err != nil {
		return
	}
	err = cp.Copy(sources.AssetsPath, config.DIST_PUBLIC_DIR)
	if err != nil {
		return
	}

	// zip dist folder
	err = as.ZipService.ZipFolder(config.DIST_DEPLOYMENT_DIR, config.DIST_ZIP_FILE)
	if err != nil {
		return
	}

	// store zipped deployment
	cid, err = as.ArtifactRepo.CreateZippedExecutable(config.DIST_ZIP_FILE)
	if err != nil {
		return
	}

	return cid, nil
}

func (as *ArtifactBuilderService) BuildDeploymentSpecification(executableCid string, signature *crypto_service.Signature, providerPublicKey *rsa.PublicKey) (cid string, err error) {
	sourceSpec, err := as.SourceService.GetSourceSpecification()
	if err != nil {
		return
	}

	// encrypt it with public key
	deploymentJson, err := json.Marshal(Artifact{
		Env:            sourceSpec.Env,
		Owner:          signature.Address,
		OwnerSignature: signature.Signature,
		DeploymentCid:  executableCid,
	})
	if err != nil {
		return
	}
	encDeploymentJson, err := as.CryptoRsaService.EncryptBytes(providerPublicKey, deploymentJson)
	if err != nil {
		return
	}

	// save and store it
	err = os.WriteFile(config.DIST_SPEC_FILE, encDeploymentJson, 0644)
	if err != nil {
		return
	}
	defer os.Remove(config.DIST_SPEC_FILE)
	cid, err = as.ArtifactRepo.CreateSpecificationFile(config.DIST_SPEC_FILE)
	if err != nil {
		return
	}

	return cid, nil
}

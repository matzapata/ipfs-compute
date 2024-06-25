package artifact

import (
	"crypto/rsa"
	"encoding/json"
	"os"

	artifact_repository "github.com/matzapata/ipfs-compute/provider/internal/artifact/repository"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/source"
	ecdsa_helpers "github.com/matzapata/ipfs-compute/provider/pkg/helpers/ecdsa"
	files_helpers "github.com/matzapata/ipfs-compute/provider/pkg/helpers/files"
	rsa_helpers "github.com/matzapata/ipfs-compute/provider/pkg/helpers/rsa"
	cp "github.com/otiai10/copy"
)

type ArtifactBuilderService struct {
	SourceService *source.SourceService
	ArtifactRepo  artifact_repository.ArtifactRepository
}

func NewArtifactBuilderService(sourceService *source.SourceService, artifactRepository artifact_repository.ArtifactRepository) *ArtifactBuilderService {
	return &ArtifactBuilderService{
		SourceService: sourceService,
		ArtifactRepo:  artifactRepository,
	}
}

// TODO: these 2 should write to ipfs repo and read from source service
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
	err = files_helpers.ZipFolder(config.DIST_DEPLOYMENT_DIR, config.DIST_ZIP_FILE)
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

func (as *ArtifactBuilderService) BuildDeploymentSpecification(executableCid string, signature *ecdsa_helpers.Signature, providerPublicKey *rsa.PublicKey) (cid string, err error) {
	// get deployment spec
	source, err := as.SourceService.GetSource()
	if err != nil {
		return
	}

	deploymentSpecJson, err := os.ReadFile(source.SpecPath)
	if err != nil {
		return
	}
	var deploymentSpec ArtifactSpecification
	err = json.Unmarshal(deploymentSpecJson, &deploymentSpec)
	if err != nil {
		return
	}

	// encrypt it with public key
	deploymentJson, err := json.Marshal(Artifact{
		ArtifactSpecification: deploymentSpec,
		Owner:                 signature.Address,
		OwnerSignature:        signature.Signature,
		DeploymentCid:         executableCid,
	})
	if err != nil {
		return
	}
	encDeploymentJson, err := rsa_helpers.EncryptBytes(providerPublicKey, deploymentJson)
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

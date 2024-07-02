package services

import (
	"encoding/json"
	"os"
	"path"

	"github.com/matzapata/ipfs-compute/provider/internal/domain"

	"github.com/matzapata/ipfs-compute/provider/pkg/archive"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
	cp "github.com/otiai10/copy"
)

type ArtifactBuilderService struct {
	SourceService domain.ISourceService
	ArtifactRepo  domain.IArtifactRepository
	Zipper        archive.IZipper
}

func NewArtifactBuilderService(
	sourceService domain.ISourceService,
	artifactRepository domain.IArtifactRepository,
	zipper archive.IZipper,
) *ArtifactBuilderService {
	return &ArtifactBuilderService{
		SourceService: sourceService,
		ArtifactRepo:  artifactRepository,
		Zipper:        zipper,
	}
}

func (as *ArtifactBuilderService) BuildDeploymentZip() (cid string, err error) {
	// temp folder for the deployment
	tempDirPath, err := os.MkdirTemp("", "deployment-*")
	if err != nil {
		return
	}
	deploymentDirPath := path.Join(tempDirPath, "deployment")
	deploymentBinPath := path.Join(deploymentDirPath, "main")
	deploymentPublicPath := path.Join(deploymentDirPath, "public")
	deploymentZipPath := path.Join(tempDirPath, "deployment.zip")
	defer os.RemoveAll(tempDirPath)

	// get source files
	sources, err := as.SourceService.GetSource()
	if err != nil {
		return
	}

	// copy binary and public folder to dist folder
	err = os.Mkdir(deploymentDirPath, 0755)
	if err != nil {
		return
	}
	err = cp.Copy(sources.ExecutablePath, deploymentBinPath)
	if err != nil {
		return
	}
	err = cp.Copy(sources.AssetsPath, deploymentPublicPath)
	if err != nil {
		return
	}

	// zip dist folder
	err = as.Zipper.ZipFolder(deploymentDirPath, deploymentZipPath)
	if err != nil {
		return
	}

	// store zipped deployment
	cid, err = as.ArtifactRepo.CreateZippedExecutable(deploymentZipPath)
	if err != nil {
		return
	}

	return cid, nil
}

func (as *ArtifactBuilderService) BuildDeploymentSpecification(
	executableCid string,
	signature *crypto.EcdsaSignature,
	providerPublicKey *crypto.RsaPublicKey,
) (cid string, err error) {
	sourceSpec, err := as.SourceService.GetSourceSpecification()
	if err != nil {
		return
	}

	// encrypt it with public key
	deploymentJson, err := json.Marshal(domain.Artifact{
		Env:            sourceSpec.Env,
		Owner:          signature.Address,
		OwnerSignature: signature.Signature,
		DeploymentCid:  executableCid,
	})
	if err != nil {
		return
	}
	encDeploymentJson, err := crypto.RsaEncryptBytes(providerPublicKey, deploymentJson)
	if err != nil {
		return
	}

	// create temp file for the spec file
	distSpecFile, err := os.CreateTemp("", "spec-*.json")
	if err != nil {
		return
	}
	distSpecFilePath := distSpecFile.Name()

	// save and store it
	err = os.WriteFile(distSpecFilePath, encDeploymentJson, 0644)
	if err != nil {
		return
	}
	defer os.Remove(distSpecFilePath)
	cid, err = as.ArtifactRepo.CreateSpecificationFile(distSpecFilePath)
	if err != nil {
		return
	}

	return cid, nil
}

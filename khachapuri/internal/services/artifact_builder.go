package services

import (
	"encoding/base64"
	"os"
	"path"

	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/domain"

	"github.com/matzapata/ipfs-compute/provider/pkg/archive"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
	cp "github.com/otiai10/copy"
)

// run (runs the artifact build (loads artifact (local or remote)) (loads config (local or remote)))

type ArtifactBuilderService struct {
	Config          *config.Config
	SourceService   domain.ISourceService
	ArtifactRepo    domain.IArtifactRepository
	RegistryService domain.IRegistryService
	Zipper          archive.IZipper
}

func NewArtifactBuilderService(
	cfg *config.Config,
	sourceService domain.ISourceService,
	registryService domain.IRegistryService,
	artifactRepository domain.IArtifactRepository,
) *ArtifactBuilderService {
	return &ArtifactBuilderService{
		Config:          cfg,
		SourceService:   sourceService,
		ArtifactRepo:    artifactRepository,
		Zipper:          archive.NewZipper(),
		RegistryService: registryService,
	}
}

// generates artifact zip in outDir (./build by default)
func (as *ArtifactBuilderService) BuildArtifact(outDir string) error {
	if outDir == "" {
		outDir = "./build"
	}

	// check if folder exists
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err := os.Mkdir(outDir, 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		// remove all files in the folder
		err := os.RemoveAll(outDir)
		if err != nil {
			return err
		}
		err = os.Mkdir(outDir, 0755)
		if err != nil {
			return err
		}
	}

	// temp folder for the deployment
	deploymentDirPath := path.Join(outDir, "artifact")
	deploymentBinPath := path.Join(deploymentDirPath, "main")
	deploymentPublicPath := path.Join(deploymentDirPath, "public")
	deploymentZipPath := path.Join(outDir, "artifact.zip")
	defer os.RemoveAll(deploymentDirPath)

	// get source files
	sources, err := as.SourceService.GetSource()
	if err != nil {
		return err
	}

	// copy binary and public folder to dist folder
	err = os.Mkdir(deploymentDirPath, 0755)
	if err != nil {
		return err
	}
	err = cp.Copy(sources.ExecutablePath, deploymentBinPath)
	if err != nil {
		return err
	}
	err = cp.Copy(sources.AssetsPath, deploymentPublicPath)
	if err != nil {
		return err
	}

	// zip dist folder
	err = as.Zipper.ZipFolder(deploymentDirPath, deploymentZipPath)
	if err != nil {
		return err
	}

	return nil
}

func (as *ArtifactBuilderService) PublishArtifact(artPath string, adminSignature *crypto.EcdsaSignature, providerDomain string) (string, error) {
	sourceSpec, err := as.SourceService.GetSourceSpecification()
	if err != nil {
		return "", err
	}

	// resolve domain
	provider, err := as.RegistryService.ResolveDomain(providerDomain)
	if err != nil {
		return "", err
	}

	// encrypt env vars if any
	var envVars string
	for _, s := range sourceSpec.Env {
		envVars += s
	}
	encEnvVars, err := crypto.RsaEncryptBytes(
		crypto.RsaLoadPublicKeyFromString(provider.RsaPublicKey),
		[]byte(envVars),
	)
	if err != nil {
		return "", nil
	}

	// assemble artifact specification
	artSpec := domain.ArtifactSpec{
		Env:            base64.StdEncoding.EncodeToString(encEnvVars),
		Owner:          adminSignature.Address,
		OwnerSignature: adminSignature.Signature,
	}

	// publish artifact and update specification with cid
	artifactCid, err := as.ArtifactRepo.PublishArtifact(artPath)
	if err != nil {
		return "", err
	}
	artSpec.ArtifactCid = artifactCid

	// publish specification
	return as.ArtifactRepo.PublishArtifactSpecification(&artSpec)
}

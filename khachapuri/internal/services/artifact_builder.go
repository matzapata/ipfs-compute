package services

import (
	"encoding/base64"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/domain"

	"github.com/matzapata/ipfs-compute/provider/pkg/archive"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
	"github.com/matzapata/ipfs-compute/provider/pkg/system"
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

func (as *ArtifactBuilderService) BuildArtifact(serviceName string) error {
	// check if folder exists
	if err := system.EnsureDirExists(as.Config.ArtifactsPath, true); err != nil {
		return err
	}

	// temp folder for the deployment
	deploymentDirPath := path.Join(as.Config.ArtifactsPath, serviceName)
	deploymentZipPath := path.Join(as.Config.ArtifactsPath, serviceName+".zip")
	defer os.RemoveAll(deploymentDirPath)

	// get source files
	sources, err := as.SourceService.GetSource()
	if err != nil {
		return err
	}
	spec, err := as.SourceService.GetSourceSpecification()
	if err != nil {
		return err
	}

	// run build command
	if spec.Services[serviceName].Build.Command != "" {
		cmd := strings.Split(spec.Services[serviceName].Build.Command, " ")
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
			return err
		}
	}

	// copy assets recursively
	for _, asset := range spec.Services[serviceName].Assets {
		if err := cp.Copy(path.Join(sources.SourcePath, asset), path.Join(deploymentDirPath, asset)); err != nil {
			return err
		}
	}

	// zip dist folder
	err = as.Zipper.ZipFolder(deploymentDirPath, deploymentZipPath)
	if err != nil {
		return err
	}

	return nil
}

func (as *ArtifactBuilderService) PublishArtifact(serviceName string, providerDomain string, adminPk *crypto.EcdsaPrivateKey) (string, error) {
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
	for _, s := range sourceSpec.Services[serviceName].Env {
		envVars += s.Name + "=" + s.Value + "\n"
	}
	encEnvVars, err := crypto.RsaEncryptBytes(
		crypto.RsaLoadPublicKeyFromString(provider.RsaPublicKey),
		[]byte(envVars),
	)
	if err != nil {
		return "", nil
	}

	// publish artifact and update specification with cid
	artPath := path.Join(as.Config.ArtifactsPath, serviceName+".zip")
	if _, err := os.Stat(artPath); os.IsNotExist(err) {
		return "", err
	}
	artifactCid, err := as.ArtifactRepo.PublishArtifact(artPath)
	if err != nil {
		return "", err
	}

	// assemble artifact specification
	adminSignature, err := crypto.EcdsaSignMessage([]byte(artifactCid), adminPk)
	if err != nil {
		return "", err
	}
	artSpec := domain.ArtifactSpec{
		Env:            base64.StdEncoding.EncodeToString(encEnvVars),
		Owner:          adminSignature.Address,
		OwnerSignature: adminSignature.Signature,
		ArtifactCid:    artifactCid,
	}

	// publish specification
	return as.ArtifactRepo.PublishArtifactSpecification(&artSpec)
}

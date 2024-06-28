package domain

import "crypto/rsa"

type Artifact struct {
	Env            []string `json:"env"`
	Owner          string   `json:"owner"`
	OwnerSignature string   `json:"owner_signature"`
	DeploymentCid  string   `json:"deployment_cid"`
}

type IArtifactService interface {
	GetArtifactExecutable(cid string) (executablePath string, err error)
	GetArtifactSpecification(cid string, providerRsaPrivateKey *rsa.PrivateKey) (*Artifact, error)
}

type IArtifactRepository interface {
	GetZippedExecutable(cid string, maxSize uint) (zipPath string, err error)
	GetSpecificationFile(cid string) (specPath string, err error)
	CreateZippedExecutable(zipPath string) (cid string, err error)
	CreateSpecificationFile(specPath string) (cid string, err error)
}

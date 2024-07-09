package domain

import "crypto/rsa"

type ArtifactSpec struct {
	Env            string `json:"env"` // encrypted env vars
	Owner          string `json:"owner"`
	OwnerSignature string `json:"owner_signature"`
	ArtifactCid    string `json:"artifact_cid"`
}

type IArtifactService interface {
	GetArtifact(cid string) (artifactPath string, err error)
	GetArtifactSpecification(cid string, providerRsaPrivateKey *rsa.PrivateKey) (*ArtifactSpec, error)
}

type IArtifactRepository interface {
	GetZippedExecutable(cid string, maxSize uint) (zipPath string, err error)
	GetSpecificationFile(cid string) (specPath string, err error)
	PublishArtifact(artPath string) (cid string, err error)
	PublishArtifactSpecification(spec *ArtifactSpec) (cid string, err error)
}

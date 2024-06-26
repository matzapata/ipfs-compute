package repositories

type ArtifactRepository interface {
	GetZippedExecutable(cid string, maxSize uint) (zipPath string, err error)
	GetSpecificationFile(cid string) (specPath string, err error)
	CreateZippedExecutable(zipPath string) (cid string, err error)
	CreateSpecificationFile(specPath string) (cid string, err error)
}

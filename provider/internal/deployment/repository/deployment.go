package deployment_repository

type DeploymentsRepository interface {
	GetZippedDeployment(cid string, maxSize uint) ([]byte, error)
	GetDeploymentSpecFile(cid string) ([]byte, error)
}

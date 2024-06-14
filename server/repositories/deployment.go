package repositories

type DeploymentsRepository interface {
	GetZippedDeployment(cid string, maxSize uint) ([]byte, error)
}

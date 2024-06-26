package repositories

import "github.com/stretchr/testify/mock"

type MockArtifactRepository struct {
	mock.Mock
}

func NewArtifactRepository() *MockArtifactRepository {
	return &MockArtifactRepository{}
}

func (*MockArtifactRepository) GetZippedExecutable(cid string, maxSize uint) (zipPath string, err error) {
	panic("not implemented")
}

func (*MockArtifactRepository) GetSpecificationFile(cid string) (specPath string, err error) {
	panic("not implemented")
}

func (*MockArtifactRepository) CreateZippedExecutable(zipPath string) (cid string, err error) {
	panic("not implemented")
}

func (*MockArtifactRepository) CreateSpecificationFile(specPath string) (cid string, err error) {
	panic("not implemented")
}

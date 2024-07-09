package domain_mocks

import (
	"crypto/rsa"

	"github.com/matzapata/ipfs-compute/provider/internal/domain"
	"github.com/stretchr/testify/mock"
)

// mock service ========================================

type MockArtifactService struct {
	mock.Mock
}

func (m *MockArtifactService) GetArtifactExecutable(cid string) (executablePath string, err error) {
	args := m.Called(cid)
	return args.String(0), args.Error(1)
}

func (m *MockArtifactService) GetArtifactSpecification(cid string, providerRsaPrivateKey *rsa.PrivateKey) (*domain.ArtifactSpec, error) {
	args := m.Called(cid, providerRsaPrivateKey)
	return args.Get(0).(*domain.ArtifactSpec), args.Error(1)
}

// mock repository ========================================

type MockArtifactRepository struct {
	mock.Mock
}

func (m *MockArtifactRepository) GetZippedExecutable(cid string, maxSize uint) (zipPath string, err error) {
	args := m.Called(cid, maxSize)
	return args.String(0), args.Error(1)
}

func (m *MockArtifactRepository) GetSpecificationFile(cid string) (specPath string, err error) {
	args := m.Called(cid)
	return args.String(0), args.Error(1)
}

func (m *MockArtifactRepository) CreateZippedExecutable(zipPath string) (cid string, err error) {
	args := m.Called(zipPath)
	return args.String(0), args.Error(1)
}

func (m *MockArtifactRepository) CreateSpecificationFile(specPath string) (cid string, err error) {
	args := m.Called(specPath)
	return args.String(0), args.Error(1)
}

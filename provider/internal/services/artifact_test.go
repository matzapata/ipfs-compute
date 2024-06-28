package services_test

import (
	"errors"
	"testing"

	"github.com/matzapata/ipfs-compute/provider/internal/config"
	domain_mocks "github.com/matzapata/ipfs-compute/provider/internal/domain/mocks"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	crypto_service "github.com/matzapata/ipfs-compute/provider/pkg/crypto"
	zip_service "github.com/matzapata/ipfs-compute/provider/pkg/zip"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ArtifactTestSuite struct {
	suite.Suite
	mockRepo          *domain_mocks.MockArtifactRepository
	mockZipService    *zip_service.MockZipService
	mockCryptoService *crypto_service.MockCryptoRsaService
	service           *services.ArtifactService
}

func (suite *ArtifactTestSuite) SetupTest() {
	suite.mockRepo = new(domain_mocks.MockArtifactRepository)
	suite.mockZipService = new(zip_service.MockZipService)
	suite.mockCryptoService = new(crypto_service.MockCryptoRsaService)
	suite.service = services.NewArtifactService(suite.mockRepo, suite.mockCryptoService, suite.mockZipService)
}

// ==========================================================================================================

func (suite *ArtifactTestSuite) TestGetArtifactExecutable() {
	// service should return the unzipped executable path extracted from the zipped artifact repository

	cid := "test-cid"
	zippedPath := "zipped-executable-path"
	unzippedPath := "unzipped-executable-path"

	suite.mockRepo.On("GetZippedExecutable", cid, config.MAX_ZIPPED_DEPLOYMENT).Return(zippedPath, nil)
	suite.mockZipService.On("Unzip", zippedPath).Return(unzippedPath, nil)

	executablePath, err := suite.service.GetArtifactExecutable(cid)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), unzippedPath, executablePath)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockZipService.AssertExpectations(suite.T())
	assert.Equal(suite.T(), 1, 1)
}

func (suite *ArtifactTestSuite) TestGetArtifactExecutableUnzipError() {
	// service should return error if cannot unzip the artifact

	cid := "test-cid"
	zippedPath := "zipped-executable-path"
	expectedError := errors.New("error")

	suite.mockRepo.On("GetZippedExecutable", cid, config.MAX_ZIPPED_DEPLOYMENT).Return(zippedPath, nil)
	suite.mockZipService.On("Unzip", zippedPath).Return("", expectedError)

	executablePath, err := suite.service.GetArtifactExecutable(cid)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "", executablePath)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockZipService.AssertExpectations(suite.T())
}

func (suite *ArtifactTestSuite) TestGetArtifactExecutableRepoError() {
	// service should return error if cannot get the zipped artifact from the repository

	cid := "test-cid"
	expectedError := errors.New("error")

	suite.mockRepo.On("GetZippedExecutable", cid, config.MAX_ZIPPED_DEPLOYMENT).Return("", expectedError)

	executablePath, err := suite.service.GetArtifactExecutable(cid)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "", executablePath)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockZipService.AssertExpectations(suite.T())
}

// ==========================================================================================================

func TestArtifactTestSuite(t *testing.T) {
	suite.Run(t, new(ArtifactTestSuite))
}

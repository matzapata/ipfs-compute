package artifact

import (
	"testing"

	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/repositories"
	crypto_service "github.com/matzapata/ipfs-compute/provider/pkg/crypto"
	zip_service "github.com/matzapata/ipfs-compute/provider/pkg/zip"
	"github.com/stretchr/testify/assert"
)

func TestGetArtifactExecutable(t *testing.T) {
	mockRepo := new(repositories.MockArtifactRepository)
	mockZipService := new(zip_service.MockZipService)
	mockCryptoService := new(crypto_service.MockCryptoRsaService)

	service := NewArtifactService(mockRepo, mockCryptoService, mockZipService)

	cid := "test-cid"
	zippedPath := "zipped-executable-path"
	unzippedPath := "unzipped-executable-path"

	mockRepo.On("GetZippedExecutable", cid, config.MAX_ZIPPED_DEPLOYMENT).Return(zippedPath, nil)
	mockZipService.On("Unzip", zippedPath).Return(unzippedPath, nil)

	executablePath, err := service.GetArtifactExecutable(cid)

	assert.NoError(t, err)
	assert.Equal(t, unzippedPath, executablePath)
	mockRepo.AssertExpectations(t)
	mockZipService.AssertExpectations(t)
}

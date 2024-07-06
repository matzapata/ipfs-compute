package services_test

import (
	"errors"
	"testing"

	"github.com/matzapata/ipfs-compute/provider/internal/domain"
	domain_mocks "github.com/matzapata/ipfs-compute/provider/internal/domain/mocks"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestGetSource(t *testing.T) {
	mockSourceRepo := new(domain_mocks.MockSourceRepository)
	sourceService := &services.SourceService{
		SourceRepository: mockSourceRepo,
	}

	t.Run("success", func(t *testing.T) {
		expectedSource := &domain.Source{}
		mockSourceRepo.On("GetSource").Return(expectedSource, nil).Once()

		source, err := sourceService.GetSource()

		assert.NoError(t, err)
		assert.Equal(t, source, expectedSource)
		mockSourceRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		expectedError := errors.New("Not found")
		mockSourceRepo.On("GetSource").Return(nil, expectedError).Once()

		source, err := sourceService.GetSource()

		assert.Nil(t, source)
		assert.Equal(t, err, expectedError)
		mockSourceRepo.AssertExpectations(t)
	})
}

func TestGetSourceSpecification(t *testing.T) {
	mockSourceRepo := new(domain_mocks.MockSourceRepository)
	sourceService := &services.SourceService{
		SourceRepository: mockSourceRepo,
	}

	t.Run("success", func(t *testing.T) {
		expectedSource := &domain.SourceSpecification{}
		mockSourceRepo.On("GetSourceSpecification").Return(expectedSource, nil).Once()

		source, err := sourceService.GetSourceSpecification()

		assert.NoError(t, err)
		assert.Equal(t, source, expectedSource)
		mockSourceRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		expectedError := errors.New("Not found")
		mockSourceRepo.On("GetSourceSpecification").Return(nil, expectedError).Once()

		source, err := sourceService.GetSourceSpecification()

		assert.Nil(t, source)
		assert.Equal(t, err, expectedError)
		mockSourceRepo.AssertExpectations(t)
	})

}

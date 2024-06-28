package services_test

import (
	"errors"
	"testing"

	"github.com/matzapata/ipfs-compute/provider/internal/domain"
	domain_mocks "github.com/matzapata/ipfs-compute/provider/internal/domain/mocks"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	"github.com/stretchr/testify/suite"
)

type SourceSuiteTest struct {
	suite.Suite
	mockRepo *domain_mocks.MockSourceRepository
	service  *services.SourceService
}

func (s *SourceSuiteTest) SetupTest() {
	s.mockRepo = new(domain_mocks.MockSourceRepository)
	s.service = services.NewSourceService(s.mockRepo)
}

// ==========================================================================================================

func (s *SourceSuiteTest) TestGetSource() {
	// service should return the source from the repository

	expectedSource := &domain.Source{}
	s.mockRepo.On("GetSource").Return(expectedSource, nil)

	source, err := s.service.GetSource()

	s.mockRepo.AssertExpectations(s.T())
	s.Equal(expectedSource, source)
	s.NoError(err)
}

func (s *SourceSuiteTest) TestGetSourceError() {
	// service should return error if repository returns error

	expectedError := errors.New("error")
	s.mockRepo.On("GetSource").Return(nil, expectedError)

	source, err := s.service.GetSource()

	s.mockRepo.AssertExpectations(s.T())
	s.Nil(source)
	s.Error(err)
}

// ==========================================================================================================

func (s *SourceSuiteTest) TestGetSourceSpecification() {
	// service should return the source specification from the repository

	expectedSourceSpecification := &domain.SourceSpecification{}
	s.mockRepo.On("GetSourceSpecification").Return(expectedSourceSpecification, nil)

	sourceSpecification, err := s.service.GetSourceSpecification()

	s.mockRepo.AssertExpectations(s.T())
	s.Equal(expectedSourceSpecification, sourceSpecification)
	s.NoError(err)
}

func (s *SourceSuiteTest) TestGetSourceSpecificationError() {
	// service should return error if repository returns error

	expectedError := errors.New("error")
	s.mockRepo.On("GetSourceSpecification").Return(nil, expectedError)

	sourceSpecification, err := s.service.GetSourceSpecification()

	s.mockRepo.AssertExpectations(s.T())
	s.Nil(sourceSpecification)
	s.Error(err)
}

// ==========================================================================================================

func TestSourceSuiteTest(t *testing.T) {
	suite.Run(t, new(SourceSuiteTest))
}

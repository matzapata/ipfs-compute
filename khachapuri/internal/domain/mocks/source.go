package domain_mocks

import (
	"github.com/matzapata/ipfs-compute/provider/internal/domain"
	"github.com/stretchr/testify/mock"
)

// mock service ========================================

type MockSourceService struct {
	mock.Mock
}

func (m *MockSourceService) GetSource() (*domain.Source, error) {
	args := m.Called()
	return args.Get(0).(*domain.Source), args.Error(1)
}

func (m *MockSourceService) GetSourceSpecification() (*domain.SourceSpecification, error) {
	args := m.Called()
	return args.Get(0).(*domain.SourceSpecification), args.Error(1)
}

// mock repository ========================================

type MockSourceRepository struct {
	mock.Mock
}

func (m *MockSourceRepository) GetSource() (*domain.Source, error) {
	args := m.Called()

	var source *domain.Source
	if args.Get(0) != nil {
		source = args.Get(0).(*domain.Source)
	}
	return source, args.Error(1)
}

func (m *MockSourceRepository) GetSourceSpecification() (*domain.SourceSpecification, error) {
	args := m.Called()

	var sourceSpecification *domain.SourceSpecification = nil
	if args.Get(0) != nil {
		sourceSpecification = args.Get(0).(*domain.SourceSpecification)
	}
	return sourceSpecification, args.Error(1)
}

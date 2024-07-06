package services

import (
	"github.com/matzapata/ipfs-compute/provider/internal/domain"
)

type SourceService struct {
	SourceRepository domain.ISourceRepository
}

func NewSourceService(sourceRepository domain.ISourceRepository) *SourceService {
	return &SourceService{
		SourceRepository: sourceRepository,
	}
}

func (s *SourceService) GetSource() (*domain.Source, error) {
	return s.SourceRepository.GetSource()
}

func (s *SourceService) GetSourceSpecification() (*domain.SourceSpecification, error) {
	return s.SourceRepository.GetSourceSpecification()
}

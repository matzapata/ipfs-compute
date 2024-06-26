package source

import "github.com/matzapata/ipfs-compute/provider/internal/repositories"

type ISourceService interface {
	GetSource() (*repositories.Source, error)
	GetSourceSpecification() (*repositories.SourceSpecification, error)
}

type SourceService struct {
	Repo repositories.SourceRepository
}

func NewSourceService(repo repositories.SourceRepository) *SourceService {
	return &SourceService{
		Repo: repo,
	}
}

func (s *SourceService) GetSource() (*repositories.Source, error) {
	return s.Repo.GetSource()
}

func (s *SourceService) GetSourceSpecification() (*repositories.SourceSpecification, error) {
	return s.Repo.GetSourceSpecification()
}

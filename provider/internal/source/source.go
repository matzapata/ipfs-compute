package source

import (
	source_repository "github.com/matzapata/ipfs-compute/provider/internal/source/repository"
)

type SourceService struct {
	Repo source_repository.SourceRepository
}

func NewSourceService(repo source_repository.SourceRepository) *SourceService {
	return &SourceService{
		Repo: repo,
	}
}

func (s *SourceService) GetSource() (*source_repository.Source, error) {
	return s.Repo.GetSource()
}

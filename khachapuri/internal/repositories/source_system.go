package repositories

import (
	"encoding/json"
	"os"

	"github.com/matzapata/ipfs-compute/provider/internal/domain"
	files_helpers "github.com/matzapata/ipfs-compute/provider/pkg/helpers/files"
)

type SystemSourceRepository struct {
}

func NewSystemSourceRepository() *SystemSourceRepository {
	return &SystemSourceRepository{}
}

func (r *SystemSourceRepository) GetSource() (*domain.Source, error) {
	return &domain.Source{
		ExecutablePath: files_helpers.BuildCwdPath("main"),
		AssetsPath:     files_helpers.BuildCwdPath("public"),
		SpecPath:       files_helpers.BuildCwdPath("khachapuri.json"),
	}, nil
}

func (r *SystemSourceRepository) GetSourceSpecification() (*domain.SourceSpecification, error) {
	source, err := r.GetSource()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(source.SpecPath)
	if err != nil {
		return nil, err
	}

	var spec domain.SourceSpecification
	err = json.Unmarshal(data, &spec)
	if err != nil {
		return nil, err
	}

	return &spec, nil
}

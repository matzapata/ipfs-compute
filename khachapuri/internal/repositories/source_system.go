package repositories

import (
	"encoding/json"
	"os"

	"github.com/matzapata/ipfs-compute/provider/internal/domain"
	"github.com/matzapata/ipfs-compute/provider/pkg/system"
)

type SystemSourceRepository struct {
}

func NewSystemSourceRepository() *SystemSourceRepository {
	return &SystemSourceRepository{}
}

func (r *SystemSourceRepository) GetSource() (*domain.Source, error) {
	return &domain.Source{
		ExecutablePath: system.BuildCwdPath("main"),
		AssetsPath:     system.BuildCwdPath("public"),
		SpecPath:       system.BuildCwdPath("khachapuri.json"),
	}, nil
}

func (r *SystemSourceRepository) GetSourceSpecification() (*domain.SourceSpecification, error) {
	// TODO: load with viper config, load .env

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

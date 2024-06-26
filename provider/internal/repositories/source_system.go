package repositories

import (
	"encoding/json"
	"os"

	files_helpers "github.com/matzapata/ipfs-compute/provider/pkg/helpers/files"
)

type SystemSourceRepository struct {
}

func NewSystemSourceRepository() *SystemSourceRepository {
	return &SystemSourceRepository{}
}

func (r *SystemSourceRepository) GetSource() (*Source, error) {
	return &Source{
		ExecutablePath: files_helpers.BuildCwdPath("main"),
		AssetsPath:     files_helpers.BuildCwdPath("public"),
		SpecPath:       files_helpers.BuildCwdPath("khachapuri.json"),
	}, nil
}

func (r *SystemSourceRepository) GetSourceSpecification() (*SourceSpecification, error) {
	source, err := r.GetSource()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(source.SpecPath)
	if err != nil {
		return nil, err
	}

	var spec SourceSpecification
	err = json.Unmarshal(data, &spec)
	if err != nil {
		return nil, err
	}

	return &spec, nil
}

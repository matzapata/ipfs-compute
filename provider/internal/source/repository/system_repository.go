package source_repository

import files_helpers "github.com/matzapata/ipfs-compute/provider/pkg/helpers/files"

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

package repositories

import (
	"github.com/stretchr/testify/mock"
)

type MockSourceRepository struct {
	mock.Mock
}

func NewMockSourceRepository() *MockSourceRepository {
	return &MockSourceRepository{}
}

func (r *MockSourceRepository) GetSource() (*Source, error) {
	return &Source{
		ExecutablePath: "main",
		AssetsPath:     "public",
		SpecPath:       "khachapuri.json",
	}, nil
}

func (r *MockSourceRepository) GetSourceSpecification() (*SourceSpecification, error) {
	return &SourceSpecification{
		Env: []string{"ABC=124"},
	}, nil
}

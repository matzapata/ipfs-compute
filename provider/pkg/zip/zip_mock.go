package zip_service

import "github.com/stretchr/testify/mock"

type MockZipService struct {
	mock.Mock
}

func NewMockZipService() *MockZipService {
	return &MockZipService{}
}

func (z *MockZipService) Unzip(src string) (string, error) {
	panic("not implemented")
}

func (z *MockZipService) ZipFolder(srcFolder string, destZip string) error {
	panic("not implemented")
}

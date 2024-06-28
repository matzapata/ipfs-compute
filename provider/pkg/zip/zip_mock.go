package zip_service

import "github.com/stretchr/testify/mock"

type MockZipService struct {
	mock.Mock
}

func NewMockZipService() *MockZipService {
	return &MockZipService{}
}

func (m *MockZipService) Unzip(src string) (string, error) {
	args := m.Called(src)
	return args.String(0), args.Error(1)
}

func (m *MockZipService) ZipFolder(srcFolder string, destZip string) error {
	args := m.Called(srcFolder, destZip)
	return args.Error(1)
}

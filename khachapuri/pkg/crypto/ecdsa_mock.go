package crypto_service

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
)

type MockCryptoEcdsaService struct {
	mock.Mock
}

func (m *MockCryptoEcdsaService) SignMessage(data []byte, hexkey string) (*Signature, error) {
	args := m.Called(data, hexkey)
	return args.Get(0).(*Signature), args.Error(1)
}

func (m *MockCryptoEcdsaService) LoadPrivateKeyFromString(hexkey string) (*ecdsa.PrivateKey, error) {
	args := m.Called(hexkey)
	return args.Get(0).(*ecdsa.PrivateKey), args.Error(1)
}

func (m *MockCryptoEcdsaService) PrivateKeyToAddress(privateKey *ecdsa.PrivateKey) (common.Address, error) {
	args := m.Called(privateKey)
	return args.Get(0).(common.Address), args.Error(1)
}

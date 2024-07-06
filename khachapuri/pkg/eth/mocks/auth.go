package eth_mocks

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
	"github.com/stretchr/testify/mock"
)

type MockEthAuthenticator struct {
	mock.Mock
}

func (m *MockEthAuthenticator) Authenticate(privateKey *crypto.EcdsaPrivateKey) (*bind.TransactOpts, error) {
	args := m.Called(privateKey)
	return args.Get(0).(*bind.TransactOpts), args.Error(1)
}

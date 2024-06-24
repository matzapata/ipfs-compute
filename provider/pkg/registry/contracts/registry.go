// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package registry_contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// RegistryMetaData contains all meta data concerning the Registry contract.
var RegistryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"domain\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"resolver\",\"type\":\"address\"}],\"name\":\"NewResolver\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"domain\",\"type\":\"bytes32\"}],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"domain\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"resolverAddres\",\"type\":\"address\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"domain\",\"type\":\"bytes32\"}],\"name\":\"resolver\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"domain\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"resolverAddres\",\"type\":\"address\"}],\"name\":\"setResolver\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"domain\",\"type\":\"bytes32\"}],\"name\":\"unregister\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// RegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use RegistryMetaData.ABI instead.
var RegistryABI = RegistryMetaData.ABI

// Registry is an auto generated Go binding around an Ethereum contract.
type Registry struct {
	RegistryCaller     // Read-only binding to the contract
	RegistryTransactor // Write-only binding to the contract
	RegistryFilterer   // Log filterer for contract events
}

// RegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type RegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RegistrySession struct {
	Contract     *Registry         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RegistryCallerSession struct {
	Contract *RegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// RegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RegistryTransactorSession struct {
	Contract     *RegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// RegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type RegistryRaw struct {
	Contract *Registry // Generic contract binding to access the raw methods on
}

// RegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RegistryCallerRaw struct {
	Contract *RegistryCaller // Generic read-only contract binding to access the raw methods on
}

// RegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RegistryTransactorRaw struct {
	Contract *RegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRegistry creates a new instance of Registry, bound to a specific deployed contract.
func NewRegistry(address common.Address, backend bind.ContractBackend) (*Registry, error) {
	contract, err := bindRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Registry{RegistryCaller: RegistryCaller{contract: contract}, RegistryTransactor: RegistryTransactor{contract: contract}, RegistryFilterer: RegistryFilterer{contract: contract}}, nil
}

// NewRegistryCaller creates a new read-only instance of Registry, bound to a specific deployed contract.
func NewRegistryCaller(address common.Address, caller bind.ContractCaller) (*RegistryCaller, error) {
	contract, err := bindRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryCaller{contract: contract}, nil
}

// NewRegistryTransactor creates a new write-only instance of Registry, bound to a specific deployed contract.
func NewRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*RegistryTransactor, error) {
	contract, err := bindRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryTransactor{contract: contract}, nil
}

// NewRegistryFilterer creates a new log filterer instance of Registry, bound to a specific deployed contract.
func NewRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*RegistryFilterer, error) {
	contract, err := bindRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RegistryFilterer{contract: contract}, nil
}

// bindRegistry binds a generic wrapper to an already deployed contract.
func bindRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registry *RegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Registry.Contract.RegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registry *RegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registry.Contract.RegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registry *RegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registry.Contract.RegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registry *RegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Registry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registry *RegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registry *RegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registry.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x02571be3.
//
// Solidity: function owner(bytes32 domain) view returns(address)
func (_Registry *RegistryCaller) Owner(opts *bind.CallOpts, domain [32]byte) (common.Address, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "owner", domain)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x02571be3.
//
// Solidity: function owner(bytes32 domain) view returns(address)
func (_Registry *RegistrySession) Owner(domain [32]byte) (common.Address, error) {
	return _Registry.Contract.Owner(&_Registry.CallOpts, domain)
}

// Owner is a free data retrieval call binding the contract method 0x02571be3.
//
// Solidity: function owner(bytes32 domain) view returns(address)
func (_Registry *RegistryCallerSession) Owner(domain [32]byte) (common.Address, error) {
	return _Registry.Contract.Owner(&_Registry.CallOpts, domain)
}

// Resolver is a free data retrieval call binding the contract method 0x0178b8bf.
//
// Solidity: function resolver(bytes32 domain) view returns(address)
func (_Registry *RegistryCaller) Resolver(opts *bind.CallOpts, domain [32]byte) (common.Address, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "resolver", domain)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Resolver is a free data retrieval call binding the contract method 0x0178b8bf.
//
// Solidity: function resolver(bytes32 domain) view returns(address)
func (_Registry *RegistrySession) Resolver(domain [32]byte) (common.Address, error) {
	return _Registry.Contract.Resolver(&_Registry.CallOpts, domain)
}

// Resolver is a free data retrieval call binding the contract method 0x0178b8bf.
//
// Solidity: function resolver(bytes32 domain) view returns(address)
func (_Registry *RegistryCallerSession) Resolver(domain [32]byte) (common.Address, error) {
	return _Registry.Contract.Resolver(&_Registry.CallOpts, domain)
}

// Register is a paid mutator transaction binding the contract method 0xd22057a9.
//
// Solidity: function register(bytes32 domain, address resolverAddres) returns()
func (_Registry *RegistryTransactor) Register(opts *bind.TransactOpts, domain [32]byte, resolverAddres common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "register", domain, resolverAddres)
}

// Register is a paid mutator transaction binding the contract method 0xd22057a9.
//
// Solidity: function register(bytes32 domain, address resolverAddres) returns()
func (_Registry *RegistrySession) Register(domain [32]byte, resolverAddres common.Address) (*types.Transaction, error) {
	return _Registry.Contract.Register(&_Registry.TransactOpts, domain, resolverAddres)
}

// Register is a paid mutator transaction binding the contract method 0xd22057a9.
//
// Solidity: function register(bytes32 domain, address resolverAddres) returns()
func (_Registry *RegistryTransactorSession) Register(domain [32]byte, resolverAddres common.Address) (*types.Transaction, error) {
	return _Registry.Contract.Register(&_Registry.TransactOpts, domain, resolverAddres)
}

// SetResolver is a paid mutator transaction binding the contract method 0x1896f70a.
//
// Solidity: function setResolver(bytes32 domain, address resolverAddres) returns()
func (_Registry *RegistryTransactor) SetResolver(opts *bind.TransactOpts, domain [32]byte, resolverAddres common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "setResolver", domain, resolverAddres)
}

// SetResolver is a paid mutator transaction binding the contract method 0x1896f70a.
//
// Solidity: function setResolver(bytes32 domain, address resolverAddres) returns()
func (_Registry *RegistrySession) SetResolver(domain [32]byte, resolverAddres common.Address) (*types.Transaction, error) {
	return _Registry.Contract.SetResolver(&_Registry.TransactOpts, domain, resolverAddres)
}

// SetResolver is a paid mutator transaction binding the contract method 0x1896f70a.
//
// Solidity: function setResolver(bytes32 domain, address resolverAddres) returns()
func (_Registry *RegistryTransactorSession) SetResolver(domain [32]byte, resolverAddres common.Address) (*types.Transaction, error) {
	return _Registry.Contract.SetResolver(&_Registry.TransactOpts, domain, resolverAddres)
}

// Unregister is a paid mutator transaction binding the contract method 0x1a0919dc.
//
// Solidity: function unregister(bytes32 domain) returns()
func (_Registry *RegistryTransactor) Unregister(opts *bind.TransactOpts, domain [32]byte) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "unregister", domain)
}

// Unregister is a paid mutator transaction binding the contract method 0x1a0919dc.
//
// Solidity: function unregister(bytes32 domain) returns()
func (_Registry *RegistrySession) Unregister(domain [32]byte) (*types.Transaction, error) {
	return _Registry.Contract.Unregister(&_Registry.TransactOpts, domain)
}

// Unregister is a paid mutator transaction binding the contract method 0x1a0919dc.
//
// Solidity: function unregister(bytes32 domain) returns()
func (_Registry *RegistryTransactorSession) Unregister(domain [32]byte) (*types.Transaction, error) {
	return _Registry.Contract.Unregister(&_Registry.TransactOpts, domain)
}

// RegistryNewResolverIterator is returned from FilterNewResolver and is used to iterate over the raw logs and unpacked data for NewResolver events raised by the Registry contract.
type RegistryNewResolverIterator struct {
	Event *RegistryNewResolver // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RegistryNewResolverIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryNewResolver)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RegistryNewResolver)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RegistryNewResolverIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryNewResolverIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryNewResolver represents a NewResolver event raised by the Registry contract.
type RegistryNewResolver struct {
	Domain   [32]byte
	Resolver common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNewResolver is a free log retrieval operation binding the contract event 0x335721b01866dc23fbee8b6b2c7b1e14d6f05c28cd35a2c934239f94095602a0.
//
// Solidity: event NewResolver(bytes32 indexed domain, address resolver)
func (_Registry *RegistryFilterer) FilterNewResolver(opts *bind.FilterOpts, domain [][32]byte) (*RegistryNewResolverIterator, error) {

	var domainRule []interface{}
	for _, domainItem := range domain {
		domainRule = append(domainRule, domainItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "NewResolver", domainRule)
	if err != nil {
		return nil, err
	}
	return &RegistryNewResolverIterator{contract: _Registry.contract, event: "NewResolver", logs: logs, sub: sub}, nil
}

// WatchNewResolver is a free log subscription operation binding the contract event 0x335721b01866dc23fbee8b6b2c7b1e14d6f05c28cd35a2c934239f94095602a0.
//
// Solidity: event NewResolver(bytes32 indexed domain, address resolver)
func (_Registry *RegistryFilterer) WatchNewResolver(opts *bind.WatchOpts, sink chan<- *RegistryNewResolver, domain [][32]byte) (event.Subscription, error) {

	var domainRule []interface{}
	for _, domainItem := range domain {
		domainRule = append(domainRule, domainItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "NewResolver", domainRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryNewResolver)
				if err := _Registry.contract.UnpackLog(event, "NewResolver", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewResolver is a log parse operation binding the contract event 0x335721b01866dc23fbee8b6b2c7b1e14d6f05c28cd35a2c934239f94095602a0.
//
// Solidity: event NewResolver(bytes32 indexed domain, address resolver)
func (_Registry *RegistryFilterer) ParseNewResolver(log types.Log) (*RegistryNewResolver, error) {
	event := new(RegistryNewResolver)
	if err := _Registry.contract.UnpackLog(event, "NewResolver", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

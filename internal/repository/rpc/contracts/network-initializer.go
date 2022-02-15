// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// NetworkInitializerABI is the input ABI used to generate the binding from.
const NetworkInitializerABI = "[{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"sealedEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSupply\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_sfc\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_auth\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_driver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_evmWriter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"initializeAll\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// NetworkInitializer is an auto generated Go binding around an Ethereum contract.
type NetworkInitializer struct {
	NetworkInitializerCaller     // Read-only binding to the contract
	NetworkInitializerTransactor // Write-only binding to the contract
	NetworkInitializerFilterer   // Log filterer for contract events
}

// NetworkInitializerCaller is an auto generated read-only Go binding around an Ethereum contract.
type NetworkInitializerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkInitializerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NetworkInitializerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkInitializerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NetworkInitializerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkInitializerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NetworkInitializerSession struct {
	Contract     *NetworkInitializer // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// NetworkInitializerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NetworkInitializerCallerSession struct {
	Contract *NetworkInitializerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// NetworkInitializerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NetworkInitializerTransactorSession struct {
	Contract     *NetworkInitializerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// NetworkInitializerRaw is an auto generated low-level Go binding around an Ethereum contract.
type NetworkInitializerRaw struct {
	Contract *NetworkInitializer // Generic contract binding to access the raw methods on
}

// NetworkInitializerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NetworkInitializerCallerRaw struct {
	Contract *NetworkInitializerCaller // Generic read-only contract binding to access the raw methods on
}

// NetworkInitializerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NetworkInitializerTransactorRaw struct {
	Contract *NetworkInitializerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNetworkInitializer creates a new instance of NetworkInitializer, bound to a specific deployed contract.
func NewNetworkInitializer(address common.Address, backend bind.ContractBackend) (*NetworkInitializer, error) {
	contract, err := bindNetworkInitializer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NetworkInitializer{NetworkInitializerCaller: NetworkInitializerCaller{contract: contract}, NetworkInitializerTransactor: NetworkInitializerTransactor{contract: contract}, NetworkInitializerFilterer: NetworkInitializerFilterer{contract: contract}}, nil
}

// NewNetworkInitializerCaller creates a new read-only instance of NetworkInitializer, bound to a specific deployed contract.
func NewNetworkInitializerCaller(address common.Address, caller bind.ContractCaller) (*NetworkInitializerCaller, error) {
	contract, err := bindNetworkInitializer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkInitializerCaller{contract: contract}, nil
}

// NewNetworkInitializerTransactor creates a new write-only instance of NetworkInitializer, bound to a specific deployed contract.
func NewNetworkInitializerTransactor(address common.Address, transactor bind.ContractTransactor) (*NetworkInitializerTransactor, error) {
	contract, err := bindNetworkInitializer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkInitializerTransactor{contract: contract}, nil
}

// NewNetworkInitializerFilterer creates a new log filterer instance of NetworkInitializer, bound to a specific deployed contract.
func NewNetworkInitializerFilterer(address common.Address, filterer bind.ContractFilterer) (*NetworkInitializerFilterer, error) {
	contract, err := bindNetworkInitializer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NetworkInitializerFilterer{contract: contract}, nil
}

// bindNetworkInitializer binds a generic wrapper to an already deployed contract.
func bindNetworkInitializer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NetworkInitializerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NetworkInitializer *NetworkInitializerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NetworkInitializer.Contract.NetworkInitializerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NetworkInitializer *NetworkInitializerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkInitializer.Contract.NetworkInitializerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NetworkInitializer *NetworkInitializerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NetworkInitializer.Contract.NetworkInitializerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NetworkInitializer *NetworkInitializerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NetworkInitializer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NetworkInitializer *NetworkInitializerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkInitializer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NetworkInitializer *NetworkInitializerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NetworkInitializer.Contract.contract.Transact(opts, method, params...)
}

// InitializeAll is a paid mutator transaction binding the contract method 0xc80e1513.
//
// Solidity: function initializeAll(uint256 sealedEpoch, uint256 totalSupply, address _sfc, address _auth, address _driver, address _evmWriter, address _owner) returns()
func (_NetworkInitializer *NetworkInitializerTransactor) InitializeAll(opts *bind.TransactOpts, sealedEpoch *big.Int, totalSupply *big.Int, _sfc common.Address, _auth common.Address, _driver common.Address, _evmWriter common.Address, _owner common.Address) (*types.Transaction, error) {
	return _NetworkInitializer.contract.Transact(opts, "initializeAll", sealedEpoch, totalSupply, _sfc, _auth, _driver, _evmWriter, _owner)
}

// InitializeAll is a paid mutator transaction binding the contract method 0xc80e1513.
//
// Solidity: function initializeAll(uint256 sealedEpoch, uint256 totalSupply, address _sfc, address _auth, address _driver, address _evmWriter, address _owner) returns()
func (_NetworkInitializer *NetworkInitializerSession) InitializeAll(sealedEpoch *big.Int, totalSupply *big.Int, _sfc common.Address, _auth common.Address, _driver common.Address, _evmWriter common.Address, _owner common.Address) (*types.Transaction, error) {
	return _NetworkInitializer.Contract.InitializeAll(&_NetworkInitializer.TransactOpts, sealedEpoch, totalSupply, _sfc, _auth, _driver, _evmWriter, _owner)
}

// InitializeAll is a paid mutator transaction binding the contract method 0xc80e1513.
//
// Solidity: function initializeAll(uint256 sealedEpoch, uint256 totalSupply, address _sfc, address _auth, address _driver, address _evmWriter, address _owner) returns()
func (_NetworkInitializer *NetworkInitializerTransactorSession) InitializeAll(sealedEpoch *big.Int, totalSupply *big.Int, _sfc common.Address, _auth common.Address, _driver common.Address, _evmWriter common.Address, _owner common.Address) (*types.Transaction, error) {
	return _NetworkInitializer.Contract.InitializeAll(&_NetworkInitializer.TransactOpts, sealedEpoch, totalSupply, _sfc, _auth, _driver, _evmWriter, _owner)
}

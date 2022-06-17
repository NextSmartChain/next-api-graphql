/*
Package rpc implements bridge to NEXT Smart Chain full node API interface.

We recommend using local IPC for fast and the most efficient inter-process communication between the API server
and an NEXT Smart Chain node. Any remote RPC connection will work, but the performance may be significantly degraded
by extra networking overhead of remote RPC calls.

You should also consider security implications of opening NEXT Smart Chain RPC interface for remote access.
If you considering it as your deployment strategy, you should establish encrypted channel between the API server
and NEXT Smart Chain RPC interface with connection limited to specified endpoints.

We strongly discourage opening NEXT Smart Chain RPC interface for unrestricted Internet access.
*/
package rpc

import (
	"context"
	"next-api-graphql/internal/config"
	"next-api-graphql/internal/logger"
	"next-api-graphql/internal/repository/rpc/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	etc "github.com/ethereum/go-ethereum/core/types"
	eth "github.com/ethereum/go-ethereum/ethclient"
	next "github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/sync/singleflight"
	"strings"
	"sync"
)

// rpcHeadProxyChannelCapacity represents the capacity of the new received blocks proxy channel.
const rpcHeadProxyChannelCapacity = 10000

// NextBridge represents NEXT Smart Chain RPC abstraction layer.
type NextBridge struct {
	rpc *next.Client
	eth *eth.Client
	log logger.Logger
	cg  *singleflight.Group

	// fMintCfg represents the configuration of the fMint protocol
	sigConfig     *config.ServerSignature
	sfcConfig     *config.Staking
	uniswapConfig *config.DeFiUniswap

	// extended minter config
	fMintCfg fMintConfig
	fLendCfg fLendConfig

	// common contracts
	sfcAbi      *abi.ABI
	sfcContract *contracts.SfcContract

	// received blocks proxy
	wg       *sync.WaitGroup
	sigClose chan bool
	headers  chan *etc.Header
}

// New creates new NEXT Smart Chain RPC connection bridge.
func New(cfg *config.Config, log logger.Logger) (*NextBridge, error) {
	cli, con, err := connect(cfg, log)
	if err != nil {
		log.Criticalf("can not open connection; %s", err.Error())
		return nil, err
	}

	// build the bridge structure using the con we have
	br := &NextBridge{
		rpc: cli,
		eth: con,
		log: log,
		cg:  new(singleflight.Group),

		// special configuration options below this line
		sigConfig:     &cfg.MySignature,
		sfcConfig:     &cfg.Staking,
		uniswapConfig: &cfg.DeFi.Uniswap,
		fMintCfg: fMintConfig{
			addressProvider: cfg.DeFi.FMint.AddressProvider,
		},
		fLendCfg: fLendConfig{lendigPoolAddress: cfg.DeFi.FLend.LendingPool},

		// configure block observation loop
		wg:       new(sync.WaitGroup),
		sigClose: make(chan bool, 1),
		headers:  make(chan *etc.Header, rpcHeadProxyChannelCapacity),
	}

	// inform about the local address of the API node
	log.Noticef("using signature address %s", br.sigConfig.Address.String())

	// add the bridge ref to the fMintCfg and return the instance
	br.fMintCfg.bridge = br
	br.run()
	return br, nil
}

// connect opens connections we need to communicate with the blockchain node.
func connect(cfg *config.Config, log logger.Logger) (*next.Client, *eth.Client, error) {
	// log what we do
	log.Infof("connecting blockchain node at %s", cfg.Next.Url)

	// try to establish a connection
	client, err := next.Dial(cfg.Next.Url)
	if err != nil {
		log.Critical(err)
		return nil, nil, err
	}

	// try to establish a for smart contract interaction
	con, err := eth.Dial(cfg.Next.Url)
	if err != nil {
		log.Critical(err)
		return nil, nil, err
	}

	// log
	log.Notice("node connection open")
	return client, con, nil
}

// run starts the bridge threads required to collect blockchain data.
func (next *NextBridge) run() {
	next.wg.Add(1)
	go next.observeBlocks()
}

// terminate kills the bridge threads to end the bridge gracefully.
func (next *NextBridge) terminate() {
	next.sigClose <- true
	next.wg.Wait()
	next.log.Noticef("rpc threads terminated")
}

// Close will finish all pending operations and terminate the NEXT Smart Chain RPC connection
func (next *NextBridge) Close() {
	// terminate threads before we close connections
	next.terminate()

	// do we have a connection?
	if next.rpc != nil {
		next.rpc.Close()
		next.eth.Close()
		next.log.Info("blockchain connections are closed")
	}
}

// Connection returns open NEXT Smart Chain connection.
func (next *NextBridge) Connection() *next.Client {
	return next.rpc
}

// DefaultCallOpts creates a default record for call options.
func (next *NextBridge) DefaultCallOpts() *bind.CallOpts {
	// get the default call opts only once if called in parallel
	co, _, _ := next.cg.Do("default-call-opts", func() (interface{}, error) {
		return &bind.CallOpts{
			Pending:     false,
			From:        next.sigConfig.Address,
			BlockNumber: nil,
			Context:     context.Background(),
		}, nil
	})
	return co.(*bind.CallOpts)
}

// SfcContract returns instance of SFC contract for interaction.
func (next *NextBridge) SfcContract() *contracts.SfcContract {
	// lazy create SFC contract instance
	if nil == next.sfcContract {
		// instantiate the contract and display its name
		var err error
		next.sfcContract, err = contracts.NewSfcContract(next.sfcConfig.SFCContract, next.eth)
		if err != nil {
			next.log.Criticalf("failed to instantiate SFC contract; %s", err.Error())
			panic(err)
		}
	}
	return next.sfcContract
}

// SfcAbi returns a parse ABI of the AFC contract.
func (next *NextBridge) SfcAbi() *abi.ABI {
	if nil == next.sfcAbi {
		ab, err := abi.JSON(strings.NewReader(contracts.SfcContractABI))
		if err != nil {
			next.log.Criticalf("failed to parse SFC contract ABI; %s", err.Error())
			panic(err)
		}
		next.sfcAbi = &ab
	}
	return next.sfcAbi
}

// ObservedBlockProxy provides a channel fed with new headers observed
// by the connected blockchain node.
func (next *NextBridge) ObservedBlockProxy() chan *etc.Header {
	return next.headers
}

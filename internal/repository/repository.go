/*
Package repository implements repository for handling fast and efficient access to data required
by the resolvers of the API server.

Internally it utilizes RPC to access Opera/Lachesis full node for blockchain interaction. Mongo database
for fast, robust and scalable off-chain data storage, especially for aggregated and pre-calculated data mining
results. BigCache for in-memory object storage to speed up loading of frequently accessed entities.
*/
package repository

import (
	"fmt"
	"sync"
	"time"

	"fantom-api-graphql/internal/config"
	"fantom-api-graphql/internal/logger"
	"fantom-api-graphql/internal/repository/cache"
	"fantom-api-graphql/internal/repository/db"
	"fantom-api-graphql/internal/repository/rpc"
	"fantom-api-graphql/internal/repository/rpc/contracts"
	"fantom-api-graphql/internal/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/sync/singleflight"
)

// repo represents an instance of the Repository manager.
var repo Repository

// onceRepo is the sync object used to make sure the Repository
// is instantiated only once on the first demand.
var onceRepo sync.Once

// config represents the configuration setup used by the repository
// to establish and maintain required connectivity to external services
// as needed.
var cfg *config.Config

// log represents the logger to be used by the repository.
var log logger.Logger

// SetConfig sets the repository configuration to be used to establish
// and maintain external repository connections.
func SetConfig(c *config.Config) {
	cfg = c
}

// SetLogger sets the repository logger to be used to collect logging info.
func SetLogger(l logger.Logger) {
	log = l
}

// R provides access to the singleton instance of the Repository.
func R() Repository {
	// make sure to instantiate the Repository only once
	onceRepo.Do(func() {
		repo = newRepository()
	})
	return repo
}

// Proxy represents Repository interface implementation and controls access to data
// trough several low level bridges.
type proxy struct {
	cache *cache.MemBridge
	db    *db.MongoDbBridge
	rpc   *rpc.FtmBridge
	log   logger.Logger
	cfg   *config.Config

	// transaction estimator counter
	txCount uint64

	// we need a Group to use single flight to control price pulls
	apiRequestGroup singleflight.Group

	// governance contracts reference
	govContracts map[string]*config.GovernanceContract

	// smart contract compilers
	solCompiler string
}

// newRepository creates new instance of Repository implementation, namely proxy structure.
func newRepository() Repository {
	if cfg == nil {
		panic(fmt.Errorf("missing configuration"))
	}
	if log == nil {
		panic(fmt.Errorf("missing logger"))
	}

	// create connections
	caBridge, dbBridge, rpcBridge, err := connect(cfg, log)
	if err != nil {
		log.Fatal("repository init failed")
		return nil
	}

	// construct the proxy instance
	p := proxy{
		cache: caBridge,
		db:    dbBridge,
		rpc:   rpcBridge,
		log:   log,
		cfg:   cfg,

		// get the map of governance contracts
		govContracts: governanceContractsMap(&cfg.Governance),

		// keep reference to the SOL compiler
		solCompiler: cfg.Compiler.DefaultSolCompilerPath,
	}

	registerSystemContracts(&p)

	// return the proxy
	return &p
}

// governanceContractsMap creates map of governance contracts keyed
// by the contract address.
func governanceContractsMap(cfg *config.Governance) map[string]*config.GovernanceContract {
	// prep the result set
	res := make(map[string]*config.GovernanceContract)

	// collect all the configured governance contracts into the map
	for _, gv := range cfg.Contracts {
		res[gv.Address.String()] = &gv
	}
	return res
}

// connect opens connections to the external sources we need.
func connect(cfg *config.Config, log logger.Logger) (*cache.MemBridge, *db.MongoDbBridge, *rpc.FtmBridge, error) {
	// create new in-memory cache bridge
	caBridge, err := cache.New(cfg, log)
	if err != nil {
		log.Criticalf("can not create in-memory cache bridge, %s", err.Error())
		return nil, nil, nil, err
	}

	// create new database connection bridge
	dbBridge, err := db.New(cfg, log)
	if err != nil {
		log.Criticalf("can not connect backend persistent storage, %s", err.Error())
		return nil, nil, nil, err
	}

	// create new Lachesis RPC bridge
	rpcBridge, err := rpc.New(cfg, log)
	if err != nil {
		log.Criticalf("can not connect Lachesis RPC interface, %s", err.Error())
		return nil, nil, nil, err
	}
	return caBridge, dbBridge, rpcBridge, nil
}

// register system contracts
func registerSystemContracts(p *proxy) {

	blockNumber := hexutil.Uint64(0)
	block, err := p.BlockByNumber(&blockNumber)
	if err != nil {
		log.Criticalf("unable to retrieve block 0, %s", err.Error())
		return
	}

	// create pseudo transaction
	tx := new(types.Transaction)
	tx.TimeStamp = time.Unix(int64(block.TimeStamp), 0)
	tx.Hash = block.Hash //0 hash

	account := types.Account{
		Address:      p.cfg.Staking.NodeDriverContract,
		ContractTx:   &block.Hash, //0 hash
		Type:         types.AccountTypeSFC,
		LastActivity: block.TimeStamp,
		TrxCounter:   1,
	}

	version, err := p.SfcVersion()
	if err != nil {
		log.Criticalf("unable to retrieve sfc version, %s", err.Error())
		return
	}

	if !p.db.IsContractKnown(&p.cfg.Staking.NodeDriverContract) {
		log.Info("register node driver contract")
		if err = p.StoreAccount(&account); err == nil {
			p.StoreContract(types.NewSfcContract(&p.cfg.Staking.NodeDriverContract, uint64(version), "NodeDriver", contracts.NodeDriverABI, block, tx))
		}
	}

	if !p.db.IsContractKnown(&p.cfg.Staking.SFCContract) {
		log.Info("register sfc contract")
		account.Address = p.cfg.Staking.SFCContract
		if err = p.StoreAccount(&account); err == nil {
			p.StoreContract(types.NewSfcContract(&p.cfg.Staking.SFCContract, uint64(version), "SFC Contract", contracts.SfcContractABI, block, tx))
		}
	}
	if !p.db.IsContractKnown(&p.cfg.Staking.NetworkInitializerContract) {
		log.Info("register network initializer contract")
		account.Address = p.cfg.Staking.NetworkInitializerContract
		if err = p.StoreAccount(&account); err == nil {
			p.StoreContract(types.NewSfcContract(&p.cfg.Staking.NetworkInitializerContract, uint64(version), "Network Initializer", contracts.NetworkInitializerABI, block, tx))
		}
	}

}

// Close with close all connections and clean up the pending work for graceful termination.
func (p *proxy) Close() {
	// inform about actions
	p.log.Notice("repository is closing")

	// close connections
	p.db.Close()
	p.rpc.Close()

	// inform about actions
	p.log.Notice("repository done")
}

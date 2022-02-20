/*
Package repository implements repository for handling fast and efficient access to data required
by the resolvers of the API server.

Internally it utilizes RPC to access Opera/Lachesis full node for blockchain interaction. Mongo database
for fast, robust and scalable off-chain data storage, especially for aggregated and pre-calculated data mining
results. BigCache for in-memory object storage to speed up loading of frequently accessed entities.
*/
package repository

import (
	"bytes"
	"fantom-api-graphql/internal/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// PullValidatorInfo extracts an extended validator information from smart contact.
func (p *proxy) PullValidatorInfo(id *hexutil.Big) (*types.ValidatorInfo, error) {
	// retieve from rpc
	info, err := p.rpc.ValidatorInfo(id)
	if err != nil {
		return nil, err
	}
	if info == nil {
		info = new(types.ValidatorInfo)
		p.StoreValidatorInfo(id, info)
	}
	return info, nil
}

// StoreValidatorInfo stores validator information to in-memory cache for future use.
func (p *proxy) StoreValidatorInfo(id *hexutil.Big, sti *types.ValidatorInfo) error {
	// push to in-memory cache
	err := p.cache.PushValidatorInfo(id, sti)
	if err != nil {
		p.log.Error("validator info can net be kept")
		return err
	}
	return nil
}

// RetrieveValidatorInfo gets validator information from in-memory if available.
func (p *proxy) RetrieveValidatorInfo(id *hexutil.Big) *types.ValidatorInfo {
	info := p.cache.PullValidatorInfo(id)
	if info == nil {
		if info, err := p.PullValidatorInfo(id); err != nil || info.Name == nil {
			return nil
		}
	}
	return info
}

// IsStiContract returns true if the given address points to the STI contract.
func (p *proxy) IsStiContract(addr *common.Address) bool {
	return bytes.Equal(addr.Bytes(), p.cfg.Staking.StiContract.Bytes())
}

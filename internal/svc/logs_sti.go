// Package svc implements blockchain data processing services.
package svc

import (
	"fantom-api-graphql/internal/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

// handleDelegationLog handles a new delegation event from logs.
func handleStiInfoUpdated(lr *types.LogRecord) {
	// check for the correct recipient
	if lr.Trx.To == nil || *lr.Trx.To != cfg.Staking.StiContract {
		return
	}
	// get Staker Id
	stakerID := (*hexutil.Big)(new(big.Int).SetBytes(lr.Data))

	// is this a tx made on sti contract?
	if info, err := repo.PullStakerInfo(stakerID); err == nil && info != nil {
		err = repo.StoreStakerInfo(stakerID, info)
	}
}

// Package svc implements blockchain data processing services.
package svc

import (
	"fantom-api-graphql/internal/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

// handleValidatorInfoUpdated handles validator info updates event from logs.
func handleValidatorInfoUpdated(lr *types.LogRecord) {
	// check for address
	if lr.Address != cfg.Staking.SFCContract {
		return
	}
	// get Staker Id
	validatorID := (*hexutil.Big)(new(big.Int).SetBytes(lr.Data))

	// is this a tx made on sti contract?
	if info, err := repo.PullValidatorInfo(validatorID); err == nil && info != nil {
		err = repo.StoreValidatorInfo(validatorID, info)
	}
}

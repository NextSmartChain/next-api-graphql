// Package svc implements blockchain data processing services.
package svc

import (
	"next-api-graphql/internal/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

// handleValidatorInfoUpdated handles validator info updates event from logs.
func handleValidatorInfoUpdated(lr *types.LogRecord) {
	// check for address
	if lr.Address != cfg.Staking.SFCContract {
		return
	}
	// get validator Id
	validatorID := (*hexutil.Big)(new(big.Int).SetBytes(lr.Data))

	// update data from rpc
	repo.UpdateValidatorInfo(validatorID)
}

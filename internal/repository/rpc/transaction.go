/*
Package rpc implements bridge to Lachesis full node API interface.

We recommend using local IPC for fast and the most efficient inter-process communication between the API server
and an NEXT Smart Chain node. Any remote RPC connection will work, but the performance may be significantly degraded
by extra networking overhead of remote RPC calls.

You should also consider security implications of opening Lachesis RPC interface for a remote access.
If you considering it as your deployment strategy, you should establish encrypted channel between the API server
and Lachesis RPC interface with connection limited to specified endpoints.

We strongly discourage opening Lachesis RPC interface for unrestricted Internet access.
*/
package rpc

import (
	"next-api-graphql/internal/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	retypes "github.com/ethereum/go-ethereum/core/types"
)

// Transaction returns information about a blockchain transaction by hash.
func (next *NextBridge) Transaction(hash *common.Hash) (*types.Transaction, error) {
	// keep track of the operation
	next.log.Debugf("loading transaction %s", hash.String())

	// call for data
	var trx types.Transaction
	err := next.rpc.Call(&trx, "next_getTransactionByHash", hash)
	if err != nil {
		next.log.Error("transaction could not be extracted")
		return nil, err
	}

	// is there a block reference already?
	if trx.BlockNumber != nil {
		// get transaction receipt
		var rec struct {
			Index             hexutil.Uint64  `json:"transactionIndex"`
			CumulativeGasUsed hexutil.Uint64  `json:"cumulativeGasUsed"`
			GasUsed           hexutil.Uint64  `json:"gasUsed"`
			ContractAddress   *common.Address `json:"contractAddress,omitempty"`
			Status            hexutil.Uint64  `json:"status"`
			Logs              []retypes.Log   `json:"logs"`
		}

		// call for the transaction receipt data
		err := next.rpc.Call(&rec, "next_getTransactionReceipt", hash)
		if err != nil {
			next.log.Errorf("can not get receipt for transaction %s", hash)
			return nil, err
		}

		// copy some data
		trx.Index = &rec.Index
		trx.CumulativeGasUsed = &rec.CumulativeGasUsed
		trx.GasUsed = &rec.GasUsed
		trx.ContractAddress = rec.ContractAddress
		trx.Status = &rec.Status
		trx.Logs = rec.Logs
	}

	// keep track of the operation
	next.log.Debugf("transaction %s loaded", hash.String())
	return &trx, nil
}

// SendTransaction sends raw signed and RLP encoded transaction to the block chain.
func (next *NextBridge) SendTransaction(tx hexutil.Bytes) (*common.Hash, error) {
	// keep track of the operation
	next.log.Debug("sending new transaction to block chain")

	var hash common.Hash
	err := next.rpc.Call(&hash, "eth_sendRawTransaction", tx)
	if err != nil {
		next.log.Error("transaction could not be sent")
		return nil, err
	}

	// keep track of the operation
	next.log.Debugf("transaction has been accepted with hash %s", hash.String())
	return &hash, nil
}

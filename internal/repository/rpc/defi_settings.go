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
	"next-api-graphql/internal/repository/rpc/contracts"
	"next-api-graphql/internal/types"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

//go:generate tools/abigen.sh --abi ./contracts/abi/defi-fmint-address-provider.abi --pkg contracts --type DefiFMintAddressProvider --out ./contracts/fmint_addresses.go

// tConfigItemsLoaders defines a map between DeFi config elements and their respective loaders.
type tConfigItemsLoaders map[*hexutil.Big]func(*bind.CallOpts) (*big.Int, error)

// DefiConfiguration resolves the current DeFi contract settings.
func (next *NextBridge) DefiConfiguration() (*types.DefiSettings, error) {
	// access the contract
	contract, err := next.fMintCfg.fMintMinterContract()
	if err != nil {
		return nil, err
	}

	// create the container
	ds := types.DefiSettings{
		FMintContract:           next.fMintCfg.mustContractAddress(fMintAddressMinter),
		FMintAddressProvider:    next.fMintCfg.addressProvider,
		FMintTokenRegistry:      next.fMintCfg.mustContractAddress(fMintAddressTokenRegistry),
		FMintRewardDistribution: next.fMintCfg.mustContractAddress(fMintAddressRewardDistribution),
		FMintCollateralPool:     next.fMintCfg.mustContractAddress(fMintCollateralPool),
		FMintDebtPool:           next.fMintCfg.mustContractAddress(fMintDebtPool),
		PriceOracleAggregate:    next.fMintCfg.mustContractAddress(fMintAddressPriceOracleProxy),
	}

	// prep to load certain values
	loaders := tConfigItemsLoaders{
		&ds.MintFee4:               contract.GetFMintFee4dec,
		&ds.MinCollateralRatio4:    contract.GetCollateralLowestDebtRatio4dec,
		&ds.RewardCollateralRatio4: contract.GetRewardEligibilityRatio4dec,
	}

	// load all the configured values
	if err := next.pullSetOfDefiConfigValues(loaders); err != nil {
		next.log.Errorf("can not pull defi config values; %s", err.Error())
		return nil, err
	}

	// load the decimals correction
	if ds.Decimals, err = next.pullDefiDecimalCorrection(contract); err != nil {
		next.log.Errorf("can not pull defi decimals correction; %s", err.Error())
		return nil, err
	}

	// return the config
	return &ds, nil
}

// pullSetOfDefiConfigValues pulls set of DeFi configuration values for the given
// config loaders map.
func (next *NextBridge) pullDefiDecimalCorrection(con *contracts.DefiFMintMinter) (int32, error) {
	// load the decimals correction
	val, err := next.pullDefiConfigValue(con.FMintFeeDigitsCorrection)
	if err != nil {
		next.log.Errorf("can not pull decimals correction; %s", err.Error())
		return 0, err
	}

	// calculate number of decimals
	var dec int32
	var value = val.ToInt().Uint64()
	for value > 1 {
		value /= 10
		dec++
	}

	// convert and return
	return dec, nil
}

// pullSetOfDefiConfigValues pulls set of DeFi configuration values for the given
// config loaders map.
func (next *NextBridge) pullSetOfDefiConfigValues(loaders tConfigItemsLoaders) error {
	// collect loaders error
	var err error

	// loop the map and load the values
	for ref, fn := range loaders {
		*ref, err = next.pullDefiConfigValue(fn)
		if err != nil {
			return err
		}
	}

	return nil
}

// tradeFee4 pulls DeFi trading fee from the Liquidity Pool contract.
func (next *NextBridge) pullDefiConfigValue(cf func(*bind.CallOpts) (*big.Int, error)) (hexutil.Big, error) {
	// pull the trading fee value
	val, err := cf(nil)
	if err != nil {
		return hexutil.Big{}, err
	}

	// do we have the value? we should always have
	if val == nil {
		return hexutil.Big{}, fmt.Errorf("defi config value not available")
	}

	return hexutil.Big(*val), nil
}

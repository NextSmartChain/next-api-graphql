// Package cache implements bridge to fast in-memory object cache.
package cache

import (
	"next-api-graphql/internal/types"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"strings"
)

// sfcMaxDelegatedRatioKey represents the key used to store SFC delegation ratio.
const (
	sfcMaxDelegatedRatioKey = "sfc_dlr"
	sfcConfigurationKey     = "sfc_cfg"
	sfcValidatorAddress     = "val_adr"
	sfcTotalStakedKey       = "staked_total"
	sfcValidatorInfoPrefix  = "validator_info_"
)

// PullSfcMaxDelegatedRatio extract the ratio from cache, if possible.
func (b *MemBridge) PullSfcMaxDelegatedRatio() *big.Int {
	// try to get the account data from the cache
	data, err := b.cache.Get(sfcMaxDelegatedRatioKey)
	if err != nil {
		return nil
	}
	return new(big.Int).SetBytes(data)
}

// PushSfcMaxDelegatedRatio stores the ratio in cache, if possible.
func (b *MemBridge) PushSfcMaxDelegatedRatio(val *big.Int) {
	if val == nil {
		return
	}
	if err := b.cache.Set(sfcMaxDelegatedRatioKey, val.Bytes()); err != nil {
		b.log.Errorf("can not store SFC delegation ratio value")
	}
}

// PullSfcConfig extract the SFC configuration from cache, if possible.
func (b *MemBridge) PullSfcConfig() *types.SfcConfig {
	// try to get the account data from the cache
	data, err := b.cache.Get(sfcConfigurationKey)
	if err != nil {
		return nil
	}

	// make the data container
	val := types.SfcConfig{
		MinValidatorStake:      hexutil.Big{},
		MaxDelegatedRatio:      hexutil.Big{},
		MinLockupDuration:      hexutil.Big{},
		MaxLockupDuration:      hexutil.Big{},
		WithdrawalPeriodEpochs: hexutil.Big{},
		WithdrawalPeriodTime:   hexutil.Big{},
	}

	// decode data
	if err := val.Unmarshal(data); err != nil {
		b.log.Errorf("can not decode SFC config; %s", err.Error())
		return nil
	}
	return &val
}

// PushSfcConfig stores the SFC configuration, if possible.
func (b *MemBridge) PushSfcConfig(val *types.SfcConfig) {
	if val == nil {
		return
	}

	// get the encoded config
	data, err := val.Marshal()
	if err != nil {
		b.log.Errorf("can not encode SFC config; %s", err.Error())
		return
	}

	// store the data
	if err := b.cache.Set(sfcConfigurationKey, data); err != nil {
		b.log.Errorf("can not store SFC configuration")
	}
}

// validatorAddressKey generates cache key for address of the given validator id.
func validatorAddressKey(valID *hexutil.Big) string {
	var sb strings.Builder
	sb.WriteString(sfcValidatorAddress)
	sb.WriteString(valID.String())
	return sb.String()
}

// PushValidatorAddress stores validator address in the memory cache.
func (b *MemBridge) PushValidatorAddress(valID *hexutil.Big, adr *common.Address) {
	// empty validator ID or address? nothing to do
	if nil == valID || nil == adr {
		return
	}

	// store the address
	if err := b.cache.Set(validatorAddressKey(valID), adr.Bytes()); err != nil {
		b.log.Errorf("can not store address of validator %d", valID.ToInt().Uint64())
	}
}

// PullValidatorAddress tries to pull the validator address from memory cache.
func (b *MemBridge) PullValidatorAddress(valID *hexutil.Big) *common.Address {
	// empty validator ID?
	if nil == valID {
		return nil
	}

	// try to get the account data from the cache
	data, err := b.cache.Get(validatorAddressKey(valID))
	if err != nil {
		return nil
	}
	adr := common.BytesToAddress(data)
	return &adr
}

// sfciCacheKeyPrefix is the prefix used for cache key to store validator information.
const sfciCacheKeyPrefix = "validator_info_"

// sfciTotalStakedKey is the cache key used to store total staked amount
const sfciTotalStakedKey = "staked_total"

// PullValidatorInfo extracts validator information from the in-memory cache if available.
func (b *MemBridge) PullValidatorInfo(id *hexutil.Big) *types.ValidatorInfo {
	// try to get the account data from the cache
	data, err := b.cache.Get(getValidatorInfoKey(id))
	if err != nil {
		// cache returns ErrEntryNotFound if the key does not exist
		return nil
	}

	// do we have the data?
	info, err := types.UnmarshalValidatorInfo(data)
	if err != nil {
		b.log.Criticalf("can not decode validator information data from in-memory cache; %s", err.Error())
		return nil
	}

	return info
}

// PushValidatorInfo stores provided validator information in the in-memory cache.
func (b *MemBridge) PushValidatorInfo(id *hexutil.Big, sfci *types.ValidatorInfo) error {
	// encode account
	data, err := sfci.Marshal()
	if err != nil {
		b.log.Criticalf("can not marshal validator information to JSON; %s", err.Error())
		return err
	}

	// set the data to cache by block number
	return b.cache.Set(getValidatorInfoKey(id), data)
}

// getPriceKeyBySymbol build a cache key for the given price symbol.
func getValidatorInfoKey(id *hexutil.Big) string {
	// use the builder to make the string
	var sb strings.Builder

	sb.WriteString(sfcValidatorInfoPrefix)
	sb.WriteString(id.String())

	return sb.String()
}

// PullTotalStaked extracts total staked amount from the in-memory cache if available.
func (b *MemBridge) PullTotalStaked() *hexutil.Big {
	// try to get the account data from the cache
	data, err := b.cache.Get(sfcTotalStakedKey)
	if err != nil {
		return nil
	}

	// do we have the data?
	val := new(big.Int).SetBytes(data)
	return (*hexutil.Big)(val)
}

// PushTotalStaked stores provided total staked amount information in the in-memory cache.
func (b *MemBridge) PushTotalStaked(amount *hexutil.Big) error {
	// we must have the value
	if amount == nil {
		b.log.Criticalf("can not store invalid amount")
		return fmt.Errorf("amount not provided")
	}

	// encode account
	return b.cache.Set(sfcTotalStakedKey, amount.ToInt().Bytes())
}

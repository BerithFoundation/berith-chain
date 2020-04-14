// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"math/big"

	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/consensus"
	"github.com/BerithFoundation/berith-chain/core/types"
	"github.com/BerithFoundation/berith-chain/core/vm"
)

// ChainContext supports retrieving headers and consensus parameters from the
// current blockchain to be used during transaction processing.
type ChainContext interface {
	// Engine retrieves the chain's consensus engine.
	Engine() consensus.Engine

	// GetHeader returns the hash corresponding to their hash.
	GetHeader(common.Hash, uint64) *types.Header
}

// NewEVMContext creates a new context for use in the EVM.
func NewEVMContext(msg Message, header *types.Header, chain ChainContext, author *common.Address, isBIP4 bool) vm.Context {
	// If we don't have an explicit author (i.e. not mining), extract from the header
	var beneficiary common.Address
	if author == nil {
		beneficiary, _ = chain.Engine().Author(header) // Ignore error, we're past header validation
	} else {
		beneficiary = *author
	}

	return vm.Context{
		CanTransfer: getCanTransferFunc(isBIP4, msg.Target()),
		Transfer:    Transfer,
		GetHash:     GetHashFn(header, chain),
		Origin:      msg.From(),
		Coinbase:    beneficiary,
		BlockNumber: new(big.Int).Set(header.Number),
		Time:        new(big.Int).Set(header.Time),
		Difficulty:  new(big.Int).Set(header.Difficulty),
		GasLimit:    header.GasLimit,
		GasPrice:    new(big.Int).Set(msg.GasPrice()),
	}
}

/*
	[BERITH]
	하드포크 BIP4와 관련하여 StakeBalance가 한도를 초과할 수 없도록 설정
*/
func getCanTransferFunc(isBIP4 bool, target types.JobWallet) vm.CanTransferFunc {
	if isBIP4 && target == types.Stake {
		return CanTransferBIP4
	}
	return CanTransfer
}

// GetHashFn returns a GetHashFunc which retrieves header hashes by number
func GetHashFn(ref *types.Header, chain ChainContext) func(n uint64) common.Hash {
	var cache map[uint64]common.Hash

	return func(n uint64) common.Hash {
		// If there's no hash cache yet, make one
		if cache == nil {
			cache = map[uint64]common.Hash{
				ref.Number.Uint64() - 1: ref.ParentHash,
			}
		}
		// Try to fulfill the request from the cache
		if hash, ok := cache[n]; ok {
			return hash
		}
		// Not cached, iterate the blocks and cache the hashes
		for header := chain.GetHeader(ref.ParentHash, ref.Number.Uint64()-1); header != nil; header = chain.GetHeader(header.ParentHash, header.Number.Uint64()-1) {
			cache[header.Number.Uint64()-1] = header.ParentHash
			if n == header.Number.Uint64()-1 {
				return header.ParentHash
			}
		}
		return common.Hash{}
	}
}

// CanTransfer checks whether there are enough funds in the address' account to make a transfer.
// This does not take the necessary gas in to account to make the transfer valid.
func CanTransfer(db vm.StateDB, addr common.Address, amount *big.Int, base types.JobWallet, limitStakeBalance *big.Int) bool {
	if base == types.Main {
		return db.GetBalance(addr).Cmp(amount) >= 0
	} else if base == types.Stake {
		return db.GetStakeBalance(addr).Cmp(amount) >= 0
	}

	return false
}

/*
	[BERITH]
	하드포크 BIP4가 적용되는 블록이고, target이 Stake Balance인 경우 CanTransferBIP4 함수로 전송 가능 여부를 체크
 */
func CanTransferBIP4(db vm.StateDB, addr common.Address, amount *big.Int, base types.JobWallet, limitStakeBalance *big.Int) bool {
	stakeBalance := db.GetStakeBalance(addr)
	totalStakingBalance := new(big.Int).Add(stakeBalance, amount)

	if base == types.Main {
		return db.GetBalance(addr).Cmp(amount) >= 0 && CheckStakeBalanceAmount(totalStakingBalance, limitStakeBalance)
	} else if base == types.Stake {
		return db.GetStakeBalance(addr).Cmp(amount) >= 0
	}
	return false
}

/*
	[BERITH]
	Stake Balance 한도 초과 체크 함수
*/
func CheckStakeBalanceAmount(totalStakingAmount, maximum *big.Int) bool {
	return totalStakingAmount.Cmp(maximum) != 1
}

// Transfer subtracts amount from sender and adds amount to recipient using the given Db
func Transfer(db vm.StateDB, sender, recipient common.Address, amount, blockNumber *big.Int, base, target types.JobWallet) {
	/*
		[BERITH]
		Tx 를 state에 적용
	*/
	switch base {
	case types.Main:
		if target == types.Main {
			db.SubBalance(sender, amount)
			db.AddBalance(recipient, amount)
		} else if target == types.Stake {
			//베이스 지갑 차감
			db.SubBalance(sender, amount)
			db.AddStakeBalance(recipient, amount, blockNumber)
		}

		break
	case types.Stake:
		if target == types.Main {
			//스테이크 풀시
			db.RemoveStakeBalance(sender)
		}
		break
	}
}

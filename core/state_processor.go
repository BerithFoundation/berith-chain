// Copyright 2015 The go-ethereum Authors
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
	"berith-chain/berith/staking"
	"math/big"

	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/consensus"
	"github.com/BerithFoundation/berith-chain/consensus/misc"
	"github.com/BerithFoundation/berith-chain/core/state"
	"github.com/BerithFoundation/berith-chain/core/types"
	"github.com/BerithFoundation/berith-chain/core/vm"
	"github.com/BerithFoundation/berith-chain/crypto"
	"github.com/BerithFoundation/berith-chain/params"
)

// StateProcessor is a basic Processor, which takes care of transitioning
// state from one point to another.
//
// StateProcessor implements Processor.
type StateProcessor struct {
	config *params.ChainConfig // Chain configuration options
	bc     *BlockChain         // Canonical block chain
	engine consensus.Engine    // Consensus engine used for block rewards
}

// NewStateProcessor initialises a new StateProcessor.
func NewStateProcessor(config *params.ChainConfig, bc *BlockChain, engine consensus.Engine) *StateProcessor {
	return &StateProcessor{
		config: config,
		bc:     bc,
		engine: engine,
	}
}

// Process processes the state changes according to the Berith rules by running
// the transaction messages using the statedb and applying any rewards to both
// the processor (coinbase) and any included uncles.
//
// Process returns the receipts and logs accumulated during the process and
// returns the amount of gas that was used in the process. If any of the
// transactions failed to execute due to insufficient gas it will return an error.
func (p *StateProcessor) Process(block *types.Block, statedb *state.StateDB, cfg vm.Config) (types.Receipts, []*types.Log, uint64, error) {
	var (
		receipts types.Receipts
		usedGas  = new(uint64)
		header   = block.Header()
		allLogs  []*types.Log
		gp       = new(GasPool).AddGas(block.GasLimit())
	)
	// Mutate the block and state according to any hard-fork specs
	if p.config.DAOForkSupport && p.config.DAOForkBlock != nil && p.config.DAOForkBlock.Cmp(block.Number()) == 0 {
		misc.ApplyDAOHardFork(statedb)
	}
	// Iterate over and process the individual transactions
	for i, tx := range block.Transactions() {
		statedb.Prepare(tx.Hash(), block.Hash(), i)
		receipt, _, err := ApplyTransaction(p.config, p.bc, nil, gp, statedb, header, tx, usedGas, cfg)
		if err != nil {
			return nil, nil, 0, err
		}
		receipts = append(receipts, receipt)
		allLogs = append(allLogs, receipt.Logs...)
	}
	// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)
	_, err := p.engine.Finalize(p.bc, header, statedb, block.Transactions(), block.Uncles(), receipts)

	return receipts, allLogs, *usedGas, err
}

// ApplyTransaction attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. It returns the receipt
// for the transaction, gas used and an error if the transaction failed,
// indicating the block was invalid.
func ApplyTransaction(config *params.ChainConfig, bc ChainContext, author *common.Address, gp *GasPool, statedb *state.StateDB, header *types.Header, tx *types.Transaction, usedGas *uint64, cfg vm.Config) (*types.Receipt, uint64, error) {
	msg, err := tx.AsMessage(types.MakeSigner(config, header.Number))
	if err != nil {
		return nil, 0, err
	}

	adjustStateForBIP4(config, statedb, header, tx)

	//[BERITH]
	// Create a new context to be used in the EVM environment
	context := NewEVMContext(msg, header, bc, author)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	vmenv := vm.NewEVM(context, statedb, config, cfg)
	// Apply the transaction to the current state (included in the env)
	_, gas, failed, err := ApplyMessage(vmenv, msg, gp)
	if err != nil {
		return nil, 0, err
	}
	// Update the state with pending changes
	var root []byte
	if config.IsByzantium(header.Number) {
		statedb.Finalise(true)
	} else {
		root = statedb.IntermediateRoot(config.IsEIP158(header.Number)).Bytes()
	}
	*usedGas += gas

	// Create a new receipt for the transaction, storing the intermediate root and gas used by the tx
	// based on the eip phase, we're passing whether the root touch-delete accounts.
	receipt := types.NewReceipt(root, failed, *usedGas)
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = gas
	// if the transaction created a contract, store the creation address in the receipt.
	if msg.To() == nil {
		receipt.ContractAddress = crypto.CreateAddress(vmenv.Context.Origin, tx.Nonce())
	}
	// Set the receipt logs and create a bloom for filtering
	receipt.Logs = statedb.GetLogs(tx.Hash())
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})

	return receipt, gas, err
}

/*
[Berith]
adjust Stake balance and Selection point For hard fork BIP4
Check the Recipient's Stake Balance of the transaction to be processed, and change it if it has more than the limit.
*/
func adjustStateForBIP4(config *params.ChainConfig, statedb *state.StateDB, header *types.Header, tx *types.Transaction) {
	stakedBalance := big.NewInt(0)
	var recipient *common.Address
	if tx.To() != nil {
		recipient = tx.To()
		stakedBalance = statedb.GetStakeBalance(*recipient)
	}

	if config.IsBIP4(header.Number) && stakedBalance.Cmp(config.Bsrr.LimitStakeBalance) == 1 {
		// Adjust staking balance of accounts staking above the limit
		difference := new(big.Int).Sub(stakedBalance, config.Bsrr.LimitStakeBalance)
		statedb.AddStakeBalance(*recipient, new(big.Int).Neg(difference), header.Number)
		statedb.AddBalance(*recipient, difference)

		// Adjust selection selectionPoint of accounts staking above the limit
		currentBlock := header.Number
		lastStkBlock := new(big.Int).Set(statedb.GetStakeUpdated(*recipient))
		selectionPoint := staking.CalcPointBigint(config.Bsrr.LimitStakeBalance, big.NewInt(0), currentBlock, lastStkBlock, config.Bsrr.Period)
		statedb.SetPoint(*recipient, selectionPoint)
	}
}

/*
[Berith]
Check if the break transaction satisfies the lock up condition
The Break Transaction has a three-day grace period.
*/
func checkBreakTransaction(msg types.Message, lastStakedBlock, currentBlock *big.Int, period uint64) (bool, int64) {
	lockUpPeriod := big.NewInt(int64((60 * 60 * 24 * 3) / int64(period))) // Created blocks in 3 days
	elapsedBlockNumber := new(big.Int).Sub(currentBlock, lastStakedBlock)
	return elapsedBlockNumber.Cmp(lockUpPeriod) == 1, new(big.Int).Sub(lockUpPeriod, elapsedBlockNumber).Int64()
}

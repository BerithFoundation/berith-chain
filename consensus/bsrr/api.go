/*
d8888b. d88888b d8888b. d888888b d888888b db   db
88  `8D 88'     88  `8D   `88'   `~~88~~' 88   88
88oooY' 88ooooo 88oobY'    88       88    88ooo88
88~~~b. 88~~~~~ 88`8b      88       88    88~~~88
88   8D 88.     88 `88.   .88.      88    88   88
Y8888P' Y88888P 88   YD Y888888P    YP    YP   YP

	  copyrights by ibizsoftware 2018 - 2019
*/
/*
[BERITH]
Consensus algorithm related API
*/

package bsrr

import (
	"berith-chain/berith/selection"
	"berith-chain/common"
	"berith-chain/consensus"
	"berith-chain/core/types"
	"berith-chain/rpc"
)

// API is a user facing RPC API to allow controlling the signer and voting
// mechanisms of the proof-of-authority scheme.
type API struct {
	chain consensus.ChainReader
	bsrr  *BSRR
}

func (api *API) GetCandidates(number *rpc.BlockNumber) (*selection.JSONCandidates, error) {
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}

	if header == nil {
		return nil, errUnknownBlock
	}

	parent := api.chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return nil, consensus.ErrUnknownAncestor
	}

	target, exist := api.bsrr.getStakeTargetBlock(api.chain, parent)
	if !exist {
		return nil, consensus.ErrUnknownAncestor
	}

	stat, err := api.chain.StateAt(target.Root)

	if err != nil {
		return nil, err
	}

	stks, err := api.bsrr.getStakers(api.chain, target.Number.Uint64(), target.Hash())
	if err != nil {
		return nil, err
	}

	return selection.GetCandidates(parent.Number.Uint64(), parent.Hash(), stks, stat), nil

}

/*
[BERITH]
Function that returns the selected Block Creator on the current local block
*/
func (api *API) GetBlockCreators(number *rpc.BlockNumber) ([]common.Address, error) {
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}

	if header == nil {
		return nil, errUnknownBlock
	}

	parent := api.chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return nil, consensus.ErrUnknownAncestor
	}

	target, exist := api.bsrr.getStakeTargetBlock(api.chain, parent)
	if !exist {
		return nil, consensus.ErrUnknownAncestor
	}

	signers, err := api.bsrr.getSigners(api.chain, target)
	if err != nil {
		return nil, err
	}

	return signers, nil
}

/*
[BERITH]
A function that returns the probability of Block Creator election
*/
func (api *API) GetJoinRatio(address common.Address, number *rpc.BlockNumber) (float64, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	var num int64
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return the signers from its snapshot
	if header == nil {
		return 0, errUnknownBlock
	}

	num = header.Number.Int64()

	epoch := int64(api.bsrr.config.Epoch)
	if num <= epoch {
		return 0, errNoData
	}

	stks, err := api.bsrr.getStakers(api.chain, uint64(num), header.Hash())
	if err != nil {
		return 0, err
	}

	parent := api.chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return 0, consensus.ErrUnknownAncestor
	}

	target, exist := api.bsrr.getStakeTargetBlock(api.chain, parent)
	if !exist {
		return 0, consensus.ErrUnknownAncestor
	}

	states, err := api.chain.StateAt(target.Root)
	if err != nil {
		return 0, err
	}

	roi, err := api.bsrr.getJoinRatio(stks, address, header.Hash(), uint64(num), states)
	if err != nil {
		return 0, err
	}

	return roi, nil
}

// GetSignersAtHash retrieves the list of authorized signers at the specified block.
func (api *API) GetSignersAtHash(hash common.Hash) ([]common.Address, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}

	parent := api.chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return nil, consensus.ErrUnknownAncestor
	}

	target, exist := api.bsrr.getStakeTargetBlock(api.chain, parent)
	if !exist {
		return nil, consensus.ErrUnknownAncestor
	}

	signers, err := api.bsrr.getSigners(api.chain, target)
	if err != nil {
		return nil, err
	}
	return signers, nil
}

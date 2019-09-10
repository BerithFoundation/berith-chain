/*
d8888b. d88888b d8888b. d888888b d888888b db   db
88  `8D 88'     88  `8D   `88'   `~~88~~' 88   88
88oooY' 88ooooo 88oobY'    88       88    88ooo88
88~~~b. 88~~~~~ 88`8b      88       88    88~~~88
88   8D 88.     88 `88.   .88.      88    88   88
Y8888P' Y88888P 88   YD Y888888P    YP    YP   YP

	  copyrights by ibizsoftware 2018 - 2019
*/

package bsrr

import (
	"errors"

	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/consensus"
	"github.com/BerithFoundation/berith-chain/core/types"
	"github.com/BerithFoundation/berith-chain/rpc"
)

// API is a user facing RPC API to allow controlling the signer and voting
// mechanisms of the proof-of-authority scheme.
type API struct {
	chain consensus.ChainReader
	bsrr  *BSRR
}

//// GetSnapshot retrieves the state snapshot at a given block.
//func (api *API) GetSnapshot(number *rpc.BlockNumber) (*Snapshot, error) {
//	// Retrieve the requested block number (or current if none requested)
//	var header *types.Header
//	if number == nil || *number == rpc.LatestBlockNumber {
//		header = api.chain.CurrentHeader()
//	} else {
//		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
//	}
//	// Ensure we have an actually valid block and return its snapshot
//	if header == nil {
//		return nil, errUnknownBlock
//	}
//	return api.bsrr.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
//}
//
//// GetSnapshotAtHash retrieves the state snapshot at a given block.
//func (api *API) GetSnapshotAtHash(hash common.Hash) (*Snapshot, error) {
//	header := api.chain.GetHeaderByHash(hash)
//	if header == nil {
//		return nil, errUnknownBlock
//	}
//	return api.bsrr.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
//}

func (api *API) GetState(number *rpc.BlockNumber) (string, error) {

	header := api.chain.GetHeaderByNumber(uint64(number.Int64()))
	if header == nil {
		return "", errors.New("No search a block")
	}
	state, err := api.chain.StateAt(header.Root)

	if err != nil {
		return "", err
	}

	users, err := api.GetBlockCreators(number)

	if err != nil {
		return "", err
	}

	result := "{\n\tHASH : " + header.Hash().Hex() + ", "
	result += "\n\tROOT : " + header.Root.Hex() + ", "
	result += "\n\tRESULTS : ["

	for _, user := range users {

		info := state.GetAccountInfo(user)
		result += "\n\t\t{"
		result += "\n\t\t\tCODEHASH : " + common.Bytes2Hex(info.CodeHash) + ", "
		result += "\n\t\t\tNONCE : " + string(info.Nonce) + ", "
		result += "\n\t\t\tROOT : " + info.Root.Hex() + ", "
		result += "\n\t\t\tMAIN : " + info.Balance.String() + ", "
		result += "\n\t\t\tPOINT : " + info.Point.String() + ", "
		result += "\n\t\t\tSTAKE : " + info.StakeBalance.String()
		result += "\n\t\t}, "
	}
	result += "\n\t]\n}"

	return result, nil

}

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

	target, exist := api.bsrr.getAncestor(api.chain, int64(api.bsrr.config.Epoch), parent)

	if !exist {
		return nil, consensus.ErrUnknownAncestor
	}

	signers, err := api.bsrr.getSigners(api.chain, target)

	if err != nil {
		return nil, err
	}

	signersMap := signers.signersMap()

	result := make([]common.Address, 0)
	for k, _ := range signersMap {
		result = append(result, k)
	}
	return result, nil
}

// GetSigners retrieves the list of authorized signers at the specified block.
func (api *API) GetSigners(number *rpc.BlockNumber) ([]common.Address, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return the signers from its snapshot
	if header == nil {
		return nil, errUnknownBlock
	}

	parent := api.chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return nil, consensus.ErrUnknownAncestor
	}

	target, exist := api.bsrr.getAncestor(api.chain, int64(api.bsrr.config.Epoch), parent)

	if !exist {
		return nil, consensus.ErrUnknownAncestor
	}

	signers, err := api.bsrr.getSigners(api.chain, target)
	//snap, err := api.bsrr.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
	if err != nil {
		return nil, err
	}
	//return snap.signers(), nil
	return signers, nil
}

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

	//p := header.ParentHash
	//uint64(num - epoch - 2)
	//	prev_header := api.chain.GetHeaderByNumber(5740)

	stakingList, err := api.bsrr.getStakingList(api.chain, uint64(num), header.Hash())
	if err != nil {
		return 0, err
	}

	parent := api.chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return 0, consensus.ErrUnknownAncestor
	}

	target, exist := api.bsrr.getAncestor(api.chain, int64(api.bsrr.config.Epoch), parent)

	if !exist {
		return 0, consensus.ErrUnknownAncestor
	}

	states, err := api.chain.StateAt(target.Root)

	if err != nil {
		return 0, err
	}

	roi, err := api.bsrr.getJoinRatio(&stakingList, address, uint64(num), states)
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

	target, exist := api.bsrr.getAncestor(api.chain, int64(api.bsrr.config.Epoch), parent)

	if !exist {
		return nil, consensus.ErrUnknownAncestor
	}

	signers, err := api.bsrr.getSigners(api.chain, target)
	//snap, err := api.bsrr.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
	if err != nil {
		return nil, err
	}
	return signers, nil
}

// Proposals returns the current proposals the node tries to uphold and vote on.
func (api *API) Proposals() map[common.Address]bool {
	api.bsrr.lock.RLock()
	defer api.bsrr.lock.RUnlock()

	proposals := make(map[common.Address]bool)
	for address, auth := range api.bsrr.proposals {
		proposals[address] = auth
	}
	return proposals
}

// Propose injects a new authorization proposal that the signer will attempt to
// push through.
func (api *API) Propose(address common.Address, auth bool) {
	api.bsrr.lock.Lock()
	defer api.bsrr.lock.Unlock()

	api.bsrr.proposals[address] = auth
}

// Discard drops a currently running proposal, stopping the signer from casting
// further votes (either for or against).
func (api *API) Discard(address common.Address) {
	api.bsrr.lock.Lock()
	defer api.bsrr.lock.Unlock()

	delete(api.bsrr.proposals, address)
}

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
	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/consensus"
	"github.com/BerithFoundation/berith-chain/core/types"
	"github.com/BerithFoundation/berith-chain/rpc"
)

// API is a user facing RPC API to allow controlling the signer and voting
// mechanisms of the proof-of-authority scheme.
type API struct {
	chain  consensus.ChainReader
	bsrr *BSRR
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

	signers, err := api.bsrr.getSigners(api.chain, header.Number.Uint64(), header.Hash())
	//snap, err := api.bsrr.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
	if err != nil {
		return nil, err
	}
	//return snap.signers(), nil
	return signers, nil
}

func (api *API) GetRoundJoinRatio(address common.Address, number *rpc.BlockNumber) (int, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return the signers from its snapshot
	if header == nil {
		return 0, errUnknownBlock
	}

	stakingList, err := api.bsrr.getStakingList(api.chain, header.Number.Uint64(), header.Hash())
	if err != nil {
		return 0, err
	}

	ratio, err := api.bsrr.roundJoinRatio(&stakingList, address)
	if err != nil {
		return 0, err
	}

	return ratio, nil
}

// GetSignersAtHash retrieves the list of authorized signers at the specified block.
func (api *API) GetSignersAtHash(hash common.Hash) ([]common.Address, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}

	signers, err := api.bsrr.getSigners(api.chain, header.Number.Uint64(), header.Hash())
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

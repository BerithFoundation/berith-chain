// Copyright 2019 The berith Authors
// This file is part of berith.
//
// berith is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// berith is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with berith. If not, see <http://www.gnu.org/licenses/>.
package main

import (
	"fmt"
	berith_chain "github.com/BerithFoundation/berith-chain"
	"github.com/BerithFoundation/berith-chain/accounts"
	"github.com/BerithFoundation/berith-chain/accounts/keystore"
	"github.com/BerithFoundation/berith-chain/berithclient"
	"math/big"
	"sync"
	"time"
)

type berithenchContext struct {
	TestId           string             // test id "berithench-{15:04:05.999}"
	ChainID          *big.Int           // chain id
	NodeContexts     []*nodeContext     // node contexts
	Keystore         *keystore.KeyStore // keystore
	AddressContexts  []*addressContext  // addresses of sending tx
	Duration         time.Duration      // time duration for test execution
	TxCount          uint64             // tx count of test execution
	TxInterval       uint64             // interval of sending transactions
	InitDelay        uint64             // delay of initial
	OutputPath       string             // output dir path to write a summary file
	EnableCpuProfile bool               // enable start cpu profile
	EnableGoTrace    bool               // enable start go trace
}

type addressContext struct {
	account    accounts.Account
	startNonce uint64 // start nonce before testing
	lastNonce  uint64 // last nonce of testing
	summary    *summaryContext
}

type nodeContext struct {
	client  *berithclient.Client // client of a node
	url     string               // node rpc url
	summary *summaryContext

	quitTask        chan struct{} // flag for quitting task
	startCpuProfile bool
	startGoTrace    bool
	subscription    berith_chain.Subscription
}

type summaryContext struct {
	try     uint // count of request
	success uint // count of success
	fail    uint // count of fail

	firstTxHash string // hash of a first transaction

	mutex *sync.Mutex
}

/////////// addressContext

// NewAddressContext new instance of address context with default values
func NewAddressContext(account accounts.Account) *addressContext {
	return &addressContext{
		account:    account,
		startNonce: 0,
		lastNonce:  0,
		summary: &summaryContext{
			mutex: &sync.Mutex{},
		},
	}
}

/////////// nodeContext

// NewNodeContext new instance of node context given url
func NewNodeContext(url string) *nodeContext {
	return &nodeContext{
		url: url,
		summary: &summaryContext{
			mutex: &sync.Mutex{},
		},
	}
}

// retryFailTransaction retry to send transaction until success to request or included in block
func (n *nodeContext) retryFailTransaction() {

}

// subscribeBlocks subscribe new blocks.
// start to a new go routine for getting a new block with polling interval i.e blockTime
// because http client is not supported SubscribeNewHead function,
// FIXME : if support http subscription, then have to change
func (n *nodeContext) subscribeBlocks(blockTime int) {
	//if n.client == nil {
	//	utils.Fatalf("must exist node.client before subscription")
	//}
	//
	//ctx := context.Background()
	//ch := make(chan<- *types.Header)
	//sub, err := n.client.SubscribeNewHead(ctx, ch)
	//if err != nil {
	//	utils.Fatalf("failed to subscribe new head:%v", err)
	//}
	//n.subscription = sub
}

// close clear up node context
func (n *nodeContext) close() {
	if n.client != nil {
		n.client.Close()
	}
	if n.subscription != nil {
		n.subscription.Unsubscribe()
	}
}

/////////// summaryContext

func (r *summaryContext) String() string {
	return fmt.Sprintf("Request : %d / Success : %d / Fail : %d", r.try, r.success, r.fail)
}

// addResult record request summary such as counts of try, success, fail
func (r *summaryContext) addResult(success bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.try++
	if success {
		r.success++
	} else {
		r.fail++
	}
}

// returns summary context with record
func (r *summaryContext) getSummary() summaryContext {
	r.mutex.Lock()
	defer r.mutex.Lock()
	return summaryContext{
		try:     r.try,
		success: r.success,
		fail:    r.fail,
		mutex:   nil,
	}
}

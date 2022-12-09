// Copyright 2014 The go-ethereum Authors
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

package state

import (
	"bytes"
	"fmt"
	"io"
	"math/big"

	"github.com/pkg/errors"

	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/crypto"
	"github.com/BerithFoundation/berith-chain/rlp"
)

var emptyCodeHash = crypto.Keccak256(nil)

type Code []byte

func (s Code) String() string {
	return string(s) //strings.Join(Disassemble(s), " ")
}

type Storage map[common.Hash]common.Hash

func (s Storage) String() (str string) {
	for key, value := range s {
		str += fmt.Sprintf("%X : %X\n", key, value)
	}

	return
}

func (s Storage) Copy() Storage {
	cpy := make(Storage)
	for key, value := range s {
		cpy[key] = value
	}

	return cpy
}

// stateObject represents an Ethereum account which is being modified.
// The usage pattern is as follows:
// First you need to obtain a state object.
// Finally, call CommitTrie to write the modified storage trie into a database.
// Account values can be accessed and modified through the object.
//
// stateObjecct는 수정중인 베리드 계정을 대변한다
// 먼저 state object를 얻어야한다.
// 계정값은 객체를 통해 접근되고 수정될 수 있다.
// 마지막으로 CommitTrie를 호출하여 수정된 스토리지 트리를 DB에 기록한다.
type stateObject struct {
	address  common.Address
	addrHash common.Hash // hash of berith address of the account
	data     Account
	db       *StateDB

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	// Write caches.
	trie Trie // storage trie, which becomes non-nil on first access
	code Code // contract bytecode, which gets set when code is loaded

	// 마지막 stateCommit에 따른 Storage 상태
	originStorage Storage // Storage cache of original entries to dedup rewrites

	// 수정된 상태를 저장한다.
	dirtyStorage Storage // Storage entries that need to be flushed to disk

	// Cache flags.
	// When an object is marked suicided it will be delete from the trie
	// during the "update" phase of the state transition.
	dirtyCode bool // true if the code was updated
	suicided  bool
	deleted   bool
}

// empty returns whether the account is considered empty.
func (s *stateObject) empty() bool {
	return s.data.Nonce == 0 && s.data.Balance.Sign() == 0 && bytes.Equal(s.data.CodeHash, emptyCodeHash)
}

/*
[BERITH]
Account is the Berith consensus representation of accounts.
These objects are stored in the main account trie.
Add StakeBalance, BehindBalance, Selection Point

Account는 Berith에서 합의된 계정 표현이다.
이 객체는 메인 어카운트 트리에 저장되고
StakeBalance, BehindBalance, SelectionPoint를 더한다.
*/
type Account struct {
	Nonce          uint64
	Balance        *big.Int
	Root           common.Hash // merkle root of the storage trie
	CodeHash       []byte      // 스마트 컨트랙트 바이트 코드의 해시
	StakeBalance   *big.Int    //brt staking balance
	StakeUpdated   *big.Int    //Block number when the stake balance was updated
	Point          *big.Int    //selection Point, 스테이킹에 대한 Point
	BehindBalance  []Behind    //behind balance
	Penalty        uint64
	PenlatyUpdated *big.Int //Block Number when the penalty was updated
}

/*
[BERITH]
Balance for Reward Payment
*/
type Behind struct {
	Number  *big.Int
	Balance *big.Int
}

/*
[BERITH]
Function to create stateObject object
Account initialization processing
*/
func newObject(db *StateDB, address common.Address, data Account) *stateObject {
	if data.Balance == nil {
		data.Balance = new(big.Int)
	}
	if data.CodeHash == nil {
		data.CodeHash = emptyCodeHash
	}
	if data.StakeBalance == nil {
		data.StakeBalance = new(big.Int)
	}

	if data.StakeUpdated == nil {
		data.StakeUpdated = new(big.Int)
	}
	if data.Point == nil {
		data.Point = new(big.Int)
	}

	if data.BehindBalance == nil {
		data.BehindBalance = make([]Behind, 0)
	}

	if data.PenlatyUpdated == nil {
		data.PenlatyUpdated = new(big.Int)
	}

	return &stateObject{
		db:            db,
		address:       address,
		addrHash:      crypto.Keccak256Hash(address[:]),
		data:          data,
		originStorage: make(Storage),
		dirtyStorage:  make(Storage),
	}
}

// EncodeRLP implements rlp.Encoder.
func (c *stateObject) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, c.data)
}

// setError remembers the first non-nil error it is called with.
func (s *stateObject) setError(err error) {
	if s.dbErr == nil {
		s.dbErr = err
	}
}

func (s *stateObject) markSuicided() {
	s.suicided = true
}

func (c *stateObject) touch() {
	c.db.journal.append(touchChange{
		account: &c.address,
	})
	if c.address == ripemd {
		// Explicitly put it in the dirty-cache, which is otherwise generated from
		// flattened journals.
		c.db.journal.dirty(c.address)
	}
}

func (c *stateObject) getTrie(db Database) Trie {
	if c.trie == nil {
		var err error
		c.trie, err = db.OpenStorageTrie(c.addrHash, c.data.Root)
		if err != nil {
			c.trie, _ = db.OpenStorageTrie(c.addrHash, common.Hash{})
			c.setError(fmt.Errorf("can't create storage trie: %v", err))
		}
	}
	return c.trie
}

// GetState retrieves a value from the account storage trie.
func (s *stateObject) GetState(db Database, key common.Hash) common.Hash {
	// If we have a dirty value for this state entry, return it
	value, dirty := s.dirtyStorage[key]
	if dirty {
		return value
	}
	// Otherwise return the entry's original value
	return s.GetCommittedState(db, key)
}

// GetCommittedState retrieves a value from the committed account storage trie.
func (s *stateObject) GetCommittedState(db Database, key common.Hash) common.Hash {
	// If we have the original value cached, return that
	value, cached := s.originStorage[key]
	if cached {
		return value
	}
	// Otherwise load the value from the database
	enc, err := s.getTrie(db).TryGet(key[:])
	if err != nil {
		s.setError(err)
		return common.Hash{}
	}
	if len(enc) > 0 {
		_, content, _, err := rlp.Split(enc)
		if err != nil {
			s.setError(err)
		}
		value.SetBytes(content)
	}
	s.originStorage[key] = value
	return value
}

// SetState updates a value in account storage.
func (s *stateObject) SetState(db Database, key, value common.Hash) {
	// If the new value is the same as old, don't set
	prev := s.GetState(db, key)
	if prev == value {
		return
	}
	// New value is different, update and journal the change
	s.db.journal.append(storageChange{
		account:  &s.address,
		key:      key,
		prevalue: prev,
	})
	s.setState(key, value)
}

func (s *stateObject) setState(key, value common.Hash) {
	s.dirtyStorage[key] = value
}

// updateTrie writes cached storage modifications into the object's storage trie.
func (s *stateObject) updateTrie(db Database) Trie {
	tr := s.getTrie(db)
	for key, value := range s.dirtyStorage {
		delete(s.dirtyStorage, key)

		// Skip noop changes, persist actual changes
		if value == s.originStorage[key] {
			continue
		}
		s.originStorage[key] = value

		if (value == common.Hash{}) {
			s.setError(tr.TryDelete(key[:]))
			continue
		}
		// Encoding []byte cannot fail, ok to ignore the error.
		v, _ := rlp.EncodeToBytes(bytes.TrimLeft(value[:], "\x00"))
		s.setError(tr.TryUpdate(key[:], v))
	}
	return tr
}

// UpdateRoot sets the trie root to the current root hash of
func (s *stateObject) updateRoot(db Database) {
	s.updateTrie(db)
	s.data.Root = s.trie.Hash()
}

// CommitTrie the storage trie of the object to db.
// This updates the trie root.
func (s *stateObject) CommitTrie(db Database) error {
	s.updateTrie(db)
	if s.dbErr != nil {
		return s.dbErr
	}
	root, err := s.trie.Commit(nil)
	if err == nil {
		s.data.Root = root
	}
	return err
}

// AddBalance removes amount from c's balance.
// It is used to add funds to the destination account of a transfer.
func (c *stateObject) AddBalance(amount *big.Int) {
	// EIP158: We must check emptiness for the objects such that the account
	// clearing (0,0,0 objects) can take effect.
	if amount.Sign() == 0 {
		if c.empty() {
			c.touch()
		}

		return
	}
	c.SetBalance(new(big.Int).Add(c.Balance(), amount))
}

// SubBalance removes amount from c's balance.
// It is used to remove funds from the origin account of a transfer.
func (c *stateObject) SubBalance(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	c.SetBalance(new(big.Int).Sub(c.Balance(), amount))
}

func (s *stateObject) SetBalance(amount *big.Int) {
	s.db.journal.append(balanceChange{
		account: &s.address,
		prev:    new(big.Int).Set(s.data.Balance),
	})
	s.setBalance(amount)
}

func (s *stateObject) setBalance(amount *big.Int) {
	s.data.Balance = amount
}

// Return the gas back to the origin. Used by the Virtual machine or Closures
func (c *stateObject) ReturnGas(gas *big.Int) {}

func (s *stateObject) deepCopy(db *StateDB) *stateObject {
	stateObject := newObject(db, s.address, s.data)
	if s.trie != nil {
		stateObject.trie = db.db.CopyTrie(s.trie)
	}
	stateObject.code = s.code
	stateObject.dirtyStorage = s.dirtyStorage.Copy()
	stateObject.originStorage = s.originStorage.Copy()
	stateObject.suicided = s.suicided
	stateObject.dirtyCode = s.dirtyCode
	stateObject.deleted = s.deleted
	return stateObject
}

// Returns the address of the contract/account
func (c *stateObject) Address() common.Address {
	return c.address
}

// Code returns the contract code associated with this object, if any.
func (s *stateObject) Code(db Database) []byte {
	if s.code != nil {
		return s.code
	}
	if bytes.Equal(s.CodeHash(), emptyCodeHash) {
		return nil
	}
	code, err := db.ContractCode(s.addrHash, common.BytesToHash(s.CodeHash()))
	if err != nil {
		s.setError(fmt.Errorf("can't load code hash %x: %v", s.CodeHash(), err))
	}
	s.code = code
	return code
}

func (s *stateObject) SetCode(codeHash common.Hash, code []byte) {
	prevcode := s.Code(s.db.db)
	s.db.journal.append(codeChange{
		account:  &s.address,
		prevhash: s.CodeHash(),
		prevcode: prevcode,
	})
	s.setCode(codeHash, code)
}

func (s *stateObject) setCode(codeHash common.Hash, code []byte) {
	s.code = code
	s.data.CodeHash = codeHash[:]
	s.dirtyCode = true
}

func (s *stateObject) SetNonce(nonce uint64) {
	s.db.journal.append(nonceChange{
		account: &s.address,
		prev:    s.data.Nonce,
	})
	s.setNonce(nonce)
}

func (s *stateObject) setNonce(nonce uint64) {
	s.data.Nonce = nonce
}

func (s *stateObject) CodeHash() []byte {
	return s.data.CodeHash
}

func (s *stateObject) Balance() *big.Int {
	return s.data.Balance
}

func (s *stateObject) Nonce() uint64 {
	return s.data.Nonce
}

// Never called, but must be present to allow stateObject to be used
// as a vm.Account interface that also satisfies the vm.ContractRef
// interface. Interfaces are awesome.
func (s *stateObject) Value() *big.Int {
	panic("Value on stateObject should never be called")
}

/*
[BERITH]
set staking balance
*/
func (s *stateObject) SetStaking(amount, blockNumber *big.Int) {
	s.db.journal.append(stakingChange{
		account:     &s.address,
		prevBalance: new(big.Int).Set(s.data.StakeBalance),
		prevBlock:   new(big.Int).Set(s.data.StakeUpdated),
	})
	s.setStaking(amount, blockNumber)
}

func (s *stateObject) setStaking(amount, blockNumber *big.Int) {
	s.data.StakeBalance = amount
	s.data.StakeUpdated = blockNumber
}

func (s *stateObject) StakeBalance() *big.Int {
	return s.data.StakeBalance
}

func (s *stateObject) StakeUpdated() *big.Int {
	return s.data.StakeUpdated
}

func (c *stateObject) RemoveStakeBalance() {
	stakeBalance := c.StakeBalance()
	if stakeBalance.Sign() == 0 {
		return
	}

	c.SetStaking(common.Big0, c.data.StakeUpdated)
	c.AddBalance(stakeBalance)
}

func (c *stateObject) AddStakeBalance(amount, blockNumber *big.Int) {
	// EIP158: We must check emptiness for the objects such that the account
	// clearing (0,0,0 objects) can take effect.
	if amount.Sign() == 0 {
		if c.empty() {
			c.touch()
		}
		return
	}
	c.SetStaking(new(big.Int).Add(c.StakeBalance(), amount), blockNumber)
}

/*
[BERITH]
BehindBalance values are arrays
Function to add Behind object including block number and coin quantity to the array
*/
func (c *stateObject) AddBehindBalance(number, amount *big.Int) {
	// EIP158: We must check emptiness for the objects such that the account
	// clearing (0,0,0 objects) can take effect.
	if amount.Sign() == 0 {
		if c.empty() {
			c.touch()
		}
		return
	}
	c.SetBehind(number, amount)
}

func (s *stateObject) SetBehind(number, amount *big.Int) {

	ch := behindChange{}
	ch.account = &s.address

	behind := Behind{}
	behind.Number = number
	behind.Balance = amount

	ch.prev = append(ch.prev, behind)
	s.db.journal.append(ch)

	s.setBehind(append(s.data.BehindBalance, behind))
}

func (s *stateObject) setBehind(behind []Behind) {
	s.data.BehindBalance = behind
}

func (s *stateObject) BehindBalance() []Behind {
	return s.data.BehindBalance
}

/*
[BERITH]
Function that returns 0th in BehindBalance array
Function to process FIFO
*/
func (s *stateObject) GetFirstBehindBalance() (Behind, error) {
	behind := s.data.BehindBalance
	if behind == nil || len(behind) == 0 {
		return Behind{}, errors.New("nil behind")
	}

	return behind[0], nil
}

/*
[BERITH]
Function to delete 0th value from BehindBalance array
Function to process FIFO
*/
func (s *stateObject) RemoveFirstBehindBalance() {
	behind := s.data.BehindBalance
	s.setBehind(behind[1:])
}

/*
[BERITH]
Function to assign the value of Selection Point
*/
func (s *stateObject) SetPoint(amount *big.Int) {
	s.db.journal.append(pointChange{
		account: &s.address,
		prev:    new(big.Int).Set(s.data.Point),
	})
	s.setPoint(amount)
}

func (s *stateObject) setPoint(amount *big.Int) {
	s.data.Point = amount
}

/*
[BERITH]
Function that returns Selection Point
*/
func (s *stateObject) GetPoint() *big.Int {
	return s.data.Point
}

/*
[BERITH]
Function to add the corresponding amount value to Selection Point value
*/
func (c *stateObject) AddPoint(amount *big.Int) {
	// EIP158: We must check emptiness for the objects such that the account
	// clearing (0,0,0 objects) can take effect.
	if amount.Sign() == 0 {
		if c.empty() {
			c.touch()
		}
		return
	}
	c.SetPoint(new(big.Int).Add(c.GetPoint(), amount))
}

/*
[BERITH]
Function that returns Account information
*/
func (s *stateObject) AccountInfo() Account {
	return s.data
}

// [BERITH] Penalty-related function definition

// [BERITH] Returns the current penalty value
func (s *stateObject) Penalty() uint64 {
	return s.data.Penalty
}

// [BERITH] Returns the block number where the penalty was last changed
func (s *stateObject) PenaltyUpdated() *big.Int {
	return s.data.PenlatyUpdated
}

// [BERITH] Function to increase penalty by 1
func (s *stateObject) AddPenalty(blockNumber *big.Int) {
	s.db.journal.append(penaltyChange{
		account:     &s.address,
		prevPenalty: s.data.Penalty,
		prevBlock:   s.data.PenlatyUpdated,
	})
	s.setPenalty(s.data.Penalty+1, blockNumber)
}

// [BERITH] Function to remove penalty for a specific account
func (s *stateObject) RemovePenalty(blockNumber *big.Int) {
	s.db.journal.append(penaltyChange{
		account:     &s.address,
		prevPenalty: s.data.Penalty,
		prevBlock:   s.data.PenlatyUpdated,
	})

	s.setPenalty(0, blockNumber)
}

func (s *stateObject) setPenalty(amount uint64, blockNumber *big.Int) {
	s.data.Penalty = amount
	s.data.PenlatyUpdated = blockNumber
}

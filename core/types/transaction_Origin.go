package types

import (
	"errors"
	"io"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/crypto"
	"github.com/BerithFoundation/berith-chain/rlp"
)

type TransactionInterface interface {
	ChainId() *big.Int
	Protected() bool
	EncodeRLP(w io.Writer) error
	DecodeRLP(s *rlp.Stream) error
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(input []byte) error
	Data() []byte
	Gas() uint64
	GasPrice() *big.Int
	Value() *big.Int
	Nonce() uint64
	CheckNonce() bool
	Base() JobWallet
	Target() JobWallet
	To() *common.Address
	Hash() common.Hash
	Size() common.StorageSize
	AsMessage(s Signer) (Message, error)
	WithSignature(signer Signer, sig []byte) (*Transaction, error)
	Cost() *big.Int
	MainFee() *big.Int
	RawSignatureValues() (*big.Int, *big.Int, *big.Int)
	From() *atomic.Value
	IsEthTransaction() bool
}

type OriginTransaction struct {
	inner *LegacyTx
	time  time.Time // Time first seen locally (spam avoidance)

	// caches
	hash    atomic.Value
	size    atomic.Value
	from    atomic.Value
	IsEthTx bool
}

// LegacyTx is the transaction data of regular Ethereum transactions.
type LegacyTx struct {
	Nonce    uint64          // nonce of sender account
	GasPrice *big.Int        // wei per gas
	Gas      uint64          // gas limit
	To       *common.Address `rlp:"nil"` // nil means contract creation
	Value    *big.Int        // wei amount
	Data     []byte          // contract invocation input data
	V, R, S  *big.Int        // signature values
}

func (o *LegacyTx) ToTxdata() txdata {
	return txdata{
		AccountNonce: o.Nonce,
		Price:        o.GasPrice,
		GasLimit:     o.Gas,
		Recipient:    o.To,
		Amount:       o.Value,
		Payload:      o.Data,
		Base:         JobWallet(1),
		Target:       JobWallet(1),
		V:            o.V,
		R:            o.R,
		S:            o.S,
		Hash:         &common.Hash{}}

}

// MarshalJSON encodes the web3 RPC transaction format.
func (o *OriginTransaction) MarshalJSON() ([]byte, error) {
	inner := o.inner
	return inner.MarshalJSON()
}

// UnmarshalJSON decodes the web3 RPC transaction format.
func (o *OriginTransaction) UnmarshalJSON(input []byte) error {
	var dec LegacyTx
	if err := dec.UnmarshalJSON(input); err != nil {
		return err
	}

	withSignature := dec.V.Sign() != 0 || dec.R.Sign() != 0 || dec.S.Sign() != 0
	if withSignature {
		var V byte
		if isProtectedV(dec.V) {
			chainID := deriveChainId(dec.V).Uint64()
			V = byte(dec.V.Uint64() - 35 - 2*chainID)
		} else {
			V = byte(dec.V.Uint64() - 27)
		}
		if !crypto.ValidateSignatureValues(V, dec.R, dec.S, false) {
			return ErrInvalidSig
		}
	}

	*o = OriginTransaction{inner: &dec}
	return nil
}

// UnmarshalBinary decodes the canonical encoding of transactions.
// It supports legacy RLP transactions and EIP2718 typed transactions.
func (tx *OriginTransaction) UnmarshalBinary(b []byte) error {
	if len(b) > 0 && b[0] > 0x7f {
		// It's a legacy transaction.
		var data LegacyTx
		err := rlp.DecodeBytes(b, &data)
		if err != nil {
			return err
		}
		tx.setDecoded(&data, len(b))
		return nil
	} else {
		return errors.New("cannot unmarshal binary : not berith transaction type")
	}
}

// setDecoded sets the inner transaction and size after decoding.
func (tx *OriginTransaction) setDecoded(inner *LegacyTx, size int) {
	tx.inner = inner
	tx.time = time.Now()
	// [Berith]
	// check original eth transaction
	tx.IsEthTx = true
	if size > 0 {
		tx.size.Store(common.StorageSize(size))
	}
}

func (o *OriginTransaction) ChainId() *big.Int {
	return deriveChainId(o.inner.V)
}

// Protected returns whether the transaction is protected from replay protection.
func (o *OriginTransaction) Protected() bool {
	return isProtectedV(o.inner.V)
}

func (o *OriginTransaction) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &o.inner)
}

// DecodeRLP implements rlp.Decoder
func (o *OriginTransaction) DecodeRLP(s *rlp.Stream) error {
	_, size, _ := s.Kind()
	err := s.Decode(&o.inner)
	if err == nil {
		o.size.Store(common.StorageSize(rlp.ListSize(size)))
	}
	o.IsEthTx = true
	return err
}

func (o *OriginTransaction) Data() []byte       { return common.CopyBytes(o.inner.Data) }
func (o *OriginTransaction) Gas() uint64        { return o.inner.Gas }
func (o *OriginTransaction) GasPrice() *big.Int { return new(big.Int).Set(o.inner.GasPrice) }
func (o *OriginTransaction) Value() *big.Int    { return new(big.Int).Set(o.inner.Value) }
func (o *OriginTransaction) Nonce() uint64      { return o.inner.Nonce }
func (o *OriginTransaction) CheckNonce() bool   { return true }
func (o *OriginTransaction) Base() JobWallet    { return JobWallet(1) } //[Berith] Tx JobWallet Base
func (o *OriginTransaction) Target() JobWallet  { return JobWallet(1) } //[Berith] Tx JobWallet Target

func (o *OriginTransaction) From() *atomic.Value {
	return &o.from
}
func (o *OriginTransaction) IsEthTransaction() bool { return o.IsEthTx } //[Berith] Tx JobWallet Target

// To returns the recipient address of the transaction.
// It returns nil if the transaction is a contract creation.
func (o *OriginTransaction) To() *common.Address {
	if o.inner.To == nil {
		return nil
	}
	to := *o.inner.To
	return &to
}

// Hash hashes the RLP encoding of o.
// It uniquely identifies the transaction.
func (o *OriginTransaction) Hash() common.Hash {
	var txForhash = struct {
		inner *LegacyTx
		time  time.Time // Time first seen locally (spam avoidance)

		// caches
		hash atomic.Value
		size atomic.Value
		from atomic.Value
	}{
		inner: o.inner,
		time:  o.time,
		hash:  o.hash,
		size:  o.size,
		from:  o.from,
	}
	if hash := txForhash.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlpHash(txForhash)
	o.hash.Store(v)
	return v
}

// Size returns the true RLP encoded storage size of the transaction, either by
// encoding and returning it, or returning a previsouly cached value.
func (o *OriginTransaction) Size() common.StorageSize {
	if size := o.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, &o.inner)
	o.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

// AsMessage returns the transaction as a core.Message.
//
// AsMessage requires a signer to derive the sender.
//
// XXX Rename message to something less arbitrary?
func (o *OriginTransaction) AsMessage(s Signer) (Message, error) {
	msg := Message{
		nonce:      o.inner.Nonce,
		gasLimit:   o.inner.Gas,
		gasPrice:   new(big.Int).Set(o.inner.GasPrice),
		to:         o.inner.To,
		amount:     o.inner.Value,
		data:       o.inner.Data,
		checkNonce: true,
		base:       JobWallet(1),
		target:     JobWallet(1),
	}

	var err error
	msg.from, err = Sender(s, o)
	return msg, err
}

// WithSignature returns a new transaction with the given signature.
// This signature needs to be formatted as described in the yellow paper (v+27).
func (o *OriginTransaction) WithSignature(signer Signer, sig []byte) (*Transaction, error) {
	r, s, v, err := signer.SignatureValues(o, sig)
	if err != nil {
		return nil, err
	}
	cpy := &Transaction{data: o.inner.ToTxdata()}
	cpy.data.R, cpy.data.S, cpy.data.V = r, s, v
	return cpy, nil
}

// Cost returns amount + gasprice * gaslimit.
func (o *OriginTransaction) Cost() *big.Int {
	total := new(big.Int).Mul(o.inner.GasPrice, new(big.Int).SetUint64(o.inner.Gas))
	total.Add(total, o.inner.Value)
	return total
}

func (o *OriginTransaction) MainFee() *big.Int {
	total := new(big.Int).Mul(o.inner.GasPrice, new(big.Int).SetUint64(o.inner.Gas))
	return total
}

func (o *OriginTransaction) RawSignatureValues() (*big.Int, *big.Int, *big.Int) {
	return o.inner.V, o.inner.R, o.inner.S
}

func NewOriginTransaction(tx *Transaction) *OriginTransaction {
	originTx := &OriginTransaction{
		inner: &LegacyTx{
			Nonce:    tx.data.AccountNonce,
			GasPrice: tx.data.Price,
			Gas:      tx.data.GasLimit,
			To:       tx.data.Recipient,
			Value:    tx.data.Amount,
			Data:     tx.data.Payload,
			V:        tx.data.V,
			R:        tx.data.R,
			S:        tx.data.S},
		hash:    tx.hash,
		size:    tx.size,
		from:    tx.from,
		IsEthTx: true,
	}
	return originTx
}

// [Berith]
// ToBerithTransaction converts transactions sent from metamask into Berith transactions
// with Base and Target set to Main(1) to register with the Berithxpool.
func (o *OriginTransaction) ToBerithTransaction() *Transaction {
	return &Transaction{
		data: txdata{
			AccountNonce: o.inner.Nonce,
			Price:        o.inner.GasPrice,
			GasLimit:     o.inner.Gas,
			Recipient:    o.inner.To,
			Amount:       o.inner.Value,
			Payload:      o.inner.Data,
			Base:         EthTx, // [Berith]
			Target:       EthTx, // [Berith]
			V:            o.inner.V,
			R:            o.inner.R,
			S:            o.inner.S,
			Hash:         &common.Hash{}},
		hash: o.hash,
		size: o.size,
		from: o.from,
		// [Berith]
		// Metamask Transaction에 대해 Base와 Target을 제외하고 Signing하기 위함
		IsEthTx: true,
	}
}

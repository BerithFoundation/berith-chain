package walletdb

import (
	"github.com/BerithFoundation/berith-chain/berithdb"
	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/rlp"
	"io"
	"strings"
)

type Member struct {
	Address    common.Address
	PrivateKey [256]byte
	ID         string
	Password   string
}

type Contact map[common.Address]string

func (c Contact) EncodeRLP(w io.Writer) error {
	result := make([]string,0)

	for k, v := range c {
		result = append(result, k.Hex()+":"+v)
	}
	return rlp.Encode(w, result)
}

func (c Contact) DecodeRLP(r *rlp.Stream) error {
	arr := make([]string, 0)

	err := r.Decode(&arr)

	if err != nil {
		return err
	}

	for _, v := range arr {
		inner := strings.Split(v, ":")
		key := common.HexToAddress(inner[0])
		val := inner[1]

		c[key] = val
	}
	return nil
}


type Transactions struct {
	txs []common.Hash
}

type WalletDB struct {
	db berithdb.Database
}

func NewWalletDB(dir string) (*WalletDB, error) {
	db, err := berithdb.NewLDBDatabase(dir, 128, 1024)

	if err != nil {
		return nil, err
	}

	return &WalletDB{
		db: db,
	}, nil
}

func (db *WalletDB) Insert(key []byte, value interface{}) error {
	data, err := rlp.EncodeToBytes(value)
	if err != nil {
		return err
	}

	return db.db.Put(key, data)
}

func (db *WalletDB) Select(key []byte, holder interface{}) error {
	data, err := db.db.Get(key)
	if err != nil {
		return err
	}

	return rlp.DecodeBytes(data, holder)
}

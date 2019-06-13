package walletdb

import (
	"github.com/BerithFoundation/berith-chain/berithdb"
	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/rlp"
)

type Member struct {
	Address    common.Address
	PrivateKey [256]byte
	ID         string
	Password   string
}

type Contact struct {
	Name    string
	Address common.Address
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

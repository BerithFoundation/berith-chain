package walletdb

import (
	"github.com/BerithFoundation/berith-chain/berithdb"
	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/rlp"
	"io"
	"strings"
)

// 로그인 계정 table ( key 값은 id 로 설정)
type Member struct {
	Address    common.Address // 계정주소
	PrivateKey string		  // 개인키
	ID         string		  // id
	Password   string		  // pwd
}
// 트랜잭션 리스트 Detail table ( key 값은 해당 transaction을 전파받은 블록넘버값)
type TxHistory struct {
	TxAddress common.Address // tx 주소값
	TxType string	// tx 타입 ex ) send, receive , stake 등등
	TxAmount string // tx value 값
	Txtime string // tx 시간
	TxState string // tx 상태
}

// 주소록 table ( key 는 저장 주소값 )
type Contact map[common.Address]string
// 트랜잭션 리스트 Master table ( key 값은 해당 transaction을 전파받은 블록넘버값)
type TxHistoryMaster map[string]string

// Contact table RLP encode 함수
func (c Contact) EncodeRLP(w io.Writer) error {
	result := make([]string,0)

	for k, v := range c {
		result = append(result, k.Hex()+":"+v)
	}
	return rlp.Encode(w, result)
}

// txHistoryMaster table RLP encode 함수
func (c TxHistoryMaster) EncodeRLP(w io.Writer) error {
	result := make([]string,0)

	for k, v := range c {
		result = append(result, k+":"+v)
	}
	return rlp.Encode(w, result)
}

// Contact table RLP decode 함수
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
// txHistoryMaster table RLP decode 함수
func (c TxHistoryMaster) DecodeRLP(r *rlp.Stream) error {
	arr := make([]string, 0)

	err := r.Decode(&arr)

	if err != nil {
		return err
	}

	for _, v := range arr {
		inner := strings.Split(v, ":")
		key := inner[0]
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

// Ldb 생성 함수
func NewWalletDB(dir string) (*WalletDB, error) {
	db, err := berithdb.NewLDBDatabase(dir, 128, 1024)

	if err != nil {
		return nil, err
	}

	return &WalletDB{
		db: db,
	}, nil
}
// db insert 함수
func (db *WalletDB) Insert(key []byte, value interface{}) error {
	data, err := rlp.EncodeToBytes(value)
	if err != nil {
		return err
	}

	return db.db.Put(key, data)
}
// db select 함수
func (db *WalletDB) Select(key []byte, holder interface{}) error {
	data, err := db.db.Get(key)
	if err != nil {
		return err
	}

	return rlp.DecodeBytes(data, holder)
}
//func (db *WalletDB) Select2(key []byte, holder interface{}) error {
//	data, err := db.db.Get(key)
//	if err != nil {
//		return err
//	}
//
//	return rlp.DecodeBytes(data, holder)
//}


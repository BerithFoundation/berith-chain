package walletdb

import (
	"github.com/BerithFoundation/berith-chain/common"
	"testing"
)

func Test01(t *testing.T) {

	member := Member{
		Address:  common.BigToAddress(common.Big257),
		ID:       "kukugi",
		Password: "1234",
	}

	member.PrivateKey[0] = 12

	var contact Contact
	contact = make(map[common.Address]string,0)
	//hexToAddress
	contact[common.BigToAddress(common.Big32)] = "kimmegi"
	contact[common.BigToAddress(common.Big3)] = "gorilla"
	db, err := NewWalletDB("/Users/kimmegi/test.ldb")
	if err != nil {
		t.Error(err)
		return
	}

	err = db.Insert([]byte("swk"), member)
	if err != nil {
		t.Error(err)
		return
	}
	err = db.Insert([]byte("soni"), contact)
	if err != nil {
		t.Error(err)
		return
	}

	db.db.Close()

}
//
//func Test01(t *testing.T) {
//
//	var (
//		err     error
//		member  Member
//		contact = make(Contact,0)
//	)
//
//	db, err := NewWalletDB("/Users/kimmegi/test.ldb")
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = db.Select([]byte("swk"), &member)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = db.Select([]byte("soni"), &contact)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	t.Log(member)
//	t.Log(contact)
//
//	db.db.Close()
//
//}

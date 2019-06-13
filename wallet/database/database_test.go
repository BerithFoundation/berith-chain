package walletdb

import (
	"testing"

	"github.com/BerithFoundation/berith-chain/common"
)

func Test01(t *testing.T) {

	member := Member{
		Address:  common.BigToAddress(common.Big257),
		ID:       "kukugi",
		Password: "1234",
	}

	member.PrivateKey[0] = 12

	contact := Contact{
		Name:    "taejun_bae",
		Address: common.BigToAddress(common.Big32),
	}

	db, err := NewWalletDB("/Users/swk/test.ldb")
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

func Test02(t *testing.T) {

	var (
		err     error
		member  Member
		contact Contact
	)

	db, err := NewWalletDB("/Users/swk/test.ldb")
	if err != nil {
		t.Error(err)
		return
	}

	err = db.Select([]byte("swk"), &member)
	if err != nil {
		t.Error(err)
		return
	}

	err = db.Select([]byte("soni"), &contact)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(member)
	t.Log(contact)

	db.db.Close()

}

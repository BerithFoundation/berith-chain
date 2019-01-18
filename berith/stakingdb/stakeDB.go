package stakingdb

import (
	"fmt"

	"bitbucket.org/ibizsoftware/berith-chain/ethdb"
)

type StakingDB struct {
	stakeDB *ethdb.LDBDatabase
}

/**
DB Create
*/
func (s *StakingDB) CreateDB(filename string) error {
	if s.stakeDB != nil {
		return nil
	}

	db, err := ethdb.NewLDBDatabase(filename, 128, 1024)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	s.stakeDB = db
	return nil
}

/**
DB Get Value
*/
func (s *StakingDB) GetValue(key string) ([]byte, error) {
	k := []byte(key)

	bt, err := s.stakeDB.Get(k)
	if err != nil {
		return nil, err
	}

	return bt, nil
}

/**
DB Insert Value
*/
func (s *StakingDB) PushValue(k string, v []byte) error {
	key := []byte(k)
	return s.stakeDB.Put(key, v)
}

/**
DB Close
*/
func (s *StakingDB) Close() {
	if s.stakeDB == nil {
		return
	}
	s.stakeDB.Close()
}

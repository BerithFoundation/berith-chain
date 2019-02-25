package stakingdb

import (
	"fmt"

	"bitbucket.org/ibizsoftware/berith-chain/berith/staking"
	"bitbucket.org/ibizsoftware/berith-chain/ethdb"
)

type StakingDB struct {
	Encoder encodeFunc
	Decoder decodeFunc
	creator createFunc
	stakeDB *ethdb.LDBDatabase
}

type decodeFunc func(val []byte) (staking.StakingList, error)

type encodeFunc func(list staking.StakingList) ([]byte, error)

type createFunc func() staking.StakingList

/**
DB Create
*/
func (s *StakingDB) CreateDB(filename string, decoder decodeFunc, encoder encodeFunc, creator createFunc) error {
	if s.stakeDB != nil {
		return nil
	}

	db, err := ethdb.NewLDBDatabase(filename, 128, 1024)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	s.stakeDB = db
	s.Encoder = encoder
	s.Decoder = decoder
	s.creator = creator
	return nil
}

/**
DB Get Value
*/
func (s *StakingDB) getValue(key string) ([]byte, error) {
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
func (s *StakingDB) pushValue(k string, v []byte) error {
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

func (s *StakingDB) GetStakingList(key string) (staking.StakingList, error) {
	rlpVal, rlpErr := s.getValue(key)
	if rlpErr != nil {
		return nil, rlpErr
	}

	result, err := s.Decoder(rlpVal)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *StakingDB) Commit(key string, stakingList staking.StakingList) error {
	bytes, err := s.Encoder(stakingList)
	if err != nil {
		return err
	}

	err = s.pushValue(key, bytes)

	if err != nil {
		return err
	}

	return nil
}

func (s *StakingDB) NewStakingList() staking.StakingList {
	return s.creator()
}

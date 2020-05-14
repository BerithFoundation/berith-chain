package staking

import (
	"fmt"

	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/rlp"

	"github.com/BerithFoundation/berith-chain/berithdb"
	"github.com/BerithFoundation/berith-chain/consensus"
	"github.com/BerithFoundation/berith-chain/core/types"
	"github.com/syndtr/goleveldb/leveldb/util"
)

/*
[Berith]
Database that stores staker information
*/
type StakingDB struct {
	creator createFunc
	stakeDB *berithdb.LDBDatabase
}

// staker type creation function
type createFunc func() Stakers

func (s *StakingDB) CreateDB(filename string, creator createFunc) error {
	if s.stakeDB != nil {
		return nil
	}

	db, err := berithdb.NewLDBDatabase(filename, 128, 1024)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	s.stakeDB = db
	s.creator = creator
	return nil
}

/*
[Berith]
Get staker data of a specific block.
*/
func (s *StakingDB) getValue(key string) ([]byte, error) {
	k := []byte(key)

	stakers, err := s.stakeDB.Get(k)
	if err != nil {
		return nil, err
	}
	return stakers, nil
}

/*
[Berith]
Save stakers data in database with block number as key
*/
func (s *StakingDB) pushValue(k string, stakers Stakers) error {
	key := []byte(k)

	v, err := rlp.EncodeToBytes(stakers)

	if err != nil {
		return err
	}

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

/*
[Berith]
After importing the staker data, it is processed into an appropriate data structure and returned.
*/
func (s *StakingDB) GetStakers(key string) (Stakers, error) {
	val, err := s.getValue(key)
	if err != nil {
		return nil, err
	}

	holder := make([]common.Address, 0)
	if err := rlp.DecodeBytes(val, &holder); err != nil {
		return nil, err
	}

	stakers := s.creator()
	stakers.FetchFromList(holder)

	return stakers, nil
}

/*
[Berith]
Save stakers data in database with block number as key
*/
func (s *StakingDB) Commit(key string, value Stakers) error {
	if err := s.pushValue(key, value); err != nil {
		return err
	}
	return nil
}

func (s *StakingDB) NewStakers() Stakers {
	return s.creator()
}

func (s *StakingDB) Clean(chain consensus.ChainReader, header *types.Header) error {
	fmt.Println("Clean stakingDB")

	for {
		key := []byte(header.Hash().Hex())
		exist, err := s.isExist(key)
		if err != nil {
			return err
		}

		if !exist { break }

		err = s.delete(key)
		if err != nil {
			return err
		}
		header = chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	}
	return s.stakeDB.LDB().CompactRange(util.Range{})
}

func (s *StakingDB) isExist(key []byte) (bool, error) {
	return s.stakeDB.Has(key)
}

func (s *StakingDB) delete(key []byte) error {
	return s.stakeDB.Delete(key)
}
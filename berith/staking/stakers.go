package staking

import (
	"io"

	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/rlp"
)

type stakers map[common.Address]struct{}

func (s stakers) Put(addr common.Address) {
	if !s.IsContain(addr) {
		s[addr] = struct{}{}
	}
}

func (s stakers) Remove(addr common.Address) {
	if s.IsContain(addr) {
		delete(s, addr)
	}
}

func (s stakers) IsContain(addr common.Address) bool {
	_, exist := s[addr]
	return exist
}

func (s stakers) AsList() []common.Address {
	result := make([]common.Address, 0)
	for k := range s {
		result = append(result, k)
	}
	return result
}

func (s stakers) FetchFromList(list []common.Address) {
	for _, staker := range list {
		s.Put(staker)
	}
}

func NewStakers() Stakers {
	return make(stakers)
}

func (s stakers) EncodeRLP(w io.Writer) error {
	arr := s.AsList()
	return rlp.Encode(w, arr)
}

func (s stakers) DecodeRLP(stream *rlp.Stream) error {
	holder := make(stakers, 0)
	return stream.Decode(&holder)
}

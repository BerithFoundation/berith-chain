package selection

import (
	"bytes"
	"github.com/BerithFoundation/berith-chain/common"
)

/*
[Berith]
sortableList is a data type for sorting the address list alphabetically.
*/
type sortableList []common.Address

func (s sortableList) Len() int {
	return len(s)
}

func (s sortableList) Swap(a, b int) {
	s[a], s[b] = s[b], s[a]
}

func (s sortableList) Less(a, b int) bool {
	return bytes.Compare(s[a][:], s[b][:]) == -1
}

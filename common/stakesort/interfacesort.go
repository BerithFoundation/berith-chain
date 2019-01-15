package stakesort

import "math/big"

type Stake struct{
	addr interface{}
	value *big.Int
}

type Stakelist []*Stake


func (s Stakelist) Len() int{
	return len(s)
}
func (s Stakelist) Less(i,j int)bool{
	return s[i].value.Uint64() > s[j].value.Uint64()
}
func (s Stakelist)Swap(i,j int){
	s[i],s[j] = s[j],s[i]
}
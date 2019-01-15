package stakesort

import (
	"fmt"
	"math/big"
	"sort"
	"testing"
)

func TestInterfacesort(t *testing.T){

	items := []Stake{
		{"a",big.NewInt(10)},
		{"b",big.NewInt(11)},
		{"c",big.NewInt(20)},
		{"d",big.NewInt(21)},
	}
	//interfacesort := make(Stakelist,0)


	sort.Sort(Stakelist(items))
	fmt.Println(items)

}
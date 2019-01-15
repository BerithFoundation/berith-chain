package stakesort

import (
	"fmt"
	"sort"
	"testing"
)

func TestInterfacesort(t *testing.T){

	items := []Stake{
		{"a",10},
		{"b",11},
		{"c",20},
		{"d",21},
	}
	//interfacesort := make(Stakelist,0)


	sort.Sort(Stakelist(items))
	fmt.Println(items)

}
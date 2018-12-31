package test

import (
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
)

func Test1(t *testing.T) {
	type A struct {
		a string
	}
	type B struct {
		a *A
		b string
	}

	a := &A{a: "aaaa"}
	b := B{a: a, b: "bbbb"}

	data, err := rlp.EncodeToBytes(b)
	if err != nil {
		t.Error("encoding fail : ", err)
	}
	var a2 A
	err = rlp.DecodeBytes(data, &a2)
	if err != nil {
		t.Error("decoding fail : ", err)
	}

	if &a2 == nil {
		t.Error("decoding fail")
	}

}

func Test2(t *testing.T) {
	ch := make(chan bool, 1)

	select {
	case value := <-ch:
		t.Error("value : ", value)
	}
}

type Person struct {
	name string
}

type People []*Person

func (p *Person) Name() string {
	return p.name
}

func (p People) Len() int {
	return len(p)
}

func Test4(t *testing.T) {
	people := make(People, 0)

	people = append(people, &Person{name: "swk"})

	value := reflect.ValueOf(people)

	typ := value.Kind()

	isSlice := typ == reflect.Slice

	t.Error("result : ", isSlice)
}

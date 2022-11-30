/*
[BERITH]
Tx 타입을 지정하기 위한 열거형
*/
package types

type JobWallet uint8

const (
	Main = 1 + iota
	Stake
)


var values = [...]string {
	"main",
	"stake",
}

func (m JobWallet) String() string {
	return values[(m-1)%2]
}

func ConvertJobWallet(s string) JobWallet {
	switch s {
	case "main":
		return Main
	case "stake":
		return Stake
	default:
		return Main
	}
}

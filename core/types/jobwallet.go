package types

type JobWallet uint8

const (
	Main = 1 + iota
	Stake
)


var values = [...]string {
	"main",
	"stake",
	"reward",
}

func (m JobWallet) String() string {
	return values[(m-1)%3]
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

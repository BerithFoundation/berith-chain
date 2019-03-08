package types

type JobWallet uint8

const (
	Main = 1 + iota
	Stake
	Reward
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
	case "reward":
		return Reward
	default:
		return Main
	}
}

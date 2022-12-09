/*
[BERITH]
Enumeration to specify Tx type
*/
package types

import "errors"

type JobWallet uint8

const (
	Main = 1 + iota
	Stake

	// // [Vote] 임시 상수 추가
	// Vote
	end
)

var (
	values = [...]string{
		"main",
		"stake",
		// // [vote] 임시 string
		// "vote",
	}

	ErrInvalidJobWallet = errors.New("invalid wallet type")
)

func (m JobWallet) String() string {
	// // [vote] (m-1)%3으로 변경해야 하나?
	return values[(m-1)%2]
}

func ConvertJobWallet(s string) JobWallet {
	switch s {
	case "main":
		return Main
	case "stake":
		return Stake

		// // [Vote] 임시 condition
	// case "vote":
	// return Vote

	default:
		return Main
	}
}

func ValidateJobWallet(base JobWallet, target JobWallet) error {
	if base == 0 || base >= end {
		return ErrInvalidJobWallet
	}

	if target == 0 || target >= end {
		return ErrInvalidJobWallet
	}

	if base == Stake && target == Stake {
		return ErrInvalidJobWallet
	}

	// [Vote] 임시 condition
	// if base == Stake && target == Vote {
	// return ErrInvalidJobWallet
	// }

	// if base == Vote && target == Stake || base == Vote && target == Vote {
	// return ErrInvalidJobWallet
	// }

	return nil
}

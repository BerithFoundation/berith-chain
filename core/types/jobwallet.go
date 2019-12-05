/*
[BERITH]
Tx 타입을 지정하기 위한 열거형
*/
package types

import "errors"

type JobWallet uint8

const (
	Main = 1 + iota
	Stake
	end
)

var (
	values = [...]string{
		"main",
		"stake",
	}

	ErrInvalidJobWallet = errors.New("invalid wallet type")
)

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

	return nil
}

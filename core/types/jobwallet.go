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
	EthTx

	end
)

var (
	values = [...]string{
		"main",
		"stake",
		"ethtx",
	}

	ErrInvalidJobWallet = errors.New("invalid wallet type")
	ErrStakeToStake     = errors.New("cannot send balance stake to stake")
	ErrToEthTx          = errors.New("cannot send balance main/stake to ethtx")
	ErrFromEthTx        = errors.New("cannot send balance ethtx to main/stake")
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

	case "ethtx":
		return EthTx

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
		return ErrStakeToStake
	}

	if (base == Main || base == Stake) && target == EthTx {
		return ErrToEthTx
	}

	if base == EthTx && (target == Main || target == Stake) {
		return ErrFromEthTx
	}

	return nil
}

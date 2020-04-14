package staking

import (
	"fmt"
	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/params"
	"math/big"
	"testing"
)

/*
[BERITH]
선출 포인트 계산 테스트
 */
func TestCalcPoint(t *testing.T) {
	add_stake := big.NewInt(1000000)
	prev_stake := big.NewInt(10000000)
	new_block := big.NewInt(7200021)
	stake_block := big.NewInt(20)
	perioid := uint64(10)
	result := CalcPointUint(prev_stake, add_stake, new_block, stake_block, perioid)

	fmt.Println(result)
}

/*
[BERITH]
선출 포인트 계산 테스트
 */
func TestCalcPointBigint(t *testing.T) {
	type testData struct {
		add_stake *big.Int
		prev_stake *big.Int
		new_block *big.Int
		stake_block *big.Int
		perioid uint64
		isBIP4 bool
		want *big.Int
	}

	limitStakeBalanceInBer := new(big.Int).Div(params.MainnetChainConfig.Bsrr.LimitStakeBalance, big.NewInt(1e+18))

	tests := []testData {
		testData{add_stake: big.NewInt(1000000), prev_stake: big.NewInt(10000000), new_block: big.NewInt(7200021), stake_block: big.NewInt(20), perioid: uint64(10), isBIP4: params.MainnetChainConfig.IsBIP4(big.NewInt(100)), want: big.NewInt(20090909)},
		testData{add_stake: big.NewInt(1000000), prev_stake: big.NewInt(50000000), new_block: big.NewInt(7200021), stake_block: big.NewInt(20), perioid: uint64(10), isBIP4: params.MainnetChainConfig.IsBIP4(big.NewInt(100)), want: big.NewInt(100019607)},
		testData{add_stake: big.NewInt(1000000), prev_stake: big.NewInt(50000000), new_block: big.NewInt(7200021), stake_block: big.NewInt(20), perioid: uint64(10), isBIP4: params.MainnetChainConfig.IsBIP4(big.NewInt(3000000)), want: big.NewInt(99019607)},
	}

	for _, test := range tests {
		result := CalcPointBigint(test.prev_stake, test.add_stake, test.new_block, test.stake_block, limitStakeBalanceInBer, test.perioid, test.isBIP4)
		if result.Cmp(test.want) != 0 {
			t.Errorf("Expected %v but %v", test.want, result)
		}
	}
}

func TestCheckMaxStakeBalance(t *testing.T) {
	type testData struct {
		point *big.Int
		want *big.Int
	}

	limitStakeBalanceInBer := new(big.Int).Div(common.StringToBig(params.LimitStakeBalance), big.NewInt(1e+18))

	testDatas := []testData {
		testData{point: new(big.Int).Add(limitStakeBalanceInBer, big.NewInt(-1)), want: new(big.Int).Add(limitStakeBalanceInBer, big.NewInt(-1))},
		testData{point: limitStakeBalanceInBer, want: limitStakeBalanceInBer},
		testData{point: new(big.Int).Add(limitStakeBalanceInBer, big.NewInt(1)), want: limitStakeBalanceInBer},
		testData{point: new(big.Int).Add(limitStakeBalanceInBer, big.NewInt(100)), want: limitStakeBalanceInBer},
	}

	for _, testData := range testDatas {
		result := checkMaxStakeBalance(testData.point, limitStakeBalanceInBer)
		if result.Cmp(testData.want) != 0 {
			t.Errorf("expected %v but %v", testData.want, result)
		}
	}
}
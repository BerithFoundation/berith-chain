package bsrr

import (
	"testing"
	"time"

	"github.com/BerithFoundation/berith-chain/berith/selection"
	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/params"
)

func TestGetMaxMiningCandidates(t *testing.T) {
	var c = &BSRR{
		config: &params.BSRRConfig{
			Period:       10,
			Epoch:        360,
			Rewards:      common.StringToBig("20000"),
			StakeMinimum: common.StringToBig("100000000000000000000000"),
			SlashRound:   1000,
			ForkFactor:   0.3,
		},
	}
	tests := []struct {
		holders  int
		expected int
	}{
		{0, 0},                        // no holders
		{1, 1},                        // only one holders
		{10, 3},                       // equals to 0 point
		{8, 2},                        // less than 0.5 point
		{9, 3},                        // greater than or equals 0.5 point
		{35000, selection.MAX_MINERS}, // greater than staking.MAX_MINERS
	}

	for i, test := range tests {
		result := c.getMaxMiningCandidates(test.holders)
		if result != test.expected {
			t.Errorf("test #%d: expected : %d but %d", i, test.expected, result)
		}
	}
}

func TestGetDelay(t *testing.T) {
	var c = &BSRR{
		config: &params.BSRRConfig{
			Period:       0,
			Epoch:        360,
			Rewards:      common.StringToBig("20000"),
			StakeMinimum: common.StringToBig("100000000000000000000000"),
			SlashRound:   1000,
			ForkFactor:   1.0,
		},
		rankGroup: &common.ArithmeticGroup{CommonDiff: commonDiff},
	}

	tests := []struct {
		rank  int
		delay time.Duration
	}{
		{-1, time.Duration(0)},
		{0, time.Duration(0)},
		{1, time.Duration(0)},

		{2, 1 * groupDelay},
		{3, 1*groupDelay + 1*termDelay},
		{4, 1*groupDelay + 2*termDelay},

		{5, 2 * groupDelay},
		{6, 2*groupDelay + 1*termDelay},
		{7, 2*groupDelay + 2*termDelay},
	}

	for i, tt := range tests {
		result, _ := c.getDelay(tt.rank)
		if result != tt.delay {
			t.Errorf("test #%d: rank : %d expected : %d but %d", i, tt.rank, tt.delay, result)
		}
	}
}

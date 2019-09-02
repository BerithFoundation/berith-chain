package staking

import (
	"github.com/BerithFoundation/berith-chain/core/types"
	"github.com/BerithFoundation/berith-chain/params"
	"math/big"
)

func CalcPoint(pStake, addStake, now_block, stake_block *big.Int, config *params.ChainConfig, header *types.Header) uint64{

	b := float64(now_block.Uint64()) //블록넘버
	p := float64(pStake.Uint64()) //이전스테이킹
	n := float64(addStake.Uint64()) //추가스테이킹
	s := float64(stake_block.Uint64()) //이전 스테이킹 블록넘버

	ratio := (b * 100)  / (7200000 + s) //100은 소수점 처리

	if ratio > 100 {
		ratio = 100
	}
	adv := p * ( (p / (p + n)) * ratio) / 100
	result := p + adv + n

	return uint64(result)
}
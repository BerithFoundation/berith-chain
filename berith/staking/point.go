/**
[BERITH]
- Selection Point 를 계산하기 위함
- 10초 기준의 공식
- 10초주기로 블록 생성시 1년에 3600000블록을 생성함
**/

package staking

import (
	"math/big"
)

const (
	BLOCK_YEAR = 3600000 //10초 기준의 1년 블록 생성 수
)


/*
now_block : 블록넘버
pStake : 이전스테이킹 수량
addStake : 추가스테이킹 수량
stake_block : 이전 스테이킹 블록넘버
epoch : 블록당 생성시간
 */
func CalcPoint(pStake, addStake, now_block, stake_block *big.Int, period uint64) uint64{

	b := float64(now_block.Uint64()) //블록넘버
	p := float64(pStake.Uint64()) //이전스테이킹
	n := float64(addStake.Uint64()) //추가스테이킹
	s := float64(stake_block.Uint64()) //이전 스테이킹 블록넘버

	d := float64(period) / 10 //공식이 10초 단위 이기때문에 맞추기 위함 (perioid 를 제네시스로 변경하면 자동으로 변경되기 위함)

	bb := BLOCK_YEAR / d //기준 블록

	ratio := (b * 100)  / (bb + s) //100은 소수점 처리

	if ratio > 100 {
		ratio = 100
	}
	adv := p * ( (p / (p + n)) * ratio) / 100
	result := p + adv + n

	return uint64(result)
}
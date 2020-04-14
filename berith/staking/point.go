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
func CalcPointUint(pStake, addStake, now_block, stake_block *big.Int, period uint64) uint64{

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

func CalcPointBigint(pStake, addStake, now_block, stake_block, limitStakeBalanceInBer *big.Int, period uint64, isBIP4 bool) *big.Int {
	p := pStake //이전스테이킹
	n := addStake //추가스테이킹
	s := stake_block //이전 스테이킹 블록넘버
	b := now_block //블록넘버

	d := float64(period) / 10 //공식이 10초 단위 이기때문에 맞추기 위함 (perioid 를 제네시스로 변경하면 자동으로 변경되기 위함)

	bb := float64(BLOCK_YEAR / d) //기준 블록

	//ratio := (b * 100)  / (bb + s) //100은 소수점 처리
	ratio := new(big.Float).Mul(new(big.Float).SetInt(b), big.NewFloat(100))
	ratio.Quo(ratio, new(big.Float).Add(big.NewFloat(bb), new(big.Float).SetInt(s)))

	/*
	if ratio > 100 {
		ratio = 100
	}
	*/
	if ratio.Cmp(big.NewFloat(100)) == 1 {
		ratio = big.NewFloat(100)
	}


	temp1 := new(big.Float).Quo(new(big.Float).SetInt(p), new(big.Float).Add(new(big.Float).SetInt(p), new(big.Float).SetInt(n)))
	temp2 := new(big.Float).Mul(new(big.Float).SetInt(p), temp1)
	temp3 := new(big.Float).Mul(temp2, ratio)
	adv := new(big.Int)
	new(big.Float).Quo(temp3, big.NewFloat(100)).Int(adv)

	/*
		[Berith]
		Stake Balance를 한도 이상가지고 있는 경우 한도 + Advantage 만큼만 Selection Point를 갖도록 처리
	*/
	r1 := new(big.Int).Add(p, n)
	if isBIP4 {
		r1 = checkMaxStakeBalance(r1, limitStakeBalanceInBer)
	}

	//result := p + adv + n
	r2 := new(big.Int).Add(r1, adv)

	return r2
}

/*
	[Berith]
	Stake Balance 한도 초과 보유자에 대한 선출 포인트 Advantage를 새로 계산하는 함수
*/
func CalcAdvForExceededPoint(nowBlockNumber, stakeBlockNumber *big.Int, period uint64, limitStakeBalanceInBer *big.Float) *big.Int {
	d := float64(period) / 10 //공식이 10초 단위 이기때문에 맞추기 위함 (perioid 를 제네시스로 변경하면 자동으로 변경되기 위함)

	bb := float64(BLOCK_YEAR / d) //기준 블록

	//ratio := (b * 100)  / (bb + s) //100은 소수점 처리
	ratio := new(big.Float).Mul(new(big.Float).SetInt(nowBlockNumber), big.NewFloat(100))
	ratio.Quo(ratio, new(big.Float).Add(big.NewFloat(bb), new(big.Float).SetInt(stakeBlockNumber)))

	/*
		if ratio > 100 {
			ratio = 100
		}
	*/
	if ratio.Cmp(big.NewFloat(100)) == 1 {
		ratio = big.NewFloat(100)
	}

	temp1 := new(big.Float).Quo(limitStakeBalanceInBer, new(big.Float).Add(limitStakeBalanceInBer, big.NewFloat(0)))
	temp2 := new(big.Float).Mul(limitStakeBalanceInBer, temp1)
	temp3 := new(big.Float).Mul(temp2, ratio)
	adv := new(big.Int)
	new(big.Float).Quo(temp3, big.NewFloat(100)).Int(adv)

	return adv
}

/*
	[Berith]
	BIP4에서 Stake Balance의 한도를 설정한 것과 관련하여 Selection Point 또한 한도를 갖도록 설정
*/
func checkMaxStakeBalance(stake, limitStakeBalanceInBer *big.Int ) *big.Int {
	if stake.Cmp(limitStakeBalanceInBer) == 1 {
		return limitStakeBalanceInBer
	}
	return stake
}
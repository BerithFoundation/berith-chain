## 베리드의 합의엔진
베리드의 노드들은 PoS 합의 알고리즘을 통해 블록을 합의한다. 블록을 생성하는 권한을 토큰을 예치한 유저에게 부여하여, 돈을 예치한 유저들이 토큰 가격의 변동의 위험성을 감수하면서 까지 공격을 시도하지 않을 것이라는 배경에서 만들어진 알고리즘이다.


### 토큰을 스테이킹한 유저 관리

블록의 정당성을 증명하기 위해 모든 노드는 토큰을 스테이킹한 유저의 목록을 관리할 필요가 있다.

#### 스테이크 트랜잭션 추적

스테이크는 트랜잭션을 통해 이루어진다. 베리드는 스테이킹 기능을 구현하기 위해 계정의 잔액을 두가지로 구분한다. 

```
type Account struct {
    Nonce          uint64
    Balance        *big.Int
    Root           common.Hash // merkle root of the storage trie
    CodeHash       []byte
    StakeBalance   *big.Int //brt staking balance
    StakeUpdated   *big.Int //Block number when the stake balance was updated
    Point          *big.Int //selection Point
    BehindBalance  []Behind //behind balance
    Penalty        uint64
    PenlatyUpdated *big.Int //Block Number when the penalty was updated
}
```

위의 코드는 머클패트리샤 트리에 저장되는 계정 정보 구조체이다. 베리드의 계정이 일반적인 잔액인 ```MainBalance``` 를 저장하는 ```Balance``` 필드와, 스테이킹된 토큰의 잔액을 저장하는 ```StakeBalance``` 필드, 두 가지로 구분되는 것을 나타낸다.

그래서 베리드의 트랜잭션은 계정의 어떤 잔액을 사용할 것인지를 나타내 주어야한다. 이를 위해서 계정의 지갑타입을 나타내는 정수를 사용한다.
```
type JobWallet uint8

const (
    Main = 1 + iota
    Stake
    end
)
```
위의 코드는 ```JobWallet``` 타입의 선언부이다.  ```JobWallet``` 타입은 정수타입으로 1인 경우 ```MainBalance```, 2인 경우 ```StakeBalance``` 임을 표시한다.

베리드는 ```JobWallet``` 을 이용하여 트랜잭션을 통해 토큰을 보내고 받는 유저가 사용할 지갑의 종류를 구분한다. 

```
type txdata struct {
    AccountNonce uint64          `json:"nonce"    gencodec:"required"`
    Price        *big.Int        `json:"gasPrice" gencodec:"required"`
    GasLimit     uint64          `json:"gas"      gencodec:"required"`
    Recipient    *common.Address `json:"to"       rlp:"nil"` // nil means contract creation
    Amount       *big.Int        `json:"value"    gencodec:"required"`
    Payload      []byte          `json:"input"    gencodec:"required"`
    Base         JobWallet       `json:"base"  gencodec:"required"`   //[Berith] 작업 주체  ex) 스테이킹시 : Main
    Target       JobWallet       `json:"target"  gencodec:"required"` //[Berith] 작업타겟 ex) 스테이킹시 : Stake

    // Signature values
    V *big.Int `json:"v" gencodec:"required"`
    R *big.Int `json:"r" gencodec:"required"`
    S *big.Int `json:"s" gencodec:"required"`

    // This is only used when marshaling to JSON.
    Hash *common.Hash `json:"hash" rlp:"-"`
}
```

위의 코드는 트랜잭션 구조체이다. 이더리움의 트랜잭션과 다르게 ```Base``` 와 ```Target``` 이라는 필드가 추가된 것을 확인할 수 있다. ```Base``` 는 ```From``` 계정이 트랜잭션을 통해 보낼 토큰이 들어있는 잔액의 종류를 표시한다. ```Target``` 은 ```To``` 계정이 트랜잭션을 통해 전송된 토큰이 들어갈 잔액의 종류를 표시한다. 

이를 통해, ```Base``` 와 ```Target``` 이 ```StakeBalane``` 를 지정하는 경우 해당 트랜잭션이 일반적인 트랜잭션이 아니고 스테이크와 관련된 트랜잭션임을 확인할 수 있다.  베리드에서 스테이크와 관련된 트랜잭션은 From 과 To 가 같아야하는 조건이 있으며, 두가지 형태로 나누어진다.

```Stake``` : ```Base``` 가 ```MainBalance``` 이고 ```Target``` 이 ```StakeBalance``` 인 트랜잭션이다. 유저가 ```MainBalance``` 의 토큰을 스테이크할 때 사용한다. 

```Unstake``` : ```Base``` 가 ```StakeBalance``` 이고 ```Target이``` ```MainBalance``` 인 트랜잭션이다. 유저가 스테이크된 토큰을 다시 ```MainBalance``` 로 옮기고 싶을 때 사용한다. ```Unstake``` 의 ```Value``` 의 지정과 상관없이, 모든 ```StakeBalance``` 의 잔액이 ```Unstake``` 된다. 

 베리드는 블록 바디를 유효성을 검사하는 함수안에서 블록에 포함될 트랜잭션의 ```Base``` 와 ```Target``` 필드를 확인하여 스테이크 트랜잭션을 추적한다.

#### 블록별 토큰을 스테이크한 유저의 목록 관리

베리드는 모든 블록에 대해 해당 블록 시점에 토큰을 스테이크하고 있는 계정의 목록을 저장한다. 이를 ```Stakers``` 라고 한다. 위에서 언급 한 블록 바디의 유효성을 검사하는 함수에서 추적한 스테이크 트랜잭션 정보를 이용하여 이전 블록의 ```Stakers``` 에서 새롭게 스테이크한 계정을 추가하고, 스테이크를 취소한 계정을 삭제하여 새로운  ```Stakers``` 를 만들어 저장한다.

```
if chain.Config().IsBIP1(number) {
            if msg.Base() == types.Main && msg.Target() == types.Stake {
                stkChanged[msg.From()] = true
            } else if msg.Base() == types.Stake && msg.Target() == types.Main {
                stkChanged[msg.From()] = false
            } else {
                continue
            }
        } else {
            if msg.Base() == types.Main && msg.Target() == types.Stake {
                stkChanged[msg.From()] = true
            } else {
                continue
            }
        }
}

for addr, isAdd := range stkChanged {
	
		...
			
		if isAdd {
			stks.Put(addr)
		} else {
			stks.Remove(addr)
		}

	}

```
위의 코드는 새로운 ```Stakers``` 를 만드는 코드의 일부로 기존의 ```Stakers``` 에서 계정을 추가하거나, 제거하는 것을 확인 할 수 있다.

만들어진 ```Stakers``` 는 현재 블록의 넘버보다 하나 작은 블록까지만 레벨DB에 저장한다. 레벨 DB에는 블록의 해쉬값을 키로, ```Stakers``` 를 인코딩한 값을 밸류로 저장한다. 현재 블록을 저장하지 않는 이유는 두가지가 있다.

현재 블록넘버에 해당하는 블록은 변경되는 경우가 많기 때문

블록 바디의 유효성을 검사하는 함수에 자기가 생성한 블록을 첫번째로 검사 할 때 블록에 서명이 포함되어있지 않기 때문에 블록의 해쉬값이 실제값과 다르기 때문

따라서, 블록에 해당하는 ```Stakers``` 가 저장되는 시점은 현재 블록을 부모로 갖는 블록을 받았을 때이다. 

### 블록 생성자 추첨

베리드에서는 토큰을 스테이크한 모든 유저가 블록을 생성 및 전파 할 수 있다. 하지만 추첨 결과에 따라 블록의 우선순위와 전파 딜레이가 결정된다.

![select](./selector&#32;table.png)

위의 그림은 블록 생성자를 추첨하는 과정을 나타낸다. 이에 대해 아래에서 설명한다.

#### 블록 생성자 추첨과 Stakers, Epoch

블록 생성자 추첨은 특정 블록의 ```Stakers``` 를 기준으로 진행된다. 특정 블록은 현재 블록의 부모 블록을 ```Epoch``` 만큼 거슬러 올라간 블록이 된다. ```Epoch는``` 정수로 노드의 ```ChainConfig``` 에 지정되어있다. 
```
MainnetChainConfig = &ChainConfig{

        ...

        Bsrr: &BSRRConfig{
            Period:       5,
            Epoch:        360,
            Rewards:      common.StringToBig("360"),
            StakeMinimum: common.StringToBig("100000000000000000000000"),
            SlashRound:   0,
            ForkFactor:   1.0,
        },
    }
```
위의 코드는 설정파일 없이 노드를 실행 시켰을 때 기본적으로 적용되는 ```ChainConfig``` 를 선언한 코드이다. 베리드의 ```Epoch``` 이 기본적으로 360으로 설정된 것을 확인할 수 있다.

![epoch](./Round.png)

위의 그림은 `Epoch` 만큼 이전의 블록은 무엇인가를 나타낸다. 그림에서 확인할 수 있듯이 `Block 361` 에 대한 추첨은 `Block 1` 을 기준으로 계산되는 것임을 알 수 있다.

블록 생성자 추첨에서 현재 블록이 아닌 ```Epoch``` 만큼 이전의 블록을 사용하는 이유는 두가지가 있다.

1. 해당 블록의 트랜잭션에 따라 추첨결과가 달라지는 것을 방지
2. 스테이크를 한 유저가 바로 블록을 생성할 수 없게 하기 위해

#### 블록 생성자 추첨과 난수

블록 생성자 추첨은 지정된 ```Stakers``` 의 내용을 랜덤한 순서로 정렬하는 것과 같다. 하지만 이를 블록체인에서 적용하기 위해서는 두가지 조건을 만족해야 한다.

1. 랜덤 값이 모든 노드에서 동일해야 한다는 조건
2. 특정한 유저가 랜덤값을 결정할 수 없어야 한다는 조건

베리드는 이 두가지 조건을 만족하기 위해 블록의 넘버를 시드로 하여 난수를 생성한다. 블록의 넘버는 같은 블록을 받은 노드라면 모두 같은 값을 가지고 있고, 블록 내용이 어떤 것이라도 변하지 않는 값이기 때문에 두가지 조건을 만족한다.

#### 블록 생성자 추첨과 Point

```Point``` 는 추첨에서 사용되는 값으로 스테이크한 토큰의 갯수, 토큰을 스테이크한 블록을 이용하여 계산된다. 어떤 계정이 추첨에서 뽑히는 확률은 (계정의 Point / 전체 Point) * 100 % 이다. 
```
func CalcPointBigint(pStake, addStake, now_block, stake_block *big.Int, period uint64) *big.Int {
    b := now_block //블록넘버
    p := pStake //이전스테이킹
    n := addStake //추가스테이킹
    s := stake_block //이전 스테이킹 블록넘버

    d := float64(period) / 10 //공식이 10초 단위 이기때문에 맞추기 위함 (perioid 를 제네시스로 변경하면 자동으로 변경되기 위함)

    bb := int64(BLOCK_YEAR / d) //기준 블록

    //ratio := (b * 100)  / (bb + s) //100은 소수점 처리
    ratio := new(big.Int).Mul(b, big.NewInt(100))
    ratio.Div(ratio, new(big.Int).Add(big.NewInt(bb), s))

    /*
    if ratio > 100 {
        ratio = 100
    }
     */
    if ratio.Cmp(big.NewInt(100)) == 1 {
        ratio = big.NewInt(100)
    }

    //adv := p * ((p / (p + n)) * ratio) / 100
    temp1 := new(big.Int).Div(p, new(big.Int).Add(p, n))
    temp2 := new(big.Int).Mul(p, temp1)
    temp3 := new(big.Int).Mul(temp2, ratio)
    adv := new(big.Int).Div(temp3, big.NewInt(100))

    //result := p + adv + n
    r1 := new(big.Int).Add(p, adv)
    r2 := new(big.Int).Add(r1, n)

    return r2
}
```
위의 코드는 ```Point``` 를 계산하는 함수이다.

#### 블록 생성자 추첨 과정

추첨을 위해 현재 블록의 ```Epoch``` 만큼 이전의 블록을 얻어낸다. 
```
// [BERITH] getStakeTargetBlock 주어진 parent header에 대하여 miner를 결정 할 target block을 반환한다.
// 1) [0, epoch-1] : target == 블록 넘버 0(즉, genesis block) 인 블록
// 2) [epoch, 2epoch] : target == 블록 넘버 epoch 인 블록
// 3) [2epoch +1, ~) : target == 블록 넘버 - epoch 인 블록
func (c *BSRR) getStakeTargetBlock(chain consensus.ChainReader, parent *types.Header) (*types.Header, bool) {
    if parent == nil {
        return &types.Header{}, false
    }

    var targetNumber uint64
    blockNumber := parent.Number.Uint64()
    d := blockNumber / c.config.Epoch

    if d > 1 {
        return c.getAncestor(chain, int64(c.config.Epoch), parent)
    }

    switch d {
    case 0:
        targetNumber = 0
    case 1:
        targetNumber = c.config.Epoch
    }

    target := chain.GetHeaderByNumber(targetNumber)
    if target != nil {
        return target, chain.HasBlockAndState(target.Hash(), targetNumber)
    }
    return target, false
}
```
위의 코드는 ```Epoch``` 만큼 이전의 블록을 얻어내는 함수의 내용이다. 

얻어낸 블록의 해쉬값을 이용하여 레벨DB에서 ```Stakers``` 를 얻는다. 이후 얻어낸 ```Stakers``` 를 계정의 해쉬값을 기준으로 정렬한다. 그리고 정렬한 계정의 목록을 순회하여 추첨 테이블인 ```Candidates``` 를 만들어낸다. 추첨 테이블은 계정의 ```Point``` 와 이전 인덱스의 계정들의 ```Point``` 합으로 구성된다.

만약 계정 A, B, C 가 각 각 10, 20 ,30의 ```Point``` 를 가지고 있다면 추첨 테이블은 다음과 같이 구성된다.

| field \ Index | 0 | 1 | 2 |
|:---|:---|:---|:---|
|account|A|B|C|
|point|10|20|30|
|val|10|30|60|

```
func SelectBlockCreator(config *params.ChainConfig, number uint64, hash common.Hash, stks staking.Stakers, state *state.StateDB) VoteResults {
    result := make(VoteResults)

    list := sortableList(stks.AsList())
    if len(list) == 0 {
        return result
    }

    sort.Sort(list)

    cddts := NewCandidates()

    for _, stk := range list {
        point := state.GetPoint(stk).Uint64()
        cddts.Add(Candidate{
            point:   point,
            address: stk,
        })
    }
    if config.IsBIP3(big.NewInt(int64(number))) {
        result = cddts.selectBIP3BlockCreator(config, number)
    } else {
        result = cddts.selectBlockCreator(config, number)
    }

    return result

}
```
위의 코드는 ```Stakers``` 를 정렬해서 추첨 테이블인 ```Candidates``` 를 만드는 함수의 내용이다. 

추첨은 계정을 키로 추첨 결과인 ```VoteResult``` 구조체를 값으로 갖는 맵을 만들어 내는 과정이다.
```
type VoteResult struct {
    Score *big.Int `json:"score"`
    Rank  int      `json:"rank"`
}
```
위의 코드는 ```VoteResult``` 구조체의 선언부이다. 점수인 ```Score```, 순위인 ```Rank``` 를 갖는 구조체임을 확인할 수 있다. ```Rank``` 는 추첨된 순서, 점수 ```Score``` 는 아래의 식에 의해 결정된다.

```len``` : 추첨 테이블의 길이

```Score``` : ```5000000 - ((5000000 - 10000)  / len * (Rank-1))```

추첨은 ```Candaidates``` 의 길이가 0 보다 큰 동안 아래의 과정을 반복하며 결과값으로 계정을 키로 ```VoteResult``` 를 값을 갖는 맵 ```VoteResults``` 를 구성한다.

* 난수 ```x``` 를 생성한다.

* ```rank``` : 1

* ```len``` : ```Candidates``` 의 길이

* ```diff``` : 5000000

* ```diff_r``` : ```(diff - 10000) / len```

* 이분 탐색으로 원하는 인덱스 ```n``` 을 찾을 때 까지 반복하며 아래의 과정을 반복한다.

    + 생성된 난수가 ```Candidates``` 에서 인덱스 ```n-1```, ```n``` 에 해당하는 ```val``` 를 구한다. 각각 ```val(n-1)```, ```val(n)```

    + 난수 x 가 ```val(n-1) <= x <= val(n)```  를 만족하는 지 확인한다.

    + 만족하는 경우 아래의 과정을 실행한다. 아닌 경우 생략하고 반복한다.

        - ```n``` 에 해당하는 ```account``` 를 키로 ```VoteResult = { Score : diff, Rank : rank }``` 를 값으로 갖는 구조체를 맵에 저장한다.

        - ``diff -= diff_r, rank++``

        - ```n``` 에 해당하는 ```accounts``` 를 ```Candidates``` 에서 제거하고 이분탐색을 종료한다.
```
func (cs *Candidates) selectBIP3BlockCreator(config *params.ChainConfig, number uint64) VoteResults {
    result := make(VoteResults)

    DIF := DIF_MAX
    DIF_R := (DIF_MAX - DIF_MIN) / int64(len(cs.selections))
    rank := 1
    rand.Seed(cs.GetSeed(config, number))

    for len(cs.selections) > 0 {

        target := uint64(rand.Int63n(int64(cs.total)))

        var chosen int
        start := 0
        end := len(cs.selections) - 1

        for {
            mid := (start + end) / 2
            a := uint64(0)
            if mid > 0 {
                a = cs.selections[mid-1].val
            }
            b := cs.selections[mid].val

            if target >= a && target <= b {
                chosen = mid
                cddt := cs.selections[mid]
                result[cddt.address] = VoteResult{
                    Rank:  rank,
                    Score: big.NewInt(DIF),
                }
                DIF -= DIF_R
                rank++
                break
            }

            if target < a {
                end = mid - 1
            }
            if target > b {
                start = mid + 1
            }
        }

        out := cs.selections[chosen]
        for i := chosen; i+1 < len(cs.selections); i++ {
            newCddt := cs.selections[i+1]
            newCddt.val -= out.point
            cs.selections[i] = newCddt
        }

        cs.selections = cs.selections[:len(cs.selections)-1]
        cs.total -= out.point
    }

    //fmt.Println(DIF)
    return result
}
```
위의 코드는 블록 생성자를 추첨하는 함수의 내용이다.

#### 블록 생성자 추첨과 블록 생성 우선순위

블록 생성자는 자신의 계정에 해당하는 ```Score```, ```Rank```를 ```VoteResults``` 에서 가져올 수 있다. 이 값들은 블록의 우선순위를 정하는 것에 관련이 있다. Score 는 블록의 ```Difiiculty``` 로 사용된다. ```Difficulty```는 노드가 블록체인을 선택 할 때 사용된다. 노드는 ```Difiiculty``` 의 합이 가장 높은 블록체인을 선택한다. ```Rank``` 는 블록의 전파 딜레와 관련있다. 베리드의 노드들은 부모 블록의 타임스탬프 보다 ```ChainConfig``` 의 ```Period``` 초 만큼 지난 시점에 블록을 전파하는데, ```Rank``` 에 따라 추가적인 딜레이가 부여된다.

따라서, 추첨에서 먼저 뽑힌 계정이 생성한 블록일 수록 우선순위가 높다고 할 수 있다.

### 블록생성보상

베리드는 블록을 생성한 계정에게 보상을 지급한다. 블록생성보상은 100년 간, 50억개의 코인을 보상하도록 계산된다. 

![reward](./berith_reward.png)
위의 그림은 블록생성보상이 기간에 따라 어떻게 변하는지 나타내는 그래프이다. 베리드에서 블록생성보상은 시간이 지남에 따라 점점 줄어드는 것을 학인할 수 있다.

```
func getReward(config *params.ChainConfig, header *types.Header) *big.Int {
	number := header.Number.Uint64()
	// 특정 블록 이후로 보상을 지급
	if number < config.Bsrr.Rewards.Uint64() {
		return big.NewInt(0)
	}

	//공식이 10초 단위 이기때문
	d := float64(config.Bsrr.Period) / 10
	n := float64(number) * d

	var z float64 = 0
	if n <= 3150000 {
		z = 5
	}

	re := (26 - math.Round(n/(7370000))*0.5 + z) * d
	if re <= 0 {
		re = 0

		return big.NewInt(0)
	} else {
		temp := re * 1e+10
		return new(big.Int).Mul(big.NewInt(int64(temp)), big.NewInt(1e+8))
	}
}
```
위의 코드는 블록생성보상의 값을 계산하는 함수의 내용이다.

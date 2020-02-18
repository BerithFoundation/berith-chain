## 베리드와 하드포크

### 하드포크란?

패치된 버전에서 생성된 블록이 이전 버전에서 호환되지 않을 경우를 말한다. 사용자는 새로운 버전을 다운 받지 않으면 이후 부터 만들어지는 블록을 블록체인에 등록할 수 없다.

#### 하드포크 패치

하드포크를 처리하기 위해서는 패치된 버전이 이전 버전의 알고리즘과 수정된 알고리즘을 동시에 가지고 있어야 한다 왜냐하면 여태껏 쌓여온 블록들을 검증 할 때는 이전의 알고리즘을 사용하고, 하드포크 시점 이후의 블록에는 수정된 알고리즘을 적용해야 하기 때문이다.

#### 베리드 하드포크

베리드는 하드포크를 위해 ```ChainConfig``` 를 사용한다. ```ChainConfig``` 에 하드포크 할 시점의 블록 번호를 지정하여 해당 시점부터는 새로운 알고리즘을 적용하는 방식으로 하드포크를 만들 수 있다.
```
MainnetChainConfig = &ChainConfig{
        ChainID:             big.NewInt(106),
        HomesteadBlock:      big.NewInt(0),
        DAOForkBlock:        nil,
        DAOForkSupport:      true,
        EIP150Block:         big.NewInt(0),
        EIP155Block:         big.NewInt(0),
        EIP158Block:         big.NewInt(0),
        ByzantiumBlock:      big.NewInt(0),
        ConstantinopleBlock: big.NewInt(0),
        BIP1Block:           big.NewInt(508000),
        BIP2Block:           big.NewInt(545000),
        BIP3Block:           big.NewInt(1168000),

        ...

    }
```
위의 코드는 설정 파일 없이 베리드를 실행하는 경우 사용되는 기본 ```ChainConfig``` 를 선언하는 부분이다. 코드에서 보이는 ```BIP1Block```, ```BIP2Block```, ```BIP3Block``` 은 모두 하드포크를 했던 블록 넘버를 나타내고 있다.

#### chaindata 호환

하드포크 시점을 ```ChainConfig``` 에 등록 하는 것 만으로도 새롭게 만들어지는 ```chaindata``` 에 하드포크 설정을 적용할 수 있지만, 기존의 ```chaindata``` 를 가지고 실행되는 경우에는 하드포크를 적용할 수 없다. 왜냐하면 ```ChainConfig``` 에 대한 정보는 ```genesisBlock``` 을 생성할 때 같이 저장되고 새롭게 실행하는 경우 이전의 ```ChainConfig``` 를 사용하기 때문이다. 베리드는 이를 해결하기 위해서 ```chaindata``` 가 사용하고 있는 ```genesisHash``` 가 메인넷에서 사용하는 기본 값을 가지고 있고, 코드에 등록된 기본 ```ChainConfig``` 와 ```chaindata``` 에 저장되어 있는 ```ChainConfig``` 가 서로 호환이 안되는 경우에 ```chaindata``` 에 저장된 ```ChainConfig``` 를 지우고 기본 ```ChainConfig``` 로 바꾸는 기능을 가지고 있다.
```
func (c *ChainConfig) checkCompatible(newcfg *ChainConfig, head *big.Int) *ConfigCompatError {
   
        ...

    if isForkIncompatible(c.BIP1Block, newcfg.BIP1Block, head) {
        return newCompatError("bip1 fork block", c.BIP1Block, newcfg.BIP1Block)
    }
    if isForkIncompatible(c.BIP2Block, newcfg.BIP2Block, head) {
        return newCompatError("bip2 fork block", c.BIP2Block, newcfg.BIP2Block)
    }
    if isForkIncompatible(c.BIP3Block, newcfg.BIP3Block, head) {
        return newCompatError("bip3 fork block", c.BIP3Block, newcfg.BIP3Block)
    }
    return nil
}
```
위의 코드는 두개의 ```ChainConfig``` 가 서로 호환되는 지 확인하는 함수의 내용 중 일부이다. 새로운 하드포크 시점을 추가한 경우, 이곳에 추가한 하드포크 시점이 없는 경우에 호환되지 않는다는 조건을 추가하는 것으로 기존 ```chaindata``` 의 ```ChainConfig``` 를 새로 등록한 ```ChainConfig``` 로 덮어 씌울 수 있다.

### 베리드 하드포크 이력

#### BIP1

스테이크를 취소해도 ```Stakers``` 에서 계정이 지워지지 않는 문제 수정
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
```
위의 코드는 트랜잭션이 타입을 확인하여 맵을 만들어내는 부분이다. 맵은 계정을 키로, 논리값을 값으로 갖는데 이후 코드에서 iterate 하여 값이 ```true``` 인 계정은 ```Stakers``` 에 추가하고, 값이 ```false``` 인 계정은 ```Stakers``` 에서 제거한다. BIP1 이 적용 된 블록부터는 스테이크 해제에 대한 트랜잭션에 다른 조건문이 쓰이는 것을 확인할 수 있다.

#### BIP2

블록 생성자 추첨을 위해 시드를 생성할 때, 해쉬 함수를 잘못 사용한 문제 수정
```
func (cs Candidates) GetSeed(config *params.ChainConfig, number uint64) int64 {

    bt := []byte{byte(number)}
    if config.IsBIP2(big.NewInt(0).SetUint64(number)) {
        bt = big.NewInt(0).SetUint64(number).Bytes()
    }
    hash := sha256.New()
    hash.Write(bt)
    md := hash.Sum(nil)
    h := common.BytesToHash(md)
    seed := h.Big().Int64()

    return seed
}
```
위의 코드는 블록 생성자의 추첨을 위해 생성할 난수의 시드값을 얻어내는 함수의 내용이다. 이전의 코드는 인자로 받은 ```number``` 의 첫번 째 바이트만을 해쉬하여 사용하기에 정확한 해쉬값을 기대할 수 없다. 그래서 BIP2 이후 부터는 ```number``` 값 전체를 해쉬하여 올바르게 해쉬값을 얻어낼 수 있도록 수정했다.

#### BIP3

블록 생성자 추첨 결과가 일정하지 않고 유리한 계정군이 생기는 문제 수정
```
func SelectBlockCreator(config *params.ChainConfig, number uint64, hash common.Hash, stks staking.Stakers, state *state.StateDB) VoteResults {
   
    ...

    if config.IsBIP3(big.NewInt(int64(number))) {
        result = cddts.selectBIP3BlockCreator(config, number)
    } else {
        result = cddts.selectBlockCreator(config, number)
    }

    return result

}
```
위의 코드는 블록 생성자 추첨 함수 중 일부이다. BIP3 이후의 블록인 경우 기존과는 다른 함수를 호출하는 것을 확인할 수 있다.

두개의 함수는 다른 알고리즘을 가지고 실행된다. 이를 정리하면

|Feature\Function|selectBlockCreator|selectBIP3BlockCreator|
|:---|:---|:---|
|처리 속도|빠름|다소 느림|
|분산|불평등함|평등함|
|처리 방법|배열 전체를 이진 탐색한 결과에 따라 배열을 나누고 나누어진 부분 배열을 다시 이진 탐색하여 추첨 결과를 얻어냄|배열 전체를 이진 탐색하여 추첨 결과를 얻어낸 뒤 추첨된 계정을 제외하는 것을 반복함|

처리속도가 추첨결과의 평등함 보다 우선시될 수 없기에 이를 수정하였다.


## 베리드와 이더리움
베리드는 이더리움의  golang 구현체인 go-ethereum 을 커스터마이징한 프로젝트이다.

 [https://github.com/ethereum/go-ethereum](https://github.com/ethereum/go-ethereum)

베리드는 패키지나 변수명 등에서 이더리움과 관련된 명칭을 베리드로 바꿨다. 
* ```ether```, ```ethereum``` → ```berith```
* ```eth``` → ```ber```
  
와 같은 변경을 소스 코드에서 찾아볼 수 있다.

### 이더리움 설정값

베리드는 이더리움을 커스터마이징한 프로젝트이지만 합의 알고리즘 특성상 이더리움에서 사용하는 설정을 베리드에서 사용할 수 없는 경우가 있다.

#### SyncMode

베리드는 블록헤더만으로 블록을 검증할 수 없기에 ```SyncMode``` 를 ```light``` 로 설정하여 노드를 실행할 수 없다. 베리드는 과거의 블록을 현재 블록의 검증에 사용할 수 있기에 ```SyncMode``` 를 ```fast``` 로 설정하여 노드를 실행할 수 없다. 따라서 베리드의 ```SyncMode``` 는 ```full``` 로 설정해야만 한다.

#### GCMode

베리드는 블록을 검증할 때, 과거 블록의 ```state``` 를 사용한다. 그렇기에 ```GCmode``` 를 ```full``` 로 설정하여 노드 종료시에 메모리에 있는 ```state``` 들이 지워진다면 블록을 검증에 문제가 발생할 수 있다. 따라서 베리드의 ```GCMode``` 는 ```archive``` 로 설정해야만 한다.

### 이더리움과 베리드의 네트워크 분리

이더리움은 다른 이더리움 노드를 찾기 위해 UDP 통신을 이용하는데 이를 ```discover``` 라 한다. 이더리움의 ```discover``` 는 요청의 종류를 요청의 첫 패킷의 값으로 구분한다. 만약 첫 패킷의 값이 요청의 종류와 매핑되어 있지 않다면 이 요청을 무시한다. 베리드는 이것을 이용하여 요청의 종류가 매핑되어 있는 상수를 수정하여 이더리움과 네트워크를 분리했다.

|network\type|ping|pong|findnode|neighbors|
|:---|:---|:---|:---|:---|
|Ethereum|1|2|3|4|
|Berith|11|12|13|14|

```
//Ethereum
const (
	// Packet type events.
	// These correspond to packet types in the UDP protocol.
	pingPacket = iota + 1
	pongPacket
	findnodePacket
	neighborsPacket
    
    ...
    	
)


//Berith
const (
    pingPacket = iota + 11 // zero is 'reserved'
    pongPacket
    findnodePacket
    neighborsPacket
)
```
위의 코드는 이더리움과 베리드의 ```discover``` 요청의 패킷 타입을 선언하는 부분이다. 베리드와 이더리움이 같은 요청에 대해 다른 값을 매핑하는 것을 확인할 수 있다.

### 베리드 RPC

베리드의 RPC 메소드는 Ethereum 과 마찬가지로 package_method 의 형식으로 이루어져 있다. 베리드의 RPC 메소드는 대부분 이더리움과 동일하지만 트랜잭션 인자의 형태가 변경되거나, 스테이크에 관한 메소드가 추가되는 것 같은 변경점이 있다. 자세한 RPC spec에 대해서는 [Berith RPC Spec](./rpc.md) 에서 확인할 수 있다.

### 주소체계

베리드의 주소는 “0x“ 로 시작하는 이더리움과 달리 “bx“로 시작하는 주소체계를 가지고있다. 이는 이더리움 계정 앞에 붙는 “0x“ 가 단순히 데이터가 16진수로 이루어져 있다는 것을 표현하는 Prefix 로, 뒤의 실제 계정을 이루는 값과는 무관하다는  것을 이용하여 표시만 바꿔준 것이다. 원격에서 보내진 json 요청을 go 에서 처리하기 위해 구조체 형태로 unmarshal 할 때, 주소 데이터의 유효성 검사 함수를 바꾼 것이다.
```
const (

      ...

	// AddressPrefix is prefix of the address
	AddressPrefix = "Bx"
	// AddressPrefixLength is the expected length of the address prefix
	AddressPrefixLength = len(AddressPrefix)
)


func CheckBerithPrefix(str string, start int) bool {
    if len(str) > start+2 {
        lower := strings.ToLower(AddressPrefix)
        upper := strings.ToUpper(AddressPrefix)
        if (str[start] == lower[0] || str[start] == upper[0]) && (str[start+1] == lower[1] || str[start+1] == upper[1]) {
            return true
        }
    }
    return false
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *Address) UnmarshalJSON(input []byte) error {
	if len(input) > 2 {
		inputStr := string(input)
		// TODO : must have addr prefix ?
		if !HasAddressPrefix(inputStr[1:]) {
			return fmt.Errorf("berith address without \"%s\" prefix", AddressPrefix)
		}
		input = []byte(inputStr[:1] + "0x" + inputStr[3:])
	}
	return hexutil.UnmarshalFixedJSON(addressT, input, a[:])
}
```
 위의 코드는 ```Address``` 타입의 데이터를 unmarshal 할 때 호출되는 함수이다. 이를 통해 “bx“ prefix 로 들어온 데이터를 받아 prefix를 “0x“ 로 바꿔서 기존의 유효성 검사를 통과하게 하는 것을 확인할 수 있다. 또 “0x“ 로 들어온 데이터는 “bx“ prefix를 검사할 때, 유효하지 않다고 판단되어 에러를 반환하는 것도 확인할 수 있다.


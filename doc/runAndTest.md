## 베리드의 실행과 테스트

베리드를 실행하고 테스트하는 방법에 대해 기술한다.

### 베리드의 빌드

베리드는 고 프로그램으로 고 바이너리를 통해 빌드되어 바이너리 실행파일이 생성된다.

#### make를 이용한 빌드

프로젝트 최상단에 make 명령을 위한 ```Makefile``` 이 정의 되어있다. 빌드는 make 명령에 커맨드를 추가 하는 것으로 특정 패키지를 따로 진행할 수 있다. 만약 노드 실행을 위한 ```cmd/berith``` 패키지를 빌드하고 싶다면 ```$ make berith``` 명령을 사용하면 된다. 

#### 고 바이너리를 이용하여 빌드하기

고 바이너리에는 패키지를 빌드시키는 커맨드인 ```build/install``` 을 내장하고 있다. 두 명령어는 모두 실행가능한 바이너리 파일을 생성하는 역할을 한다. 

### 베리드의 실행

베리드 바이너리를 실행하는 방법과 설정값을 지정하는 방법에대해 기술한다.

#### 베리드의 실행 방법

베리드는 빌드된 ```cmd/berith``` 패키지의 바이너리를 실행하는 것으로 작동한다. 콘솔창으로 로그가 출력되는 것으로 확인이 가능하다.

#### 베리드의 설정

베리드는 다양한 설정값을 갖는다. 설정값을 지정하는 방법은 두가지로 나뉜다.

1. 베리드를 실행하는 명령어에 ```Flag``` 를 지정하는 것으로 설정값을 지정하는 것이다. ```$ berith help``` 명령어를 통해 베리드가 지원하는 ```Flag``` 에 대해서 확인할 수 있다.

2. 베리드를 실행할 때 설정 파일을 넘겨줄 수 있다. ```$ berith dumpconfig``` 명령어를 통해 기본적으로 사용하고 있는 설정파일을 얻어낼 수 있다. 이를 ```toml``` 파일로 저장하여 ```$ berith --config dir/filename.toml``` 처럼 ```config``` ```Flag``` 를 이용하여 베리드에 전달할 수 있다. ```config``` 를 이용해 지정한 설정값은 ```Flag``` 로 직접 지정한 설정값보다 우선순위가 낮다.

### 베리드의 테스트

베리드의 테스트는 테스트 코드를 이용한 단위 테스트, 테스트넷을 이용한 테스트가 있다.

#### 베리드의 테스트 코드

베리드의 테스트 코드는 고언어의 단위 테스트 모듈을 이용하여 작성한다. 고언어에서는 파일명을 ```filename_test.go``` 의 형태로 짓는 것으로 해당 파일이 테스트 코드임을 나타낸다. 테스트 코드 내부에 테스트 함수를 만들어 테스트를 진행하는데 함수의 이름은 ```TestFunctionName()```  의 형태로 지어야 하며, 함수의 인자로 ```*testing.T``` 타입을 받을 수 있다. 전달된 인자는 테스트에 도움이 되는 구조체로 자세한 내용은 [https://golang.org/pkg/testing](https://golang.org/pkg/testing) 이곳에서 확인할 수 있다.

#### 베리드의 테스트넷

테스트넷이란 메인넷과 별도로 구성된 테스트를 위한 네트워크이다. 이는 공개적인 테스트넷일 수도 로컬에서 만든 단독 네트워크일 수도 있다. 

##### 공개 테스트넷

베리드는 ```testnet``` 이라는 ```Flag``` 를 이용하여 코드에 등록된 설정으로 테스트넷에 접속할 수 있다. 이는 ```testnet``` 이라는 ```Flag``` 를 지정하여 노드를 실행시킨 다른 다수의 유저들과 연결되어 실제 메인넷과 비슷한 테스트 환경을 구축할 수 있다. 하지만 베리드는 이를 지원하지 않는다.

##### 로컬 테스트넷

로컬 테스트넷이란 테스트를 위하여 일정 수의 노드를 연결시켜 네트워크를 구성하는 방법이다. 테스트넷을 구축하기 위해서는 두가지 설정이 필요하다.

제네시스 블록을 설정할 필요가 있다. 제네시스 블록은 블록체인의 가장 첫 블록이다.
```
MainnetChainConfig = &ChainConfig{
	ChainID:             big.NewInt(106),
	HomesteadBlock:      big.NewInt(0),
	DAOForkBlock:        nil,
	DAOForkSupport:      true,
	EIP150Block:         big.NewInt(0),
	EIP155Block:         big.NewInt(0),
	EIP158Block:         big.NewInt(0),
	ByzantiumBlock:      big.NewInt(0),
	ConstantinopleBlock: big.NewInt(0),
	BIP1Block:           big.NewInt(508000),
	BIP2Block:           big.NewInt(545000),
	BIP3Block:           big.NewInt(1168000),
	Bsrr: &BSRRConfig{
		Period:       5,
		Epoch:        360,
		Rewards:      common.StringToBig("360"),
		StakeMinimum: common.StringToBig("100000000000000000000000"),
		SlashRound:   0,
		ForkFactor:   1.0,
	},
}

func DefaultGenesisBlock() *Genesis {
    return &Genesis{
        Config:     params.MainnetChainConfig,
        Nonce:      0x00,
        Timestamp:  0x00,
        ExtraData:  hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000005559b45c940464df7145affec2f2ff4f691d92beefc1b29449332e6dd77cea91bc89db4ac9c43fa80000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
        GasLimit:   94000000,
        Difficulty: big.NewInt(1),
        Mixhash:    common.BytesToHash(hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000")),
        Coinbase:   common.HexToAddress("0x0000000000000000000000000000000000000000"),
        Alloc: map[common.Address]GenesisAccount{
            common.HexToAddress("Bx68dabe06325d57f2ad36cb60b754184f4f771a69"): {Balance: common.StringToBig("5000000000000000000000000000")},
        },
    }
}
```
위의 코드는 아무런 설정없이 베리드를 실행한 경우 지정되는 제네시스 블록의 정보이다. 대부분 이더리움과 똑같다. 이 중에서 테스트넷을 만드는데 필요한것은 아래와 같다.

* ExtraData : 최초에 어떤 계정도 스테이크를 하지 못하기에 블록을 생성할 계정의 정보를 등록해야한다. 이 ExtraData의 첫 64개의 “0” 뒤에 블록을 생성할 계정목록을 구분자 없이 등록한다.

* Alloc : 특정 계정의 코인잔액을 지정할 수 있다.

제네시스 파일은 위의 코드를 수정하여 지정할 수도 있지만, ```$ berith init``` 명령어와 제네시스 정보가 담긴 ```json``` 파일을 이용하여 제네시스 블록이 등록된 ```chaindata``` 를 만들어서 사용할 수 있다. 

또 부트노드를 설정해야한다. 부트노드는 노드가 가장 처음 연결한 노드 목록을 요청하는 노드로 테스트넷 설정을 가진 노드를 얻기 위해 따로 구성하는 것이 좋다. 로컬 테스트 환경에서는 ```admin_addPeer``` RPC를 이용하여 원하는 노드를 수동으로 연결할 수도 있다.
```
var MainnetBootnodes = []string{
    // Berith Foundation Go Bootnodes
    "enode://ea9b7c833a522780cb50dbb5f6e44c8d475ce8dedda44cb555e59994a5f89288908ebb288cfec9962c7321dee311a2a9bbfbadda78b1b3ef6dbcb33aea063e21@13.124.140.180:40404",
    "enode://2c7f9c316e460f98516e27ecd4bcb3e6772d2ae50db7ed11f96b4cb973aaca51b21cb485815d9f627c607e9def084c6e183cd2c12ec9dcc22fd9af198b6d34d3@15.164.130.81:40404",
}
```
위의 코드는 설정값 없이 베리드를 실행한 경우 기본적으로 사용되는 메인넷의 부트노드다. 이 코드를 수정하여 부트노드를 변경할 수도 있지만, 베리드를 실행 할 때, ```config``` 파일에서 부트노드를 지정하거나 ```bootnodes``` 라는 ```Flag``` 와 부트노드 정보를 입력하여 부트노드를 지정하여 사용할 수 있다.


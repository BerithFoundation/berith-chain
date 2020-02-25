## Execution and Testing of Berith

Describes how to run and test Berith.

### Building Berith

Berith is a go program. It is built with go binary and generates a binary executable.

#### Building with make

The ```Makefile1``` for the make command is defined at the top of the project. You can build a specific package separately by adding commands to the make command. For example, if you want to build a ```cmd/berith``` package for running nodes, you can use the ```$ make berith``` command.

#### Building with go binary

Go binary has a ```build/install``` command that builds packages. Both commands are responsible for generating executable binary files.

### Running Berith

Describes how to run the Berith binaries and how to specify setting values.

#### How to run Berith

Berith can be run by executing the binary of the ```cmd / Berith``` package that was built.

#### Setting up Berith

Berith has various setting values. There are two ways to specify setting values.

1. You can specify setting values by specifying ‘Flag’ on the command that executes Berith. You can check the ‘Flag’ supported by Berith through the ‘$ berith help’ command.

2. You can hand over a configuration file when running Berith. With the ‘$berith dumpconfig’ command, you can get the default configuration file. You can save it as a ‘toml’ file and pass it to Berith using ‘config Flag’ like ‘$ berith --config dir / filename.toml’. The setting value specified by ‘config’ has a lower priority than the setting value directly specified by Flag.

### Testing Berith

Berith includes a unit test using test code and a test using Testnet.

#### Berith test code

Berith test code is written using the unit test module of the go language. In the go language, the filename is typed as ‘filename_test.go’ to indicate that the file is a test code. Test by creating a test function inside the test code. At this time, the name of the function should be in the form of ‘TestFunctionName()’, and it can receive the ‘*testing.T’ type as a parameter. The received argument is a structure that helps test. For more information, you can visit [https://golang.org/pkg/testing](https://golang.org/pkg/testing)

#### Berith Testnet

Testnet is a network for testing that is configured separately from the mainnet. It could be the public testnet or a local standalone testnet.

##### Public Testnet

Berith can access the testnet with the settings registered in the code using the Flag ‘testnet’. By specifying a Flag called testnet, you can connect with several other users who have run nodes and create a test environment similar to the actual mainnet. But Berith doesn```t support it.

##### Local Testnet

Local testnet is a method of configuring a network by connecting a certain number of nodes for a test. Two settings are required to build a testnet.

First, we need to set up the Genesis block. The Genesis Block is the first block of the blockchain.
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
The above code is the information of the Genesis block that is specified if you run Berith without any setup. Most are the same as ethereum. Of those, the information needed to create a testnet is as follows:.

* ExtraData : Initially, no account can stake, so you need to register the information for the account that will create the block. Register the list of accounts to create a block after the first 64 zeros of this ExtraData without delimiters.

* Alloc : You can specify a coin balance for a particular account.

The genesis file can be specified by modifying the code above. Another way is to use ```$ berith init``` command and json file containing Genesis information to create and use ‘chaindata’ with Genesis Block registered.

Second, We need to set up a bootnode. The bootnode is the node that requests the list of nodes that the node first connected to. And it is recommended that you configure them separately to obtain a node with a testnet setting. In a local test environment, you can manually connect to any node using the ```admin_addPeer``` RPC.
```
var MainnetBootnodes = []string{
    // Berith Foundation Go Bootnodes
    "enode://ea9b7c833a522780cb50dbb5f6e44c8d475ce8dedda44cb555e59994a5f89288908ebb288cfec9962c7321dee311a2a9bbfbadda78b1b3ef6dbcb33aea063e21@13.124.140.180:40404",
    "enode://2c7f9c316e460f98516e27ecd4bcb3e6772d2ae50db7ed11f96b4cb973aaca51b21cb485815d9f627c607e9def084c6e183cd2c12ec9dcc22fd9af198b6d34d3@15.164.130.81:40404",
}
```
The above code is the bootnode of the mainnet that is used by default if you run Berith without a set value. You can modify this code to change the boot node. Alternatively, when running Berith, you can specify the boot node in the ```config``` file, or specify and use the boot node by entering the Flag ```bootnodes``` and boot node information.


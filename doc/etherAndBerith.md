## Berith and Ethereum
Berith is a project that customizes Ethereum```s go-ethereum, the golang implementation.
 [https://github.com/ethereum/go-ethereum](https://github.com/ethereum/go-ethereum)

Berith changed the name related to Ethereum to Berith in the names of packages and variables.
* ```ether```, ```ethereum``` → ```berith```
* ```eth``` → ```ber```
  
You can find these changes in the source code.

### Ethereum Setting value

Although Berith is a project that customized Ethereum, due to the nature of the consensus algorithm, the settings used by Ethereum are sometimes not available in Berith.

#### SyncMode

Berith cannot verify a block with the block header alone, so if you set ```SyncMode``` to ```light```, you cannot run the node. Berith can use the past block to verify the current block, so if you set ```SyncMode``` to ```fast```, you can`t run the node. That is why ```SyncMode``` must set to ```full```.

#### GCMode

Berith uses the ```state``` of the past block to verify the block. So if you set the ```GCmode``` to ```full``` and the ```state``` in memory is cleared when you shut down the node, you may have problems verifying the block. That is why ```GCMode``` must set to ```archive```.

### Network separation between Ethereum and Berith

Ethereum uses UDP to find other Ethereum nodes, called ```discover```. In Ethereum, ‘discover’ separates the type of request into the value of the first packet of the request. If the value of the first packet does not map to the type of request, the request is ignored. Berith used this to separate the network from Ethereum by modifying the constants that the request type mapped.

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
    pingPacket = iota + 11 // zero is ```reserved```
    pongPacket
    findnodePacket
    neighborsPacket
)
```
The above code is the part that declares the packet type of the ```discover``` request from Ethereum and Berith. You can see that Berith and Ethereum map different values ​​for the same request.

### Berith RPC

The RPC method in Berith consists of a package_method, just like Ethereum. RPC methods in Berith are mostly the same as Ethereum, but there are changes such as changing the form of transaction parameter and addition methods for stakes. For detailed information refer to Berith RPC Spec. For more information: [Berith RPC Spec](./rpc.md)

### Berith Address

Berith addresses start with “Bx”, unlike Ethereum with “0x”. This is a change using the fact that the “0x” in front of the Ethereum account is not related to the values ​​that make up the actual account and is simply a prefix to indicate that the data consists of hexadecimal digits. When unmarshal a remote json request in the form of a structure to handle it in go, it changed the validation function of the address data.
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
The above code is a function that is called when unmarshalling data of type ```Address```. Through this, you can see that the data received as "Bx" prefix is ​​changed to "0x" and pass the existing validation check. You can also see that the data entered as "0x" is not valid when checking the "Bx" prefix and returns an error.


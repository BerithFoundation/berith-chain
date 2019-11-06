# Berith rpc specs

## Index

- <a href="#json-rpc-methods">JSON-RPC methods</a>
- <a href="#json-rpc-api-reference">JSON RPC API Reference</a>

---  

<div id="json-rpc-methods"></div>  

## JSON-RPC methods  

- <a href="#web3_clientVersion">web3_clientVersion</a>  
- <a href="#web3_sha3">web3_sha3</a>  

- <a href="#net_version">net_version</a>
- <a href="#net_peerCount">net_peerCount</a>  
- <a href="#net_listening">net_listening</a>  

- <a href="#berith_protocolVersion">berith_protocolVersion</a>  
- <a href="#berith_coinbase">berith_coinbase</a>  
- <a href="#berith_syncing">berith_syncing</a>
- <a href="#berith_mining">berith_mining</a>  
- <a href="#berith_gasPrice">berith_gasPrice</a>  
- <a href="#berith_accounts">berith_accounts</a>  
- <a href="#berith_blockNumber">berith_blockNumber</a>  
- <a href="#berith_getBalance">berith_getBalance</a>  
- <a href="#berith_getStakeBalance">berith_getStakeBalance</a>
- <a href="#berith_getTransactionCount">berith_getTransactionCount</a>  
- <a href="#berith_getBlockTransactionCountByHash">berith_getBlockTransactionCountByHash</a>  
- <a href="#berith_getBlockTransactionCountByNumber">berith_getBlockTransactionCountByNumber</a>  
- <a href="#berith_getUncleCountByBlockHash">berith_getUncleCountByBlockHash</a>  
- <a href="#berith_getUncleCountByBlockNumber">berith_getUncleCountByBlockNumber</a>  
- <a href="#berith_getCode">berith_getCode</a>  
- <a href="#berith_sign">berith_sign</a>
- <a href="#berith_sendTransaction">berith_sendTransaction</a>  
- <a href="#berith_sendRawTransaction">berith_sendRawTransaction</a>
- <a href="#berith_call">berith_call</a>  
- <a href="#berith_eth_estimateGas">berith_eth_estimateGas</a>
- <a href="#berith_getBlockByHash">berith_getBlockByHash</a>  
- <a href="#berith_getBlockByNumber">berith_getBlockByNumber</a>  
- <a href="#berith_getTransactionByHash">berith_getTransactionByHash</a>  
- <a href="#berith_getTransactionByBlockHashAndIndex">berith_getTransactionByBlockHashAndIndex</a>  
- <a href="#berith_getTransactionByBlockNumberAndIndex">berith_getTransactionByBlockNumberAndIndex</a>
- <a href="#berith_getTransactionReceipt">berith_getTransactionReceipt</a>  
- <a href="#berith_newFilter">berith_newFilter</a>  
- <a href="#berith_newBlockFilter">berith_newBlockFilter</a>
- <a href="#berith_newPendingTransactionFilter">berith_newPendingTransactionFilter</a>  
- <a href="#berith_uninstallFilter">berith_uninstallFilter</a>  
- <a href="#berith_getFilterChanges">berith_getFilterChanges</a>  
- <a href="#berith_getFilterLogs">berith_getFilterLogs</a>  
- <a href="#berith_getLogs">berith_getLogs</a>  

- <a href="#amon_getBlockCreatorsByNumber">amon_getBlockCreatorsByNumber</a>
- <a href="#amon_getBlockCreatorsByHash">amon_getBlockCreatorsByHash</a>
- <a href="#amon_getJoinRatio">amon_getJoinRatio</a>
  
---  

<div id="json-rpc-api-reference"></div>

## JSON RPC API Reference  

<div id="web3_clientVersion"></div>  

### web3_clientVersion  
Returns the current client version.  

**Parameter**  
none  

**Returns**  
```String``` - The current client version.   

**Example**  
```bash
// Request
curl --data '{"jsonrpc":"2.0","method":"web3_clientVersion","params":[],"id":67}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":67,
  "jsonrpc":"2.0",
  "result": "Berith/v0.1.0-unstable-c6469980/linux-amd64/go1.12.1"
}
```

---  

<div id="web3_sha3"></div>  

### web3_sha3  
Returns Keccak-256 (not the standardized SHA3-256) of the given data.

**Parameter**  
1. ```DATA``` - the data to convert into a SHA3 hash. 

**Example Parameters**  
```js
params: [
  "0x68656c6c6f20776f726c64"
]
```

**Returns**  
```DATA``` - The SHA3 result of the given string.

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"web3_sha3","params":["0x68656c6c6f20776f726c64"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc": "2.0",
  "result": "0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad"
}
```
---  
<div id="net_version"></div>  

### net_version  
Returns the current network id.

**Parameter**  
none  

**Returns**  
```String``` - The current network id.


**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"net_version","params":[],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc": "2.0",
  "result": "3"
}
```  

<div id="net_peerCount"></div>  

### net_peerCount 
Returns number of peers currently connected to the client.

**Parameter**  
none

**Returns**  
```QUANTITY``` - integer of the number of connected peers.

**Example**  

```js
// Request
curl --data '{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc": "2.0",
  "result": "0x2" // 2
}
```

---  

<div id="net_listening"></div>  

### net_listening  
Returns `true` if client is actively listening for network connections.

**Parameter**  
none  

**Returns**  
`Boolean` - `true` when listening, otherwise `false`.

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"net_listening","params":[],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc":"2.0",
  "result":true
}
```

---  

<div id="berith_protocolVersion"></div>  

### berith_protocolVersion  
Returns the current berith protocol version.

**Parameter**  

none  

**Returns**  

`String` - The current berith protocol version.
 

**Example**  

```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_protocolVersion","params":[],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc": "2.0",
  "result": "0x3f"
}
```

---  

<div id="berith_coinbase"></div>  

### berith_coinbase  
Returns the client coinbase address.

**Parameter**  
none

**Returns**  
`DATA`, 20 bytes - the current coinbase address.
                                                                                                                                                                                                                         
**Example**
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_coinbase","params":[],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc": "2.0",
  "result": "Bxd8a25ff31c6174ce7bce74ca4a91c2e816dbf91e"
}
```  

---  

<div id="berith_syncing"></div>  

### berith_syncing  
Returns an object with data about the sync status or `false`.

**Parameter**  
none

**Returns**  
`Object|Boolean`, An object with sync status data or `FALSE`, when not syncing:
  - `startingBlock`: `QUANTITY` - The block at which the import started (will only be reset, after the sync reached his head)
  - `currentBlock`: `QUANTITY` - The current block, same as eth_blockNumber
  - `highestBlock`: `QUANTITY` - The estimated highest block
                                                                                                                                                                                                                         
**Example**
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_syncing","params":[],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc": "2.0",
  "result": {
    startingBlock: '0x384',
    currentBlock: '0x386',
    highestBlock: '0x454'
  }
}
// Or when not syncing
{
  "id":1,
  "jsonrpc": "2.0",
  "result": false
}
```  

---  

<div id="berith_mining"></div>  

### berith_mining  
Returns `true` if client is actively mining new blocks.

**Parameter**  
none

**Returns**  
`Boolean` - returns `true` of the client is mining, otherwise `false`.

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_mining","params":[],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":71,
  "jsonrpc": "2.0",
  "result": true
}
```  

---  

<div id="berith_gasPrice"></div>  

### berith_gasPrice  
Returns the current price per gas in wei.

**Parameter**  
none  

**Returns**  
`QUANTITY` - integer of the current gas price in wei.

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_gasPrice","params":[],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Response
{
  "id":1,
  "jsonrpc": "2.0",
  "result": "0x09184e72a000" // 10000000000000
}
```  

---  
<div id="berith_accounts"></div>  

### berith_accounts  
Returns a list of addresses owned by client.  

**Parameter**  
none

**Returns**  
`Array of DATA`, 20 Bytes - addresses owned by the client.

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_accounts","params":[],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc": "2.0",
  "result": ["Bxc94770007dda54cF92009BFF0dE90c06F603a09f"]
}
```
---  

<div id="berith_blockNumber"></div>  

### berith_blockNumber  
Returns the number of most recent block.

**Parameter**  
none

**Returns**  
`QUANTITY` - integer of the current block number the client is on.

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_blockNumber","params":[],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc": "2.0",
  "result": "0xc94" // 1207
}
```  
 
---  
  
<div id="berith_getBalance"></div>  

### berith_getBalance  
Returns the balance of the account of given address.

**Parameter**  
1. `DATA`, 20 Bytes - address to check for balance.
2. `QUANTITY|TAG` - integer block number, or the string `"latest"`, `"earliest"` or `"pending"`  

**Example Parameters**  

```js
params: [
   'Bxc94770007dda54cF92009BFF0dE90c06F603a09f',
   'latest'
]
```

**Returns**  
`QUANTITY` - integer of the current balance in wei.  


**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_getBalance","params":["Bxd8a25ff31c6174ce7bce74ca4a91c2e816dbf91e", "latest"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc": "2.0",
  "result": "0x0234c8a3397aab58" // 158972490234375000
}
```
  
---  
  
<div id="berith_getStakeBalance"></div>  

### berith_getStakeBalance  
Returns the staked amount of the account of given address.

**Parameter**  
1. `DATA`, 20 Bytes - address to check for staked amount.
2. `QUANTITY|TAG` - integer block number, or the string `"latest"`, `"earliest"` or `"pending"`

**Example Parameters**  

```js
params: [
   '0xc94770007dda54cF92009BFF0dE90c06F603a09f',
]
```

**Returns**  
`QUANTITY` - integer of the current balance in wei.  


**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_getStakeBalance","params":["Bxd8a25ff31c6174ce7bce74ca4a91c2e816dbf91e", "latest"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc": "2.0",
  "result": "0x0234c8a3397aab58" // 158972490234375000
}
```

---  
  
<div id="berith_getTransactionCount"></div>  

### berith_getTransactionCount  
Returns the number of transactions *sent* from an address.

**Parameter**  
1. `DATA`, 20 Bytes - address.
2. `QUANTITY|TAG` - integer block number, or the string `"latest"`, `"earliest"` or `"pending"`

**Example Parameters**  
```js
params: [
   'Bxd8a25ff31c6174ce7bce74ca4a91c2e816dbf91e',
   'latest' // state at the latest block
]
```

**Returns**  
`QUANTITY` - integer of the number of transactions send from this address.


**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_getTransactionCount","params":["Bxd8a25ff31c6174ce7bce74ca4a91c2e816dbf91e", "latest"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc": "2.0",
  "result": "0x1" // 1
}
```

---  
  
<div id="berith_getBlockTransactionCountByHash"></div>  

### berith_getBlockTransactionCountByHash  
Returns the number of transactions in a block from a block matching the given block hash.  

**Parameter**  
1. `DATA`, 32 Bytes - hash of a block.

**Example Parameters**  
```js
params: [
   '0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238'
]
```

**Returns**  
`QUANTITY` - integer of the number of transactions in this block.  

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_getBlockTransactionCountByHash","params":["0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc": "2.0",
  "result": "0xc" // 11
}
```  

---  
  
<div id="berith_getBlockTransactionCountByNumber"></div>  

### berith_getBlockTransactionCountByNumber  
Returns the number of transactions in a block matching the given block number.  

**Parameter**  
1. `QUANTITY|TAG` - integer of a block number, or the string `"earliest"`, `"latest"` or `"pending"`  

**Example Parameters**  
```js
params: [
   '0xe8', // 232
]
```

**Returns**  
`QUANTITY` - integer of the number of transactions in this block.  

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_getBlockTransactionCountByNumber","params":["0xe8"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc": "2.0",
  "result": "0xa" // 10
}
```  
---  
  
<div id="berith_getUncleCountByBlockHash"></div>  

### berith_getUncleCountByBlockHash  
Returns the number of uncles in a block from a block matching the given block hash.
  
**Parameter**  
1. `DATA`, 32 Bytes - hash of a block.

**Example Parameters**  
```js
params: [
   '0xfadc1760414e85645d3d1d7c9d778c86d0fa7c89506d876c1cbaa672881eaa90'
]
```

**Returns**  
`QUANTITY` - integer of the number of uncles in this block.

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_getUncleCountByBlockHash","params":["0xfadc1760414e85645d3d1d7c9d778c86d0fa7c89506d876c1cbaa672881eaa90"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc": "2.0",
  "result": "0xc" // 1
}
```  

---  
  
<div id="berith_getUncleCountByBlockNumber"></div>  

### berith_getUncleCountByBlockNumber  
Returns the number of uncles in a block from a block matching the given block number.

**Parameter**  
1. `QUANTITY|TAG` - integer of a block number, or the string "latest", "earliest" or "pending"

```js
params: [
   '0xe8', // 232
]
```

**Returns**  
`QUANTITY` - integer of the number of uncles in this block.  

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_getUncleCountByBlockNumber","params":["0xe8"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc": "2.0",
  "result": "0x1" // 1
}
```

---  

<div id="berith_getCode"></div>  

### berith_getCode  
Returns code at a given address.

**Parameter**  
1. `DATA`, 20 Bytes - address.
2. `QUANTITY|TAG` - integer block number, or the string `"latest"`, `"earliest"` or `"pending"`  

**Example Parameters**  
```js
params: [
   'Bxa94f5374fce5edbc8e2a8697c15331677e6ebf0b',
   '0x2'  // 2
]
```

**Returns**  
`DATA` - the code from the given address.

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_getCode","params":["Bxa94f5374fce5edbc8e2a8697c15331677e6ebf0b", "0x2"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545
```

---  

<div id="berith_sign"></div>  

### berith_sign  
The sign method calculates an Berith specific signature with: `sign(keccak256("\x19Berith Signed Message:\n" + len(message) + message)))`.

By adding a prefix to the message makes the calculated signature recognisable as an Berith specific signature.  
This prevents misuse where a malicious DApp can sign arbitrary data (e.g. transaction) and use the signature to impersonate the victim.

**Note** the address to sign with must be unlocked.

**Parameter**  
account, message

1. `DATA`, 20 Bytes - address.
2. `DATA`, N Bytes - message to sign.  

**Returns**  
`DATA`: Signature

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_sign","params":["Bxd8a25ff31c6174ce7bce74ca4a91c2e816dbf91e", "0xdeadbeaf"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0x45fefbfa3b81117f9033b95ca949d8fbe18ad755de7650c57bd67423e3de3f1e7c89b8f5d1764b2b7d702d4a75ef2a10e02b5be913b4e16f2e84bafa8498a5281b"
}
```

---  

<div id="berith_sendTransaction"></div>  

### berith_sendTransaction  
Creates new message call transaction or a contract creation, if the data field contains code.

**Parameter**  
1. `Object` - The transaction object
  - `from`: `DATA`, 20 Bytes - The address the transaction is send from.
  - `to`: `DATA`, 20 Bytes - (optional when creating new contract) The address the transaction is directed to.
  - `gas`: `QUANTITY`  - (optional, default: 90000) Integer of the gas provided for the transaction execution. It will return unused gas.
  - `gasPrice`: `QUANTITY`  - (optional, default: To-Be-Determined) Integer of the gasPrice used for each paid gas
  - `value`: `QUANTITY`  - (optional) Integer of the value sent with this transaction
  - `data`: `DATA`  - The compiled code of a contract OR the hash of the invoked method signature and encoded parameters.
  - `nonce`: `QUANTITY`  - (optional) Integer of a nonce. This allows to overwrite your own pending transactions that use the same nonce.
  - `base` : `QUANTITY` - (optional) transfer's from balance type (MAIN = "main", "STAKE" = "stake"). default is main 
  - `target` : `QUANTITY` - (optional) transfer's to balance type (MAIN = "main", "STAKE" = "stake"). default is main
    
**Example Parameters**  
```js
// transfer
params: [{
  "from": "Bxb60e8dd61c5d32be8058bb8eb970870f07233155",
  "to": "Bxd46e8dd67c5d32be8058bb8eb970870f07244567",
  "gas": "0x76c0", // 30400
  "gasPrice": "0x9184e72a000", // 10000000000000
  "value": "0x9184e72a", // 2441406250
  "data": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"
}]
```

**Returns**  
`DATA`, 32 Bytes - the transaction hash, or the zero hash if the transaction is not yet available.

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_sendTransaction","params":[{ 
"from": "Bxd8a25ff31c6174ce7bce74ca4a91c2e816dbf91e", "to" : "Bxd8a25ff31c6174ce7bce74ca4a91c2e816dbf91e", "value" : "0x1",
"base":"0x1", "target" : "0x2"
}],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc": "2.0",
  "result": "0xcc7b69665bb876991b7621a4c225993f3c789a7647871f71350d9f7153583039"
}
```  

---  

<div id="berith_sendRawTransaction"></div>  

### berith_sendRawTransaction   
Creates new message call transaction or a contract creation for signed transactions.

**Parameter**  
1. `DATA`, The signed transaction data.
    
**Example Parameters**  
```js
params: ["0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"]
```

**Returns**  
`DATA`, 32 Bytes - the transaction hash, or the zero hash if the transaction is not yet available.

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_sendRawTransaction","params":["0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc": "2.0",
  "result": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331"
}
```

---  
  
<div id="berith_call"></div>  

### berith_call  
Executes a new message call immediately without creating a transaction on the block chain.

**Parameter**  
1. `Object` - The transaction call object
  - `from`: `DATA`, 20 Bytes - (optional) The address the transaction is sent from.
  - `to`: `DATA`, 20 Bytes  - The address the transaction is directed to.
  - `gas`: `QUANTITY`  - (optional) Integer of the gas provided for the transaction execution. eth_call consumes zero gas, but this parameter may be needed by some executions.
  - `gasPrice`: `QUANTITY`  - (optional) Integer of the gasPrice used for each paid gas
  - `value`: `QUANTITY`  - (optional) Integer of the value sent with this transaction
  - `data`: `DATA`  - (optional) Hash of the method signature and encoded parameters.
2. `QUANTITY|TAG` - integer block number, or the string `"latest"`, `"earliest"` or `"pending"`

**Returns**  
`DATA` - the return value of executed contract.

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"","params":[{{ see above }}],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc": "2.0",
  "result": "0x"
}
```

---  
  
<div id="berith_estimateGas"></div>  

### berith_estimateGas  
Generates and returns an estimate of how much gas is necessary to allow the transaction to complete.  
The transaction will not be added to the blockchain.  Note that the estimate may be significantly  
more than the amount of gas actually used by the transaction, for a variety of reasons including EVM mechanics and node performance.

**Parameter**  
See [berith_call](#berith_call) parameters, expect that all properties are optional. If no gas limit is specified berith uses the block gas limit from the pending block as an upper bound.  
As a result the returned estimate might not be enough to executed the call/transaction when the amount of gas is higher than the pending block gas limit.

**Returns**  
`QUANTITY` - the amount of gas used.

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_estimateGas","params":[{see above}],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
// Result
{
  "id":1,
  "jsonrpc": "2.0",
  "result": "0x5208" // 21000
}
```

---  

<div id="berith_getBlockByHash"></div>  

### berith_getBlockByHash  
Returns information about a block by hash.
  
**Parameter**  
1. `DATA`, 32 Bytes - Hash of a block.
2. `Boolean` - If `true` it returns the full transaction objects, if `false` only the hashes of the transactions.  

**Example Parameters**  
```js
params: [
   '0xd4c5d0d4b37e7636617ea0731e0bd03395605654bb7cafe11663234b7e9e23c7',
   true
]
```  

**Returns**  
`Object` - A block object, or `null` when no block was found:

  - `number`: `QUANTITY` - the block number. `null` when its pending block.
  - `hash`: `DATA`, 32 Bytes - hash of the block. `null` when its pending block.
  - `parentHash`: `DATA`, 32 Bytes - hash of the parent block.
  - `nonce`: `DATA`, 8 Bytes - rank of miners in stake holers.
  - `sha3Uncles`: `DATA`, 32 Bytes - SHA3 of the uncles data in the block.
  - `logsBloom`: `DATA`, 256 Bytes - the bloom filter for the logs of the block. `null` when its pending block.
  - `transactionsRoot`: `DATA`, 32 Bytes - the root of the transaction trie of the block.
  - `stateRoot`: `DATA`, 32 Bytes - the root of the final state trie of the block.
  - `receiptsRoot`: `DATA`, 32 Bytes - the root of the receipts trie of the block.
  - `miner`: `DATA`, 20 Bytes - the address of the beneficiary to whom the mining rewards were given.
  - `difficulty`: `QUANTITY` - integer of the difficulty for this block.
  - `totalDifficulty`: `QUANTITY` - integer of the total difficulty of the chain until this block.
  - `extraData`: `DATA` - the "extra data" field of this block.
  - `size`: `QUANTITY` - integer the size of this block in bytes.
  - `gasLimit`: `QUANTITY` - the maximum gas allowed in this block.
  - `gasUsed`: `QUANTITY` - the total used gas by all transactions in this block.
  - `timestamp`: `QUANTITY` - the unix timestamp for when the block was collated.
  - `transactions`: `Array` - Array of transaction objects, or 32 Bytes transaction hashes depending on the last given parameter.
  - `uncles`: `Array` - Array of uncle hashes.


**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_getBlockByHash","params":["0xd4c5d0d4b37e7636617ea0731e0bd03395605654bb7cafe11663234b7e9e23c7", true],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "difficulty": "0x5265c0",
    "extraData": "0xd98201008662657269746888676f312e31322e31856c696e7578000000000000d23db08f84e572b58452c04950bb534d289bf9c8b1dfd178d1562bfb35342daa096ce8c25b78b6cc727a25752f22749fdf8571484bff49437dc3b7d477e9b1dc00",
    "gasLimit": "0x5f65c42",
    "gasUsed": "0x5208",
    "hash": "0xd4c5d0d4b37e7636617ea0731e0bd03395605654bb7cafe11663234b7e9e23c7",
    "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "miner": "Bxd8a25ff31c6174ce7bce74ca4a91c2e816dbf91e",
    "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "nonce": "0x0000000000000001",
    "number": "0x850",
    "parentHash": "0x4b297ac504bb9a7d74a161a3f45f48717a5acf0aaf16e1e6d2faf4b34a4ffcd9",
    "receiptsRoot": "0x056b23fbba480696b65fe5a59b8f2148a1299103c4f57df839233af2cf4ca2d2",
    "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "size": "0x2ce",
    "stateRoot": "0x416718a61a1aedf6c189fa2a87974370e0824cd116aedf82853dedf220496892",
    "timestamp": "0x5dba6667",
    "totalDifficulty": "0x2a9b602b5",
    "transactions": [
      {
        "blockHash": "0xd4c5d0d4b37e7636617ea0731e0bd03395605654bb7cafe11663234b7e9e23c7",
        "blockNumber": "0x850",
        "from": "Bxd8a25ff31c6174ce7bce74ca4a91c2e816dbf91e",
        "gas": "0x15f90",
        "gasPrice": "0x1",
        "hash": "0xcfe05741f5c30a50ceab4e847eb4f48ed7b0b9841c92e1aaaeaef12226950ade",
        "input": "0x",
        "nonce": "0x4",
        "to": "Bxd8a25ff31c6174ce7bce74ca4a91c2e816dbf91e",
        "transactionIndex": "0x0",
        "value": "0x1",
        "base": 1,
        "target": 1,
        "v": "0x11cca",
        "r": "0x6b69f34217fb9ffccee69719719ecc1afe73d634d70b7c9834b467d2afc78921",
        "s": "0x526c6133e1e400989e24365f9bb06bec009a4526be572417a86e5a96c56814d3"
      }
    ],
    "transactionsRoot": "0xc40267b67d6b4d1930cb54f913c69c6560dd35633d67a4921e6803f3dbf7ef84",
    "uncles": []
  }
}
```

---  

<div id="berith_getBlockByNumber"></div>  

### berith_getBlockByNumber  
Returns information about a block by block number.  

**Parameter**  
1. `QUANTITY|TAG` - integer of a block number, or the string `"earliest"`, `"latest"` or `"pending"`
2. `Boolean` - If `true` it returns the full transaction objects, if `false` only the hashes of the transactions.  

**Example Parameters**
```js
params: [
   '0x1b4', // 436
   true
]
```  

**Returns**  
See [berith_getBlockByHash](#berith_getBlockByHash)  

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_getBlockByNumber","params":["0x1b4", true],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545
```  

---  

<div id="berith_getTransactionByHash"></div>  

### berith_getTransactionByHash  
Returns the information about a transaction requested by transaction hash.

**Parameter**  
```js
params: [
   "0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b"
]
```

**Returns**  
`Object` - A transaction object, or `null` when no transaction was found:

  - `blockHash`: `DATA`, 32 Bytes - hash of the block where this transaction was in. `null` when its pending.
  - `blockNumber`: `QUANTITY` - block number where this transaction was in. `null` when its pending.
  - `from`: `DATA`, 20 Bytes - address of the sender.
  - `gas`: `QUANTITY` - gas provided by the sender.
  - `gasPrice`: `QUANTITY` - gas price provided by the sender in Wei.
  - `hash`: `DATA`, 32 Bytes - hash of the transaction.
  - `input`: `DATA` - the data send along with the transaction.
  - `nonce`: `QUANTITY` - the number of transactions made by the sender prior to this one.
  - `to`: `DATA`, 20 Bytes - address of the receiver. `null` when its a contract creation transaction.
  - `transactionIndex`: `QUANTITY` - integer of the transaction's index position in the block. `null` when its pending.
  - `value`: `QUANTITY` - value transferred in Wei.
  - `base` : `QUANTITY` - transfer's from balance type ("MAIN" = 1, "STAKE" = 2). default is main 
  - `target` : `QUANTITY` - transfer's to balance type ("MAIN" = 1, "STAKE" = 2). default is main
  - `v`: `QUANTITY` - ECDSA recovery id
  - `r`: `QUANTITY` - ECDSA signature r
  - `s`: `QUANTITY` - ECDSA signature s

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_getTransactionByHash","params":["0xcfe05741f5c30a50ceab4e847eb4f48ed7b0b9841c92e1aaaeaef12226950ade"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "blockHash": "0xd4c5d0d4b37e7636617ea0731e0bd03395605654bb7cafe11663234b7e9e23c7",
    "blockNumber": "0x850",
    "from": "Bxd8a25ff31c6174ce7bce74ca4a91c2e816dbf91e",
    "gas": "0x15f90",
    "gasPrice": "0x1",
    "hash": "0xcfe05741f5c30a50ceab4e847eb4f48ed7b0b9841c92e1aaaeaef12226950ade",
    "input": "0x",
    "nonce": "0x4",
    "to": "Bxd8a25ff31c6174ce7bce74ca4a91c2e816dbf91e",
    "transactionIndex": "0x0",
    "value": "0x1",
    "base": 1,
    "target": 1,
    "v": "0x11cca",
    "r": "0x6b69f34217fb9ffccee69719719ecc1afe73d634d70b7c9834b467d2afc78921",
    "s": "0x526c6133e1e400989e24365f9bb06bec009a4526be572417a86e5a96c56814d3"
  }
}
```

---  
  
<div id="berith_getTransactionByBlockHashAndIndex"></div>  

### berith_getTransactionByBlockHashAndIndex  
Returns information about a transaction by block hash and transaction index position.

**Parameter**  
1. `DATA`, 32 Bytes - hash of a block.
2. `QUANTITY` - integer of the transaction index position.

**Example Parameters**  
```js
params: [
   '0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331',
   '0x0' // 0
]
```

**Returns**  
See [berith_getTransactionByHash](#berith_getTransactionByHash)

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_getTransactionByBlockHashAndIndex","params":["0xd4c5d0d4b37e7636617ea0731e0bd03395605654bb7cafe11663234b7e9e23c7", "0x0"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545
```

---  
 
<div id="berith_getTransactionByBlockNumberAndIndex"></div>  

###  berith_getTransactionByBlockNumberAndIndex  
Returns information about a transaction by block number and transaction index position.

**Parameter**  
1. `QUANTITY|TAG` - a block number, or the string `"earliest"`, `"latest"` or `"pending"`.
2. `QUANTITY` - the transaction index position.

**Example Parameters**  
```js
params: [
   '0x29c', // 668
   '0x0' // 0
]
```  

**Returns**  
See [berith_getTransactionByHash](#berith_getTransactionByHash)

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_getTransactionByBlockNumberAndIndex","params":["0x850", "0x0"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545
```

---  

<div id="berith_getTransactionReceipt"></div>  

### berith_getTransactionReceipt  
Returns the receipt of a transaction by transaction hash.

**Note** That the receipt is not available for pending transactions.

**Parameter**  
1. `DATA`, 32 Bytes - hash of a transaction  

**Example Parameters**  
```js
params: [
   '0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238'
]
```  

**Returns**  
`Object` - A transaction receipt object, or `null` when no receipt was found:

  - `transactionHash `: `DATA`, 32 Bytes - hash of the transaction.
  - `transactionIndex`: `QUANTITY` - integer of the transaction's index position in the block.
  - `blockHash`: `DATA`, 32 Bytes - hash of the block where this transaction was in.
  - `blockNumber`: `QUANTITY` - block number where this transaction was in.
  - `from`: `DATA`, 20 Bytes - address of the sender.
  - `to`: `DATA`, 20 Bytes - address of the receiver. null when it's a contract creation transaction.
  - `cumulativeGasUsed `: `QUANTITY ` - The total amount of gas used when this transaction was executed in the block.
  - `gasUsed `: `QUANTITY ` - The amount of gas used by this specific transaction alone.
  - `contractAddress `: `DATA`, 20 Bytes - The contract address created, if the transaction was a contract creation, otherwise `null`.
  - `logs`: `Array` - Array of log objects, which this transaction generated.
  - `logsBloom`: `DATA`, 256 Bytes - Bloom filter for light clients to quickly retrieve related logs.

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_getTransactionReceipt","params":["0xcfe05741f5c30a50ceab4e847eb4f48ed7b0b9841c92e1aaaeaef12226950ade"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "blockHash": "0xd4c5d0d4b37e7636617ea0731e0bd03395605654bb7cafe11663234b7e9e23c7",
    "blockNumber": "0x850",
    "contractAddress": null,
    "cumulativeGasUsed": "0x5208",
    "from": "Bxd8a25ff31c6174ce7bce74ca4a91c2e816dbf91e",
    "gasUsed": "0x5208",
    "logs": [],
    "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "status": "0x1",
    "to": "Bxd8a25ff31c6174ce7bce74ca4a91c2e816dbf91e",
    "transactionHash": "0xcfe05741f5c30a50ceab4e847eb4f48ed7b0b9841c92e1aaaeaef12226950ade",
    "transactionIndex": "0x0"
  }
}
```

---  
  
<div id="berith_newFilter"></div>  

### berith_newFilter  
Creates a filter object, based on filter options, to notify when the state changes (logs).
To check if the state has changed, call [berith_getFilterChanges](#berith_getFilterChanges).  

**A note on specifying topic filters:**  
Topics are order-dependent. A transaction with a log with topics [A, B] will be matched by the following topic filters:
* `[]` "anything"
* `[A]` "A in first position (and anything after)"
* `[null, B]` "anything in first position AND B in second position (and anything after)"
* `[A, B]` "A in first position AND B in second position (and anything after)"
* `[[A, B], [A, B]]` "(A OR B) in first position AND (A OR B) in second position (and anything after)"

**Parameter**  
1. `Object` - The filter options:
  - `fromBlock`: `QUANTITY|TAG` - (optional, default: `"latest"`) Integer block number, or `"latest"` for the last mined block or `"pending"`, `"earliest"` for not yet mined transactions.
  - `toBlock`: `QUANTITY|TAG` - (optional, default: `"latest"`) Integer block number, or `"latest"` for the last mined block or `"pending"`, `"earliest"` for not yet mined transactions.
  - `address`: `DATA|Array`, 20 Bytes - (optional) Contract address or a list of addresses from which logs should originate.
  - `topics`: `Array of DATA`,  - (optional) Array of 32 Bytes `DATA` topics. Topics are order-dependent. Each topic can also be an array of DATA with "or" options.
  
**Example Parameters**  
```js
params: [{
  "fromBlock": "0x1",
  "toBlock": "0x2",
  "address": "Bx8888f1f195afa192cfee860698584c030f4c9db1",
  "topics": ["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b", null, ["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b", "0x0000000000000000000000000aff3454fce5edbc8cca8697c15331677e6ebccc"]]
}]
```

**Returns**  
`QUANTITY` - A filter id.
  

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_newFilter","params":[{
"fromBlock" : "0x1",
"address" : "0xd8a25ff31c6174ce7bce74ca4a91c2e816dbf91e"
}],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc": "2.0",
  "result": "0x94f027546179fe922429e172b74bbba0"
}
```

---  

<div id="berith_newBlockFilter"></div>  

### berith_newBlockFilter  
Creates a filter in the node, to notify when a new block arrives.
To check if the state has changed, call [berith_getFilterChanges](#berith_getFilterChanges).

**Parameter**  
none

**Returns**  
`QUANTITY` - A filter id.

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_newBlockFilter","params":[],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Response
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0xf005365f5d4ec0d7cbd28dbe85d0b9c3"
}
```

---  

<div id="berith_newPendingTransactionFilter"></div>  

### berith_newPendingTransactionFilter  
Creates a filter in the node, to notify when new pending transactions arrive.
To check if the state has changed, call [berith_getFilterChanges](#berith_getFilterChanges).

**Parameter**  
none

**Returns**  
`QUANTITY` - A filter id.

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_newPendingTransactionFilter","params":[],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "id":1,
  "jsonrpc":  "2.0",
  "result": "0xe5134f2bd4b107a8932e2d505e03f684"
}
```  

---  
  
<div id="berith_uninstallFilter"></div>  

### berith_uninstallFilter  
Uninstalls a filter with given id. Should always be called when watch is no longer needed.
Additonally Filters timeout when they aren't requested with [berith_getFilterChanges](#berith_getFilterChanges) for a period of time.

**Parameter**  
1. `QUANTITY` - The filter id.  

**Example Parameters**  
```js
params: [
  "0xb" // 11
]
```

**Returns**  
`Boolean` - `true` if the filter was successfully uninstalled, otherwise `false`.

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_uninstallFilter","params":["0xe5134f2bd4b107a8932e2d505e03f684"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": true
}
```  

---  
  
<div id="berith_getFilterChanges"></div>  

### berith_getFilterChanges  
Polling method for a filter, which returns an array of logs which occurred since last poll.

**Parameter**  
1. `QUANTITY` - the filter id.  

**Example Parameters**  
```js
  params: [
    "0x16" // 22
  ]
```

**Returns**  
`Array` - Array of log objects, or an empty array if nothing has changed since last poll.

- For filters created with `berith_newBlockFilter` the return are block hashes (`DATA`, 32 Bytes), e.g. `["0x3454645634534..."]`.
- For filters created with `berith_newPendingTransactionFilter ` the return are transaction hashes (`DATA`, 32 Bytes), e.g. `["0x6345343454645..."]`.
- For filters created with `berith_newFilter` logs are objects with following params:

  - `removed`: `TAG` - `true` when the log was removed, due to a chain reorganization. `false` if its a valid log.
  - `logIndex`: `QUANTITY` - integer of the log index position in the block. `null` when its pending log.
  - `transactionIndex`: `QUANTITY` - integer of the transactions index position log was created from. `null` when its pending log.
  - `transactionHash`: `DATA`, 32 Bytes - hash of the transactions this log was created from. `null` when its pending log.
  - `blockHash`: `DATA`, 32 Bytes - hash of the block where this log was in. `null` when its pending. `null` when its pending log.
  - `blockNumber`: `QUANTITY` - the block number where this log was in. `null` when its pending. `null` when its pending log.
  - `address`: `DATA`, 20 Bytes - address from which this log originated.
  - `data`: `DATA` - contains the non-indexed arguments of the log.
  - `topics`: `Array of DATA` - Array of 0 to 4 32 Bytes `DATA` of indexed log arguments. (In *solidity*: The first topic is the *hash* of the signature of the event (e.g. `Deposit(address,bytes32,uint256)`), except you declared the event with the `anonymous` specifier.)

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_getFilterChanges","params":["0xa7acb4d047ac37e3dd177de95d182a1a"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Result
{
    "jsonrpc":"2.0",
    "id":1,
    "result":[
        "0xb65ce5f9ed6d8b4ee49d1f654f8760010f28ba503ef43816ff93a3b7a5228df4",
        "0x3e44c22fa1960b6dac0fbe37745e5004b60e9404314c504aa9a41292a03af69a",
    ]
}
```

---  
<div id="berith_getFilterLogs"></div>  

### berith_getFilterLogs  
Returns an array of all logs matching filter with given id.

**Parameter**  
1. `QUANTITY` - The filter id.

**Example Parameters**  
```js
params: [
  "0x16" // 22
]
```  

**Returns**  
See [berith_getFilterChanges](#berith_getFilterChanges)

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_getFilterLogs","params":["0x2d584b85ea5ea838335020c16744639"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545
```

---  
<div id="berith_getLogs"></div>  

### berith_getLogs  
Returns an array of all logs matching a given filter object.

**Parameter**  
1. `Object` - The filter options:
  - `fromBlock`: `QUANTITY|TAG` - (optional, default: `"latest"`) Integer block number, or `"latest"` for the last mined block or `"pending"`, `"earliest"` for not yet mined transactions.
  - `toBlock`: `QUANTITY|TAG` - (optional, default: `"latest"`) Integer block number, or `"latest"` for the last mined block or `"pending"`, `"earliest"` for not yet mined transactions.
  - `address`: `DATA|Array`, 20 Bytes - (optional) Contract address or a list of addresses from which logs should originate.
  - `topics`: `Array of DATA`,  - (optional) Array of 32 Bytes `DATA` topics. Topics are order-dependent. Each topic can also be an array of DATA with "or" options.
  - `blockhash`:  `DATA`, 32 Bytes - (optional) With the addition of EIP-234 (Geth >= v1.8.13 or Parity >= v2.1.0), `blockHash` is a new filter option which restricts the logs returned to the single block with the 32-byte hash `blockHash`.  Using `blockHash` is equivalent to `fromBlock` = `toBlock` = the block number with hash `blockHash`.  If `blockHash` is present in the filter criteria, then neither `fromBlock` nor `toBlock` are allowed.

**Example Parameters**  
```js
params: [{
  "topics": ["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b"]
}]
```

**Returns**  
See [berith_getFilterChanges](#berith_getFilterChanges)

**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"berith_getLogs","params":[{"topics":["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b"]}],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545
```

---  

<div id="amon_getBlockCreatorsByNumber"></div>  

### amon_getBlockCreatorsByNumber  
Returns an array of addresses who can seal a block the given block number.

**Parameter**
1. `QUANTITY|TAG`, 32 Bytes - Hash of a block.  


**Returns**  
`Array of DATA`, 20 Bytes - addresses who can seal a block. 


**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"amon_getBlockCreatorsByNumber","params":["0x9"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Response
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": [
    "Bxca7207de79e55c1a69dbc67a4a2e81dfc62c6ac4",
    "Bxd8a25ff31c6174ce7bce74ca4a91c2e816dbf91e",
    "Bx90865e6e6737fe766dd08f39cc2cf1550b5f3875",
    "Bxbb926bbb0b15ca54d4a19dcdf44fc8940e3f6da3",
    "Bx8676fb254279ef78c53b8a781e228ab439065786"
  ]
}
```

---  

<div id="amon_getBlockCreatorsByHash"></div>  

### amon_getBlockCreatorsByHash  
Returns an array of addresses who can seal a block the given block hash.

**Parameter**  
1. `DATA`, 32 Bytes - Hash of a block.


**Returns**    
See [amon_getBlockCreatorsByNumber](#amon_getBlockCreatorsByNumber)


**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"amon_getBlockCreatorsByHash","params":["0x398514d5a403e6245f1af54f96d262f21a8a9bef1fb9b7920e46f14047752cb7"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545
```

---  

<div id="amon_getJoinRatio"></div>  

### amon_getJoinRatio  
Returns a probability of the top of block creators. probability is determined by stake amount.

**Parameter**  
1. `DATA`, 20 Bytes - address to check probability of the top of block creators
2. `QUANTITY|TAG`, 32 Bytes - Hash of a block.


**Returns**    
Returns a probability `float64`


**Example**  
```js
// Request
curl --data '{"jsonrpc":"2.0","method":"amon_getJoinRatio","params":["Bxbb926bbb0b15ca54d4a19dcdf44fc8940e3f6da3", "0x14"],"id":1}' -H "Content-Type: application/json" -X POST localhost:8545

// Response
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": 0.13333333333333333
}
```
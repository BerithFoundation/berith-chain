## Balance and State
Berith has a special balance and transaction model for implementing Stake and adding stability to paying block creation rewards.

### Balance
It is a memory space that stores the balance of coins allocated to each `Balance` account in Berith. It stored in the form of Merkle Patricia Tree in the `State` storage.

#### Type of Balance
Berith has three types of 'Balance' for each account.

![bal_type](./behind2.png)

The above figure shows three kinds of `Balance` structure. There are `MainBalance` and `StakeBalance` wallets exposed through RPC or console commands provided by Berith, and `BehindBalance` which is used in a hidden state inside the node. Below is a description of each balance.

##### MainBalance
`MainBalance` stores the balance of common coins. Coins in `MainBalance` can be transferred to `MainBalance` of another account or your `StakeBalance` at any time.
##### StakeBalance
`StakeBalance` stores the balance of staked coins. Coins in `StakeBalance` can only send the entire quantity to your own `MainBalance`. Depending on the quantity of `StakeBalance`, the probability of creating a block increases.
##### BehindBalance
`BehindBalance` is a temporary repository that stores block generation compensation, waits for as many blocks as `Epoch`, and delivers them to `MainBalance`. The contents of `BehindBalance` cannot be obtained through RPC or console commands. In other words, it is a hidden `Balance`.

The reason for delaying the payment of block generation compensation in Berith is because Berith has a consensus mechanism to allow fork of the chain. Otherwise, if two chains are forked and each creates a separate blockchain, users who are on the disappearing side when one chain is merged into the other will experience a loss of rewards.

![lose_coin](./behind1.png)

The figure above shows the blockchain is forked. If the generation of compensation. If the block creation reward is immediately paid, the user in the picture will experience the loss of the block generation compensation was paid because the chain that the user belonged to disappears. `BehindBalance` is designed to solve this. 

The block generation compensation stored in the `BehindBalance` above is stored as a structure array of block numbers and compensation. Then, when a new block is received, the stored block number is compared with the received block number, and if the difference is more than `Epoch`, it is moved to `MainBalance`.

![move_to_main](./behind3.png)

The above figure shows the process of moving the block generation compensation stored in `BehindBalance` to `MainBalance`.

#### Balance and State
State is a structure that represents information about an account. It is stored in the local DB in the form of a Merkle Patricia Tree. State also contains information about `Balance`.

```
type Account struct {
	Nonce          uint64
	Balance        *big.Int
	Root           common.Hash // merkle root of the storage trie
	CodeHash       []byte
	StakeBalance   *big.Int //brt staking balance
	StakeUpdated   *big.Int //Block number when the stake balance was updated
	Point          *big.Int //selection Point
	BehindBalance  []Behind //behind balance
	Penalty        uint64
	PenlatyUpdated *big.Int //Block Number when the penalty was updated
}
```

The code above is a declaration of a structure that represents an account among the information stored in the state. You can see that `StakeBalance` and `BehindBalance` are declared, including the `Balance` that indicates the `MainBalance` described above.

##### EVM and Transaction Processing

Transactions are handled by the EVM to change `State`. Berith modified the transaction processing of EVM so that the state changes according to he type of transacrion. Berith modified the transaction processing of EVM so that `State` was changed accordingly according to the type of the transaction. Berith modified the transaction processing of the EVM so that the state changes appropriately depending on the type of transaction.

```
func Transfer(db vm.StateDB, sender, recipient common.Address, amount, blockNumber *big.Int, base, target types.JobWallet) {
	/*
		[BERITH]
		Tx 를 state에 적용
	*/
	switch base {
	case types.Main:
		if target == types.Main {
			db.SubBalance(sender, amount)
			db.AddBalance(recipient, amount)
		} else if target == types.Stake {
			//베이스 지갑 차감
			db.SubBalance(sender, amount)
			db.AddStakeBalance(recipient, amount, blockNumber)

		}

		break
	case types.Stake:
		if target == types.Main {
			//스테이크 풀시
			db.RemoveStakeBalance(sender)
		}
		break
	}
}
```
The code above shows the contents of a function that changes `State` by processing a transaction in the EVM. You can see that each different `Balance` is modified according to the type of transaction.


/*
[BERITH]
New function implementation in berith
Where to implement functions used by CLI and RPC
*/

package brtapi

import (
	"bytes"
	"context"
	"errors"
	"strconv"

	"github.com/BerithFoundation/berith-chain/accounts/keystore"
	"github.com/BerithFoundation/berith-chain/miner"
	"github.com/BerithFoundation/berith-chain/rpc"

	"math/big"

	"github.com/BerithFoundation/berith-chain/accounts"
	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/common/hexutil"
	"github.com/BerithFoundation/berith-chain/core/types"
	"github.com/BerithFoundation/berith-chain/crypto"
	"github.com/BerithFoundation/berith-chain/log"
)

//PrivateBerithAPI struct of berith private apis
type PrivateBerithAPI struct {
	backend        Backend
	miner          *miner.Miner
	nonceLock      *AddrLocker
	accountManager *accounts.Manager
}

/*
[BERITH]
SendTxArgs represents the arguments to sumbit a new transaction into the transaction pool.
Specify the tx type by putting Base and Target in the existing transaction structure.
*/
type SendTxArgs struct {
	From     common.Address  `json:"from"`
	To       *common.Address `json:"to"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Value    *hexutil.Big    `json:"value"`
	Nonce    *hexutil.Uint64 `json:"nonce"`
	// We accept "data" and "input" for backwards-compatibility reasons. "input" is the
	// newer name and should be preferred by clients.
	Data   *hexutil.Bytes  `json:"data"`
	Input  *hexutil.Bytes  `json:"input"`
	Base   types.JobWallet `json:"Base"`
	Target types.JobWallet `json:"Target"`
}

// setDefaults is a helper function that fills in default values for unspecified tx fields.
func (args *SendTxArgs) setDefaults(ctx context.Context, b Backend) error {
	if args.Gas == nil {
		args.Gas = new(hexutil.Uint64)
		*(*uint64)(args.Gas) = 90000
	}
	if args.GasPrice == nil {
		price, err := b.SuggestPrice(ctx)
		if err != nil {
			return err
		}
		args.GasPrice = (*hexutil.Big)(price)
	}
	if args.Value == nil {
		args.Value = new(hexutil.Big)
	}
	if args.Nonce == nil {
		nonce, err := b.GetPoolNonce(ctx, args.From)
		if err != nil {
			return err
		}
		args.Nonce = (*hexutil.Uint64)(&nonce)
	}
	if args.Data != nil && args.Input != nil && !bytes.Equal(*args.Data, *args.Input) {
		return errors.New(`Both "data" and "input" are set and not equal. Please use "input" to pass transaction call data.`)
	}
	if args.To == nil {
		// Contract creation
		var input []byte
		if args.Data != nil {
			input = *args.Data
		} else if args.Input != nil {
			input = *args.Input
		}
		if len(input) == 0 {
			return errors.New(`contract creation without any data provided`)
		}
	}
	return nil
}

func (args *SendTxArgs) toTransaction() *types.Transaction {
	var input []byte
	if args.Data != nil {
		input = *args.Data
	} else if args.Input != nil {
		input = *args.Input
	}
	if args.To == nil {
		return types.NewContractCreation(uint64(*args.Nonce), (*big.Int)(args.Value), uint64(*args.Gas), (*big.Int)(args.GasPrice), input, args.Base, args.Target)
	}
	return types.NewTransaction(uint64(*args.Nonce), *args.To, (*big.Int)(args.Value), uint64(*args.Gas), (*big.Int)(args.GasPrice), input, args.Base, args.Target)
}

/*
[BERITH]
The first function called to register the implementation
*/
func NewPrivateBerithAPI(b Backend, m *miner.Miner, nonceLock *AddrLocker) *PrivateBerithAPI {
	return &PrivateBerithAPI{
		backend:        b,
		miner:          m,
		accountManager: b.AccountManager(),
		nonceLock:      nonceLock,
	}
}

// submitTransaction is a helper function that submits tx to txPool and logs a message.
func submitTransaction(ctx context.Context, b Backend, tx *types.Transaction) (common.Hash, error) {
	if err := b.SendTx(ctx, tx); err != nil {
		return common.Hash{}, err
	}
	if tx.To() == nil {
		signer := types.MakeSigner(b.ChainConfig(), b.CurrentBlock().Number())
		from, err := types.Sender(signer, tx)
		if err != nil {
			return common.Hash{}, err
		}
		addr := crypto.CreateAddress(from, tx.Nonce())
		log.Info("Submitted contract creation", "fullhash", tx.Hash().Hex(), "contract", addr.Hex())
	} else {
		log.Info("Submitted transaction", "fullhash", tx.Hash().Hex(), "recipient", tx.To())
	}
	return tx.Hash(), nil
}

type WalletTxArgs struct {
	From     common.Address  `json:"from"`
	Value    *hexutil.Big    `json:"value"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Nonce    *hexutil.Uint64 `json:"nonce"`
}

/*
[BERITH]
Function to check the elected point of the specified account
*/
func (s *PrivateBerithAPI) GetSelectionPoint(ctx context.Context, address common.Address, blockNr rpc.BlockNumber) (*hexutil.Big, error) {
	state, _, err := s.backend.StateAndHeaderByNumber(ctx, blockNr)
	if state == nil || err != nil {
		return nil, err
	}

	return (*hexutil.Big)(state.GetPoint(address)), state.Error()
}

/*
[BERITH]
- SendStaking creates a transaction for user staking
- Function to handle berith.stake request
- Create Tx with added base and target.
- Includes pre-processing logic that returns less than 100,000 errors when staking early
- WalletTxArg structure is a structure to limit Tx
*/
func (s *PrivateBerithAPI) Stake(ctx context.Context, args WalletTxArgs) (common.Hash, error) {

	state, _, err := s.backend.StateAndHeaderByNumber(ctx, rpc.LatestBlockNumber)
	if state == nil || err != nil {
		return common.Hash{}, err
	}
	stakedAmount := state.GetStakeBalance(args.From)
	stakingAmount := args.Value.ToInt()
	totalStakingAmount := new(big.Int).Add(stakingAmount, stakedAmount)

	if config := s.backend.ChainConfig(); config.IsEIP155(s.backend.CurrentBlock().Number()) {
		if totalStakingAmount.Cmp(config.Bsrr.StakeMinimum) <= -1 {
			minimum := new(big.Int).Div(config.Bsrr.StakeMinimum, common.UnitForBer)

			log.Error("The minimum number of stakes is " + strconv.Itoa(int(minimum.Uint64())))
			return common.Hash{}, errors.New("staking balance failed")
		}
	}

	// Look up the wallet containing the requested signer
	sendTx := new(SendTxArgs)

	// Create transaction
	sendTx.From = args.From
	sendTx.To = &args.From
	sendTx.Value = args.Value
	sendTx.Base = types.Main
	sendTx.Target = types.Stake
	sendTx.Gas = args.Gas
	sendTx.GasPrice = args.GasPrice
	sendTx.Nonce = args.Nonce

	return s.sendTransaction(ctx, *sendTx)
}

/*
[BERITH]
- SendStaking creates a transaction for user staking
- When this function is called, all staking is released and returned to Main
- After creating Tx and sending it, it is processed by Consensus.
*/
func (s *PrivateBerithAPI) StopStaking(ctx context.Context, args WalletTxArgs) (common.Hash, error) {
	// Look up the wallet containing the requested signer
	sendTx := new(SendTxArgs)

	sendTx.From = args.From
	sendTx.To = &args.From
	sendTx.Value = new(hexutil.Big)
	sendTx.Base = types.Stake
	sendTx.Target = types.Main
	sendTx.Gas = args.Gas
	sendTx.GasPrice = args.GasPrice

	return s.sendTransaction(ctx, *sendTx)
}

/*
[BERITH]
- private trasaction function
- Functions that deal with actual transactions
*/
func (s *PrivateBerithAPI) sendTransaction(ctx context.Context, args SendTxArgs) (common.Hash, error) {
	account := accounts.Account{Address: args.From}

	wallet, err := s.backend.AccountManager().Find(account)
	if err != nil {
		return common.Hash{}, err
	}

	s.nonceLock.LockAddr(args.From)
	defer s.nonceLock.UnlockAddr(args.From)

	// Set some sanity defaults and terminate on failure
	if err := args.setDefaults(ctx, s.backend); err != nil {
		return common.Hash{}, err
	}
	// Assemble the transaction and sign with the wallet
	tx := args.toTransaction()

	var chainID *big.Int
	if config := s.backend.ChainConfig(); config.IsEIP155(s.backend.CurrentBlock().Number()) {
		chainID = config.ChainID
	}
	signed, err := wallet.SignTx(account, tx, chainID)
	if err != nil {
		return common.Hash{}, err
	}

	return submitTransaction(ctx, s.backend, signed)
}

/*
[BERITH]
 - Function to check the staking quantity of the specified Account
 - Check and return the current local block status
*/
func (s *PrivateBerithAPI) GetStakeBalance(ctx context.Context, address common.Address, blockNr rpc.BlockNumber) (*hexutil.Big, error) {
	state, _, err := s.backend.StateAndHeaderByNumber(ctx, blockNr)
	if state == nil || err != nil {
		return nil, err
	}
	return (*hexutil.Big)(state.GetStakeBalance(address)), state.Error()
}

/*
[BERITH]
 - Structure for returning account information
*/
type AccountInfo struct {
	Balance      *big.Int //main balance
	StakeBalance *big.Int //staking balance
}

/*
[BERITH]
- Function to return account information (All Balance)
- Functions created for convenience of information verification
*/
func (s *PrivateBerithAPI) GetAccountInfo(ctx context.Context, address common.Address, blockNr rpc.BlockNumber) (*AccountInfo, error) {
	state, _, err := s.backend.StateAndHeaderByNumber(ctx, blockNr)
	if state == nil || err != nil {
		return nil, err
	}

	account := state.GetAccountInfo(address)
	info := &AccountInfo{
		Balance:      account.Balance,
		StakeBalance: account.StakeBalance,
	}

	return info, state.Error()
}

/*
[BERITH]
- Function to change the keystore password
- Made by user's request
*/
func (s *PrivateBerithAPI) UpdateAccount(ctx context.Context, address common.Address, passphrase, newPassphrase string) error {
	return fetchKeystore(s.accountManager).Update(accounts.Account{Address: address}, passphrase, newPassphrase)
}

// fetchKeystore retrives the encrypted keystore from the account manager.
func fetchKeystore(am *accounts.Manager) *keystore.KeyStore {
	return am.Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
}

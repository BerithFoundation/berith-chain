// [BERITH]
// New function implementation in berith
// Where to implement functions used by CLI and RPC

package brtapi

import (
	"berith-chain/core"
	"context"
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

// PrivateBerithAPI struct of berith private apis
type PrivateBerithAPI struct {
	backend        Backend
	miner          *miner.Miner
	nonceLock      *AddrLocker
	accountManager *accounts.Manager
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

/*
[BERITH]
Function to get the SelectionPoint of the specified account
*/
func (s *PrivateBerithAPI) GetSelectionPoint(ctx context.Context, address common.Address, blockNumber rpc.BlockNumber) (*hexutil.Big, error) {
	state, _, err := s.backend.StateAndHeaderByNumber(ctx, blockNumber)
	if state == nil || err != nil {
		return nil, err
	}

	return (*hexutil.Big)(state.GetPoint(address)), state.Error()
}

/*
[BERITH]
Stake creates a transaction for user staking
Function to handle berith.stake request
Create Tx with added base and target.
WalletTxArgs structure is a structure to limit Tx
*/
func (s *PrivateBerithAPI) Stake(ctx context.Context, wallet WalletTxArgs) (common.Hash, error) {
	state, _, err := s.backend.StateAndHeaderByNumber(ctx, rpc.LatestBlockNumber)
	if state == nil || err != nil {
		return common.Hash{}, err
	}

	stakedAmount := state.GetStakeBalance(wallet.From)
	stakingAmount := wallet.Value.ToInt()
	totalStakingAmount := new(big.Int).Add(stakingAmount, stakedAmount)

	// 본래 Stake 트랜잭션 검증을 IsEIP155 포크 기준으로 하기 때문에 별도의 Berith 포크 설정 또한 필요하지 않음.
	if config := s.backend.ChainConfig(); config.IsEIP155(s.backend.CurrentBlock().Number()) {
		err := checkStakeMinimum(totalStakingAmount, config.Bsrr.StakeMinimum)
		if err != nil {
			return common.Hash{}, err
		}
	}

	// Create transaction
	sendTx := &SendTxArgs{
		From:     wallet.From,
		To:       &wallet.From,
		Value:    wallet.Value,
		Base:     types.Main,
		Target:   types.Stake,
		Gas:      wallet.Gas,
		GasPrice: wallet.GasPrice,
		Nonce:    wallet.Nonce,
	}
	return s.sendTransaction(ctx, *sendTx)
}

func checkStakeMinimum(stakeAmount *big.Int, stakeMininum *big.Int) error {
	if stakeAmount.Cmp(stakeMininum) <= -1 {
		minimum := new(big.Int).Div(stakeMininum, common.UnitForBer)

		log.Error("The mininum number of stakes is " + strconv.Itoa(int(minimum.Uint64())))
		return core.ErrUnderStakeBalance
	}
	return nil
}

/*
[BERITH]
When this function is called, all staking is released and returned to Main
After creating Tx and sending it, it is processed by Consensus.
*/
func (s *PrivateBerithAPI) StopStaking(ctx context.Context, wallet WalletTxArgs) (common.Hash, error) {
	sendTx := &SendTxArgs{
		From:     wallet.From,
		To:       &wallet.From,
		Value:    new(hexutil.Big),
		Base:     types.Stake,
		Target:   types.Main,
		Gas:      wallet.Gas,
		GasPrice: wallet.GasPrice,
	}
	return s.sendTransaction(ctx, *sendTx)
}

/*
[BERITH]
Functions that deal with actual transactions
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

// submitTransaction is a helper function that submits tx to txPool and logs a message.
func submitTransaction(ctx context.Context, b Backend, tx *types.Transaction) (common.Hash, error) {
	if err := b.SendTx(ctx, tx); err != nil {
		return common.Hash{}, err
	}

	if err := printTxLog(b, tx); err == nil {
		return common.Hash{}, err
	}
	return tx.Hash(), nil
}

// printTxLog is a function that provides a log of transmitted transactions.
func printTxLog(b Backend, tx *types.Transaction) error {
	if tx.To() == nil {
		signer := types.MakeSigner(b.ChainConfig(), b.CurrentBlock().Number())
		from, err := types.Sender(signer, tx)
		if err != nil {
			return err
		}
		addr := crypto.CreateAddress(from, tx.Nonce())
		log.Info("Submitted contract creation", "fullhash", tx.Hash().Hex(), "contract", addr.Hex())
	} else {
		log.Info("Submitted transaction", "fullhash", tx.Hash().Hex(), "recipient", tx.To())
	}
	return nil
}

/*
[BERITH]
Function to check the staking quantity of the specified Account
Check and return the current local block status
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
Function to return account information (All Balance)
Functions created for convenience of information verification
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
Function to change the keystore password
Made by user's request
*/
func (s *PrivateBerithAPI) UpdateAccount(ctx context.Context, address common.Address, passphrase, newPassphrase string) error {
	return fetchKeystore(s.accountManager).Update(accounts.Account{Address: address}, passphrase, newPassphrase)
}

// fetchKeystore retrives the encrypted keystore from the account manager.
func fetchKeystore(am *accounts.Manager) *keystore.KeyStore {
	return am.Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
}

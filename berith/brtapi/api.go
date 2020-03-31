/*
[BERITH]
berith 에서 새롭게 추가된 함수 구현체
CLI 및 RPC 에서 사용하는 함수를 구현하는 곳
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

	// "fmt"
	"math/big"

	"github.com/BerithFoundation/berith-chain/accounts"
	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/common/hexutil"
	"github.com/BerithFoundation/berith-chain/core/types"
	"github.com/BerithFoundation/berith-chain/crypto"
	"github.com/BerithFoundation/berith-chain/log"
	// "github.com/BerithFoundation/berith-chain/berith/stake"
)

//PrivateBerithAPI struct of berith private apis
type PrivateBerithAPI struct {
	backend        Backend
	miner          *miner.Miner
	nonceLock      *AddrLocker
	accountManager *accounts.Manager
}

// SendTxArgs represents the arguments to sumbit a new transaction into the transaction pool.
/*
[BERITH]
기존 트렌젝션 구조체에 base, target 을 넣어 tx 타입을 지정한다.
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
	base   types.JobWallet `json:"base"`
	target types.JobWallet `json:"target"`
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
		return types.NewContractCreation(uint64(*args.Nonce), (*big.Int)(args.Value), uint64(*args.Gas), (*big.Int)(args.GasPrice), input, args.base, args.target)
	}
	return types.NewTransaction(uint64(*args.Nonce), *args.To, (*big.Int)(args.Value), uint64(*args.Gas), (*big.Int)(args.GasPrice), input, args.base, args.target)
}

/*
[BERITH]
구현체를 등록 하기 위해 최초로 호출되는 함수
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
지정된 어카운트의 선출 포인트를 확인 할수 있는 함수
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
- berith.stake 명령시 처리 하는 함수로 추가된 base 와 target 을 지정하여 Tx 를 만드는 함수
- 초반 스테이킹시 10만개이하로 에러를 반환하는 선처리 로직 포함
- WalletTxArg 구조체는 Tx 를 제한두기위한 구조체
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

	//Tx 를 만들어줌
	sendTx.From = args.From
	sendTx.To = &args.From
	sendTx.Value = args.Value
	sendTx.base = types.Main
	sendTx.target = types.Stake
	sendTx.Gas = args.Gas
	sendTx.GasPrice = args.GasPrice
	sendTx.Nonce = args.Nonce

	return s.sendTransaction(ctx, *sendTx)
}

/*
[BERITH]
- SendStaking creates a transaction for user staking
- 이함수를 호출하면 모든 Staking 해제 하고 Main 으로 반환됨
- Tx 를 만들어서 Send 하는 역할만 함 이후 처리는 Consensus 에서 처리
*/
func (s *PrivateBerithAPI) StopStaking(ctx context.Context, args WalletTxArgs) (common.Hash, error) {
	// Look up the wallet containing the requested signer
	sendTx := new(SendTxArgs)

	sendTx.From = args.From
	sendTx.To = &args.From
	sendTx.Value = new(hexutil.Big)
	sendTx.base = types.Stake
	sendTx.target = types.Main
	sendTx.Gas = args.Gas
	sendTx.GasPrice = args.GasPrice

	return s.sendTransaction(ctx, *sendTx)
}

/*
[BERITH]
- private trasaction function
- 실제 트렌젝션 처리를 하는 함수
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
 - 지정한 Account 의 스테이킹 수량을 확인 하는 함수
 - 현재 로컬상 블록 상태를 확인 하여 반환
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
 - 어카운트 정보를 반환 하기 위한 구조체
*/
type AccountInfo struct {
	Balance      *big.Int //main balance
	StakeBalance *big.Int //staking balance
}

/*
[BERITH]
- 어카운트 정보 (All Balance) 를 반환 하기 위한 함수
- 정보 확인 편의를 생각하여 만든 함수
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
- 키스토어 의 비밀번호를 변경하는 함수
- 유저들의 요청에 의해 만듬
*/
func (s *PrivateBerithAPI) UpdateAccount(ctx context.Context, address common.Address, passphrase, newPassphrase string) error {
	return fetchKeystore(s.accountManager).Update(accounts.Account{Address: address}, passphrase, newPassphrase)
}

// fetchKeystore retrives the encrypted keystore from the account manager.
func fetchKeystore(am *accounts.Manager) *keystore.KeyStore {
	return am.Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
}

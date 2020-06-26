package pkg

import (
	"fmt"
	"sync"

	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank"

	keys "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	rpchttp "github.com/tendermint/tendermint/rpc/client"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

type TxManager struct {
	conf      Configuration
	node      rpchttp.ABCIClient
	kb        keys.Keybase
	faucetAcc *auth.BaseAccount
	mux       sync.Mutex
}

func (tm *TxManager) queryWithData(path string, data []byte) ([]byte, int64, error) {
	res, err := tm.node.ABCIQuery(path, data)
	if err != nil {
		return nil, 0, err
	}
	return res.Response.Value, res.Response.Height, nil
}

func NewTxManager(conf Configuration, node rpchttp.ABCIClient) *TxManager {
	return &TxManager{node: node, conf: conf}
}

func (tm *TxManager) WithKeybase(kb keys.Keybase) *TxManager {
	tm.kb = kb
	return tm
}

func (tm *TxManager) Init() error {
	info, err := tm.kb.Get("faucet")
	if err != nil {
		return err
	}
	// fetch account info
	// we need these step to fetch account sequence on live chain
	acc, err := tm.fetchAccount(info.GetAddress())
	if err != nil {
		return err
	}
	tm.faucetAcc = acc
	return nil
}

func (tm *TxManager) fetchAccount(addr sdk.AccAddress) (*auth.BaseAccount, error) {
	path := fmt.Sprintf("custom/%s/%s", auth.QuerierRoute, auth.QueryAccount)
	bs, err := ModuleCdc.MarshalJSON(types.NewQueryAccountParams(addr))
	if err != nil {
		return nil, errors.Wrap(err, "codec marshalling failed")
	}
	result, err := tm.node.ABCIQuery(path, bs)
	if err != nil {
		return nil, errors.Wrap(err, "abci query failed")
	}
	resp := result.Response
	if !resp.IsOK() {
		return nil, errors.Wrapf(errors.ErrInvalidRequest, "abci response: %s %s", addr.String(), resp.Log)
	}

	var baseAcc auth.BaseAccount
	if err := ModuleCdc.UnmarshalJSON(resp.Value, &baseAcc); err != nil {
		return nil, errors.Wrap(err, "fetch account failed")
	}
	return &baseAcc, nil
}

func (tm *TxManager) BroadcastTx(tx []byte) (*coretypes.ResultBroadcastTx, error) {
	return tm.node.BroadcastTxSync(tx)
}

func (tm *TxManager) BuildAndSignTx(targetAcc sdk.AccAddress) ([]byte, error) {
	/* CONTRACT
	a faucet wallet must be used by single actor otherwise successful tx will bump
	account sequence on chain.
	*/
	//
	tm.mux.Lock()
	seq := tm.faucetAcc.GetSequence()
	err := tm.faucetAcc.SetSequence(seq + 1)
	tm.mux.Unlock()
	if err != nil {
		return nil, err
	}

	txBuilder := auth.TxBuilder{}.
		WithTxEncoder(auth.DefaultTxEncoder(ModuleCdc)).
		WithAccountNumber(tm.faucetAcc.AccountNumber).
		WithSequence(seq).
		WithGasPrices(tm.conf.GasPrices).
		WithChainID(tm.conf.ChainID).
		WithMemo(tm.conf.Memo).WithKeybase(tm.kb)

	sendMsg := bank.MsgSend{
		FromAddress: tm.faucetAcc.GetAddress(),
		ToAddress:   targetAcc,
		Amount: sdk.Coins{
			sdk.NewInt64Coin(tm.conf.CoinDenom, tm.conf.SendAmount),
		},
	}

	// adjust gas
	simTx, err := txBuilder.BuildTxForSim([]sdk.Msg{sendMsg})
	if err != nil {
		return nil, errors.Wrap(err, "tx gas adjustment failed")
	}
	_, adjusted, err := utils.CalculateGas(tm.queryWithData, ModuleCdc, simTx, tm.conf.GasAdjust)

	txBuilder = txBuilder.WithGas(adjusted)
	tx, err := txBuilder.BuildAndSign("faucet", tm.conf.Passphrase, []sdk.Msg{sendMsg})
	if err != nil {
		return nil, errors.Wrap(err, "tx signing failed")
	}

	return tx, nil
}

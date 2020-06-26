package pkg

/*
import (
	"log"
	"os"
	"testing"

	keys2 "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client"
)

var (
	conf Configuration
	node rpchttp.Client
)

const (
	armor = ""
	pass       = ""
	address    = ""
	targetAddr = ""
	tendermintRPC = "http://localhost:26657"
	coindenom = "tiov"
	send = 10
	chainID = "local"
)

func TestMain(m *testing.M) {
	conf = Configuration{
		TendermintRPC: tendermintRPC,
		CoinDenom:     coindenom,
		ChainID:       chainID,
		Passphrase:    pass,
		SendAmount:    send,
	}
	var err error
	node, err = rpchttp.NewHTTP(conf.TendermintRPC, "/websocket")
	if err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	os.Exit(code)
}

func Test_TxManagerFetchAccount(t *testing.T) {
	mn := NewTxManager(conf, node)
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		t.Fatal(err)
	}
	acc, err := mn.fetchAccount(addr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(acc)
}

func Test_TxManagerBuildAndSignTx(t *testing.T) {
	mn := NewTxManager(conf, node)
	kb := keys2.NewInMemory()
	if err := kb.ImportPrivKey("faucet", armor, pass); err != nil {
		t.Fatal(err)
	}
	mn.WithKeybase(kb)
	if err := mn.Init(); err != nil {
		t.Fatalf("tx manager: %v", err)
	}
	addr, err := sdk.AccAddressFromBech32(targetAddr)
	if err != nil {
		t.Fatal(err)
	}
	_, err = mn.BuildAndSignTx(addr)
	if err != nil {
		t.Error(err)
	}
}

func Test_TxManagerBroadcastTx(t *testing.T) {
	mn := NewTxManager(conf, node)
	kb := keys2.NewInMemory()
	if err := kb.ImportPrivKey("faucet", armor, pass); err != nil {
		t.Fatal(err)
	}
	mn.WithKeybase(kb)
	if err := mn.Init(); err != nil {
		t.Fatalf("tx manager: %v", err)
	}
	addr, err := sdk.AccAddressFromBech32(targetAddr)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := mn.BuildAndSignTx(addr)
	if err != nil {
		t.Error(err)
	}
	res, err := mn.BroadcastTx(tx)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}
*/

package signutil

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/iov-one/iovns/app/config"
	cryptoamino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"testing"
)

func TestVerify(t *testing.T) {
	config.ApplyChangesAndSeal(sdk.GetConfig())
	sdk.RegisterCodec(ModuleCdc)
	auth.RegisterCodec(ModuleCdc)
	cryptoamino.RegisterAmino(ModuleCdc)
	const testSig = `{"type":"cosmos-sdk/StdTx","value":{"msg":[{"type":"signutil/MsgSignText","value":{"message":"hello","signer":"star1ynqxwk8gcmkfg7e30p6uumnx0mrfcmzrjmfnap"}}],"fee":{"amount":[],"gas":"200000"},"signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"Akxp+TXfAnYFJIRRxjWA3m56mK+plzrECk6kG7opJ02V"},"signature":"TkAAEht40YZQkbriqzcY2IswVijoBpFfOKCF0e3CEed/M6nEAAJ9VSj18c+f9QNocQXWwjMT/fpRNCu70x9Q1Q=="}],"memo":""}}`
	var tx auth.StdTx
	ModuleCdc.MustUnmarshalJSON([]byte(testSig), &tx)
	err := Verify(tx, DefaultChainID, DefaultAccountNumber, DefaultSequence)
	if err != nil {
		t.Fatal(err)
	}
}

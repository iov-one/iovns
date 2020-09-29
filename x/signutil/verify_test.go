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
	const testSig = `{"type":"cosmos-sdk/StdTx","value":{"msg":[{"type":"signutil/MsgSignText","value":{"message":"hello","signer":"star1zrwgm6skw3j6e2tjgq4vj5u5avmzvr6ed29vsl"}}],"fee":{"amount":[],"gas":"200000"},"signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"A4NhkRqUan3iC7iCW/mArMUw0hfHw/BjWV5zCiABINXk"},"signature":"o/eXmmHnuYs/waNdjr4dx4Lh/yCWtyYjcIFL0NUKYCZDDugj9cmB9++ZHRj0i8HmYyZCbsHMOUGheyxCOCMpug=="}],"memo":""}}`
	var tx auth.StdTx
	ModuleCdc.MustUnmarshalJSON([]byte(testSig), &tx)
	err := Verify(tx, "test", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
}

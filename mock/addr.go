package mock

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"os"
)

func Addresses() (sdk.AccAddress, sdk.AccAddress) {
	keyBase := keys.NewInMemory()
	addr1, _, err := keyBase.CreateMnemonic("alice", keys.English, "", keys.Secp256k1)
	if err != nil {
		fmt.Println("unable to generate mock addresses " + err.Error())
		os.Exit(1)
	}
	addr2, _, err := keyBase.CreateMnemonic("bob", keys.English, "", keys.Secp256k1)
	if err != nil {
		fmt.Println("unable to generate mock addresses " + err.Error())
		os.Exit(1)
	}
	return addr1.GetAddress(), addr2.GetAddress()
}

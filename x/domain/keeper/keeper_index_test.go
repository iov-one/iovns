package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/types"
	"os"
	"testing"
)

func genAddress() (sdk.AccAddress, sdk.AccAddress) {
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

func Test_indexFunctionality(t *testing.T) {
	k, ctx := NewTestKeeper(t, true)
	k.CreateAccount(ctx, types.Account{
		Domain: "test",
		Name:   "1",
		Owner:  aliceAddr,
	})
	k.CreateAccount(ctx, types.Account{
		Domain: "test",
		Name:   "2",
		Owner:  bobAddr,
	})
	k.CreateAccount(ctx, types.Account{
		Domain: "test",
		Name:   "3",
		Owner:  aliceAddr,
	})
	accountKeys := k.iterAccountToOwner(ctx, aliceAddr)
	// expected two keys
	if len(accountKeys) != 2 {
		t.Fatalf("expected two keys, got: %d", len(accountKeys))
	}
	// transfer account
	acc, _ := k.GetAccount(ctx, "test", "1")
	k.TransferAccount(ctx, acc, bobAddr)
	// expected two keys for account bobAddr
	if len(k.iterAccountToOwner(ctx, bobAddr)) != 2 {
		t.Fatalf("expected two keys for %s, got: %d", bobAddr, len(k.iterAccountToOwner(ctx, bobAddr)))
	}
	// expect one key for aliceAddr
	if len(k.iterAccountToOwner(ctx, aliceAddr)) != 1 {
		t.Fatalf("expected two keys for %s, got: %d", bobAddr, len(k.iterAccountToOwner(ctx, aliceAddr)))
	}
	// delete account from bobAddr
	k.DeleteAccount(ctx, "test", "1") // belongs to bobAddr
	if len(k.iterAccountToOwner(ctx, bobAddr)) != 1 {
		t.Fatalf("expected two keys for %s, got: %d", bobAddr, len(k.iterAccountToOwner(ctx, bobAddr)))
	}

}

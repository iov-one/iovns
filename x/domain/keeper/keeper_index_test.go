package keeper

import (
	"github.com/iov-one/iovns/x/domain/types"
	"testing"
)

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

func TestKeeper_iterAccountToOwner(t *testing.T) {

}

func TestKeeper_iterDomainToOwner(t *testing.T) {

}

func TestKeeper_mapAccountToOwner(t *testing.T) {

}

func TestKeeper_mapDomainToOwner(t *testing.T) {
}

func TestKeeper_unmapAccountToOwner(t *testing.T) {

}

func TestKeeper_unmapDomainToOwner(t *testing.T) {

}

func Test_accAddrFromIndex(t *testing.T) {
	if !(aliceAddr.String() == accAddrFromIndex(indexAddr(aliceAddr)).String()) {
		t.Fatalf("mismatched addresses for: %s", aliceAddr.String())
	}
}

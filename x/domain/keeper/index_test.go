package keeper

import (
	"fmt"
	"github.com/iov-one/iovns/x/domain/types"
	"testing"
)

func Test_accountIndexing(t *testing.T) {
	var accountKeys [][]byte
	do := func(key []byte) bool {
		accountKeys = append(accountKeys, key)
		return true
	}
	k, ctx, _ := NewTestKeeper(t, true)
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
	k.iterAccountToOwner(ctx, aliceAddr, do)
	// expected two keys
	if len(accountKeys) != 2 {
		t.Fatalf("expected two keys, got: %d", len(accountKeys))
	}
	accountKeys = nil
	// transfer account
	acc, _ := k.GetAccount(ctx, "test", "1")
	k.TransferAccount(ctx, acc, bobAddr)
	// expected two keys for account bobAddr
	k.iterAccountToOwner(ctx, bobAddr, do)
	if len(accountKeys) != 2 {
		t.Fatalf("expected two keys for %s, got: %d", bobAddr, len(accountKeys))
	}
	accountKeys = nil
	// expect one key for aliceAddr
	k.iterAccountToOwner(ctx, aliceAddr, do)
	if len(accountKeys) != 1 {
		t.Fatalf("expected two keys for %s, got: %d", bobAddr, len(accountKeys))
	}
	accountKeys = nil
	// delete account from bobAddr
	k.DeleteAccount(ctx, "test", "1") // belongs to bobAddr
	k.iterAccountToOwner(ctx, bobAddr, do)
	if len(accountKeys) != 1 {
		t.Fatalf("expected two keys for %s, got: %d", bobAddr, len(accountKeys))
	}

}

func Test_domainIndexing(t *testing.T) {
	var domainKeys [][]byte
	do := func(key []byte) bool {
		domainKeys = append(domainKeys, key)
		return true
	}
	k, ctx, _ := NewTestKeeper(t, true)
	k.CreateDomain(ctx, types.Domain{
		Name:  "1",
		Admin: bobAddr,
	})
	k.CreateDomain(ctx, types.Domain{
		Name:  "2",
		Admin: aliceAddr,
	})
	// check number of keys mapped to owner
	k.iterDomainToOwner(ctx, bobAddr, do)
	if l := len(domainKeys); l != 1 {
		t.Fatalf("expected %d keys got: %d", 1, l)
	}
	domainKeys = nil
	// transfer domain
	domain, _ := k.GetDomain(ctx, "1")
	k.TransferDomain(ctx, aliceAddr, domain)
	// check if addr b has 0 keys
	k.iterDomainToOwner(ctx, bobAddr, do)
	if l := len(domainKeys); l != 0 {
		t.Fatalf("expected %d keys got: %d", 0, l)
	}
	domainKeys = nil
	// check if addr a has 2 keys
	k.iterDomainToOwner(ctx, aliceAddr, do)
	if l := len(domainKeys); l != 2 {
		t.Fatalf("expected %d keys got: %d", 2, l)
	}
	domainKeys = nil
	// delete domain
	_ = k.DeleteDomain(ctx, "2")
	// check if addr b has 1 key
	k.iterDomainToOwner(ctx, aliceAddr, do)
	if l := len(domainKeys); l != 1 {
		t.Fatalf("expected %d keys got: %d", 1, l)
	}
}

// checks if the functions that convert address to indexed address and indexed address to address
// are reversible and compatible
func Test_addressIndexing(t *testing.T) {
	if !(aliceAddr.String() == accAddrFromIndex(indexAddr(aliceAddr)).String()) {
		t.Fatalf("mismatched addresses for: %s", aliceAddr.String())
	}
}

func Test_targetsIndexing(t *testing.T) {
	accMatch := func(acc1, acc2 types.Account) error {
		if acc1.Name != acc2.Name {
			return fmt.Errorf("name mismatch: %s <-> %s", acc1.Name, acc2.Name)
		}
		if acc1.Domain != acc2.Domain {
			return fmt.Errorf("domain mismatch: %s<-> %s", acc1.Domain, acc2.Domain)
		}
		return nil
	}
	var accountKeys [][]byte
	do := func(key []byte) bool {
		accountKeys = append(accountKeys, key)
		return true
	}
	k, ctx, _ := NewTestKeeper(t, true)
	// create one target
	targetA := types.BlockchainAddress{
		ID:      "t1",
		Address: "1",
	}
	targetB := types.BlockchainAddress{
		ID:      "t2",
		Address: "2",
	}
	// create one account
	accountA := types.Account{
		Domain: "test",
		Name:   "1",
		Targets: []types.BlockchainAddress{
			targetA,
			targetB,
		},
		Owner: aliceAddr,
	}
	// insert account
	k.CreateAccount(ctx, accountA)
	// iterate targets
	k.iterateBlockchainTargetsAccounts(ctx, targetA, do)
	if len(accountKeys) != 1 {
		t.Fatalf("expected 1 keys, got: %d", len(accountKeys))
	}
	// generate test account
	acc := &types.Account{}
	err := acc.Unpack(accountKeys[0])
	if err != nil {
		t.Fatalf("unpack error: %d", err)
	}
	// check if it matches
	if err := accMatch(*acc, accountA); err != nil {
		t.Fatal(err)
	}
	// DeleteAccount
	accountKeys = nil
	k.DeleteAccount(ctx, accountA.Domain, accountA.Name)
	k.iterateBlockchainTargetsAccounts(ctx, targetA, do)
	if len(accountKeys) != 0 {
		t.Fatalf("no key expected, got: %d", len(accountKeys))
	}
	// ReplaceAccountTargets
	accountKeys = nil
	k.CreateAccount(ctx, accountA)
	k.ReplaceAccountTargets(ctx, accountA, []types.BlockchainAddress{targetB})
	k.iterateBlockchainTargetsAccounts(ctx, targetA, do)
	if len(accountKeys) != 0 {
		t.Fatalf("no key expected, got: %d", len(accountKeys))
	}
	accountKeys = nil
	k.iterateBlockchainTargetsAccounts(ctx, targetB, do)
	if len(accountKeys) != 1 {
		t.Fatalf("expected 1 key, got: %d", len(accountKeys))
	}
	if err := accMatch(*acc, accountA); err != nil {
		t.Fatal(err)
	}
	// TransferAccount
	accountKeys = nil
	accountA, _ = k.GetAccount(ctx, accountA.Domain, accountA.Name) // edited the account before, so update it
	k.TransferAccount(ctx, accountA, bobAddr)
	// check if targetA is associated with any account
	k.iterateBlockchainTargetsAccounts(ctx, targetA, do)
	if len(accountKeys) != 0 {
		t.Fatalf("expected 0 keys, got: %d", len(accountKeys))
	}
}

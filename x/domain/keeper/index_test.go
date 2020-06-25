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
	err := k.iterAccountToOwner(ctx, aliceAddr, do)
	if err != nil {
		t.Fatal(err)
	}
	// expected two keys
	if len(accountKeys) != 2 {
		t.Fatalf("expected two keys, got: %d", len(accountKeys))
	}
	accountKeys = nil
	// transfer account
	acc, _ := k.GetAccount(ctx, "test", "1")
	k.TransferAccount(ctx, acc, bobAddr)
	// expected two keys for account bobAddr
	err = k.iterAccountToOwner(ctx, bobAddr, do)
	if err != nil {
		t.Fatal(err)
	}
	if len(accountKeys) != 2 {
		t.Fatalf("expected two keys for %s, got: %d", bobAddr, len(accountKeys))
	}
	accountKeys = nil
	// expect one key for aliceAddr
	err = k.iterAccountToOwner(ctx, aliceAddr, do)
	if err != nil {
		t.Fatal(err)
	}
	if len(accountKeys) != 1 {
		t.Fatalf("expected two keys for %s, got: %d", bobAddr, len(accountKeys))
	}
	accountKeys = nil
	// delete account from bobAddr
	k.DeleteAccount(ctx, "test", "1") // belongs to bobAddr
	err = k.iterAccountToOwner(ctx, bobAddr, do)
	if err != nil {
		t.Fatal(err)
	}
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
	err := k.iterDomainToOwner(ctx, bobAddr, do)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(domainKeys); l != 1 {
		t.Fatalf("expected %d keys got: %d", 1, l)
	}
	domainKeys = nil
	// transfer domain
	domain, _ := k.GetDomain(ctx, "1")
	k.TransferDomainAll(ctx, aliceAddr, domain)
	// check if addr b has 0 keys
	err = k.iterDomainToOwner(ctx, bobAddr, do)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(domainKeys); l != 0 {
		t.Fatalf("expected %d keys got: %d", 0, l)
	}
	domainKeys = nil
	// check if addr a has 2 keys
	err = k.iterDomainToOwner(ctx, aliceAddr, do)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(domainKeys); l != 2 {
		t.Fatalf("expected %d keys got: %d", 2, l)
	}
	domainKeys = nil
	// delete domain
	_ = k.DeleteDomain(ctx, "2")
	// check if addr b has 1 key
	err = k.iterDomainToOwner(ctx, aliceAddr, do)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(domainKeys); l != 1 {
		t.Fatalf("expected %d keys got: %d", 1, l)
	}
}

func Test_resourceIndexing(t *testing.T) {
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
	// create one resource
	resourceA := types.Resource{
		URI:      "t1",
		Resource: "1",
	}
	resourceB := types.Resource{
		URI:      "t2",
		Resource: "2",
	}
	// create one account
	accountA := types.Account{
		Domain: "test",
		Name:   "1",
		Resources: []types.Resource{
			resourceA,
			resourceB,
		},
		Owner: aliceAddr,
	}
	// insert account
	k.CreateAccount(ctx, accountA)
	// iterate resources
	err := k.iterateResourceAccounts(ctx, resourceA, do)
	if err != nil {
		t.Fatal(err)
	}
	if len(accountKeys) != 1 {
		t.Fatalf("expected 1 keys, got: %d", len(accountKeys))
	}
	// generate test account
	acc := &types.Account{}
	err = acc.Unpack(accountKeys[0])
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
	err = k.iterateResourceAccounts(ctx, resourceA, do)
	if err != nil {
		t.Fatal(err)
	}
	if len(accountKeys) != 0 {
		t.Fatalf("no key expected, got: %d", len(accountKeys))
	}
	// ReplaceAccountResources
	accountKeys = nil
	k.CreateAccount(ctx, accountA)
	k.ReplaceAccountResources(ctx, accountA, []types.Resource{resourceB})
	err = k.iterateResourceAccounts(ctx, resourceA, do)
	if err != nil {
		t.Fatal(err)
	}
	if len(accountKeys) != 0 {
		t.Fatalf("no key expected, got: %d", len(accountKeys))
	}
	accountKeys = nil
	err = k.iterateResourceAccounts(ctx, resourceB, do)
	if err != nil {
		t.Fatal(err)
	}
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
	// check if resourceA is associated with any account
	err = k.iterateResourceAccounts(ctx, resourceA, do)
	if err != nil {
		t.Fatal(err)
	}
	if len(accountKeys) != 0 {
		t.Fatalf("expected 0 keys, got: %d", len(accountKeys))
	}
}

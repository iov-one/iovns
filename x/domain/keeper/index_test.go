package keeper

import (
	"fmt"
	"github.com/iov-one/iovns/pkg/index"
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
	k.TransferDomain(ctx, aliceAddr, domain)
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
	err := k.iterateBlockchainTargetsAccounts(ctx, targetA, do)
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
	err = k.iterateBlockchainTargetsAccounts(ctx, targetA, do)
	if err != nil {
		t.Fatal(err)
	}
	if len(accountKeys) != 0 {
		t.Fatalf("no key expected, got: %d", len(accountKeys))
	}
	// ReplaceAccountTargets
	accountKeys = nil
	k.CreateAccount(ctx, accountA)
	k.ReplaceAccountTargets(ctx, accountA, []types.BlockchainAddress{targetB})
	err = k.iterateBlockchainTargetsAccounts(ctx, targetA, do)
	if err != nil {
		t.Fatal(err)
	}
	if len(accountKeys) != 0 {
		t.Fatalf("no key expected, got: %d", len(accountKeys))
	}
	accountKeys = nil
	err = k.iterateBlockchainTargetsAccounts(ctx, targetB, do)
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
	// check if targetA is associated with any account
	err = k.iterateBlockchainTargetsAccounts(ctx, targetA, do)
	if err != nil {
		t.Fatal(err)
	}
	if len(accountKeys) != 0 {
		t.Fatalf("expected 0 keys, got: %d", len(accountKeys))
	}
}

func Test_certificatesIndexing(t *testing.T) {
	k, ctx, _ := NewTestKeeper(t, true)
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
	certA := types.Certificate("test")
	certB := types.Certificate("test2")
	accountA := types.Account{
		Domain:       "test",
		Name:         "1",
		Owner:        bobAddr,
		ValidUntil:   0,
		Certificates: []types.Certificate{certA},
		MetadataURI:  "",
	}
	accountB := types.Account{
		Domain:       "test",
		Name:         "2",
		Owner:        aliceAddr,
		ValidUntil:   0,
		Certificates: []types.Certificate{certA, certB},
		MetadataURI:  "",
	}
	// create accounts
	k.CreateAccount(ctx, accountA)
	k.CreateAccount(ctx, accountB)
	// get certs
	err := k.iterateCertificateAccounts(ctx, certA, do)
	if err != nil {
		t.Fatal(err)
	}
	if len(accountKeys) != 2 {
		for _, k := range accountKeys {
			t.Logf("%s", k)
		}
		t.Fatalf("expected 2 keys, got: %d", len(accountKeys))
	}
	acc := types.Account{}
	index.MustUnpack(accountKeys[0], &acc)
	// check if accounts match
	if err := accMatch(acc, accountA); err != nil {
		t.Fatal(err)
	}
	index.MustUnpack(accountKeys[1], &acc)
	if err := accMatch(acc, accountB); err != nil {
		t.Fatal(err)
	}
	// delete account
	accountKeys = nil
	k.DeleteAccount(ctx, accountB.Domain, accountB.Name)
	err = k.iterateCertificateAccounts(ctx, certA, do)
	if err != nil {
		t.Fatal(err)
	}
	if len(accountKeys) != 1 {
		t.Fatalf("expected 1 key, got: %d", len(accountKeys))
	}
	// check if A is the only account with the key
	index.MustUnpack(accountKeys[0], &acc)
	if err := accMatch(acc, accountA); err != nil {
		t.Fatal(err)
	}
	// transfer account
	accountKeys = nil
	k.TransferAccount(ctx, accountA, aliceAddr)
	// check if certs has no matches
	err = k.iterateCertificateAccounts(ctx, certA, do)
	if err != nil {
		t.Fatal(err)
	}
	if len(accountKeys) != 0 {
		t.Fatalf("expected 0 keys, got: %d", len(accountKeys))
	}
	// now add certificates
	accountKeys = nil
	accountA, _ = k.GetAccount(ctx, accountA.Domain, accountA.Name) // get updated account
	k.AddAccountCertificate(ctx, accountA, certB)                   // add cert
	// check if accountA has B cert
	err = k.iterateCertificateAccounts(ctx, certB, do)
	if err != nil {
		t.Fatal(err)
	}
	if len(accountKeys) != 1 {
		t.Fatalf("expected 0 keys, got: %d", len(accountKeys))
	}
	// check that accountA is correctly matched to certB
	index.MustUnpack(accountKeys[0], &acc)
	if err := accMatch(acc, accountA); err != nil {
		t.Fatal(err)
	}
	// delete cert
	accountKeys = nil
	accountA, _ = k.GetAccount(ctx, accountA.Domain, accountA.Name)
	k.DeleteAccountCertificate(ctx, accountA, 0)
	err = k.iterateCertificateAccounts(ctx, certB, do)
	if err != nil {
		t.Fatal(err)
	}
	if len(accountKeys) != 0 {
		t.Fatalf("expected 0 keys, got: %d", len(accountKeys))
	}
}

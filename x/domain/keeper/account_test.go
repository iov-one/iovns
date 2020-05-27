package keeper

import (
	"fmt"
	"github.com/iov-one/iovns/x/domain/types"
	"testing"
)

func TestKeeper_IterateAllAccounts(t *testing.T) {
	k, ctx, _ := NewTestKeeper(t, true)
	n := 100
	accountSet := make(map[string]string, 100)
	for i := 0; i < n; i++ {
		acc := types.Account{
			Name:       fmt.Sprintf("%d", i),
			Domain:     fmt.Sprintf("%d", i%3),
			ValidUntil: 1000,
			Broker:     nil,
		}
		k.SetAccount(ctx, acc)
		accountSet[acc.Name] = acc.Domain
	}
	accounts := k.IterateAllAccounts(ctx)
	if len(accounts) != n {
		t.Fatalf("IterateAllAccounts() expected: %d accounts, got %d", n, len(accounts))
	}
	// check if all account names are there
	for _, account := range accounts {
		domain, ok := accountSet[account.Name]
		if !ok {
			t.Fatalf("IterateAllAccounts() unwanted account: %s", account.Name)
		}
		if account.Domain != domain {
			t.Fatalf("IterateAllAccounts() expected domain name: %s, got: %s", domain, account.Domain)
		}
	}
}

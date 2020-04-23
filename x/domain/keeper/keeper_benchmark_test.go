package keeper

import (
	"fmt"
	"github.com/iov-one/iovns/x/domain/types"
	"testing"
)

func Benchmark_accountToOwnerIndexing(b *testing.B) {
	// get keeper
	k, ctx := NewTestKeeper(b, true)
	// generate mock accounts
	number := 100000
	for i := 0; i < number; i++ {
		k.CreateAccount(ctx, types.Account{
			Domain: fmt.Sprintf("test"),
			Name:   fmt.Sprintf("%d", i),
			Owner:  aliceAddr,
		})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k.iterAccountToOwner(ctx, aliceAddr, func(key []byte) bool {
			return true
		})
	}
}

func Benchmark_domainToOwnerIndexing(b *testing.B) {
	// get keeper
	k, ctx := NewTestKeeper(b, true)
	// generate mock domains
	number := 100000
	for i := 0; i < number; i++ {
		k.CreateDomain(ctx, types.Domain{
			Name:  fmt.Sprintf("%d", i),
			Admin: aliceAddr,
		})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k.iterDomainToOwner(ctx, aliceAddr, func(key []byte) bool {
			return true
		})
	}
}

func Benchmark_domainToAccountIndexing(b *testing.B) {
	// get keeper
	k, ctx := NewTestKeeper(b, true)
	// generate mock accounts
	number := 100000
	for i := 0; i < number; i++ {
		k.CreateAccount(ctx, types.Account{
			Domain: fmt.Sprintf("test"),
			Name:   fmt.Sprintf("%d", i),
			Owner:  aliceAddr,
		})
	}
	b.ResetTimer()
	// iterate domain
	for i := 0; i < b.N; i++ {
		k.GetAccountsInDomain(ctx, "test", func(key []byte) bool {
			return true
		})
	}
}

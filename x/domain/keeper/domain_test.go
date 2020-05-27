package keeper

import (
	"fmt"
	"testing"

	"github.com/iov-one/iovns/x/domain/types"
)

func TestKeeper_IterateAllDomains(t *testing.T) {
	k, ctx, _ := NewTestKeeper(t, true)
	n := 100
	domainSet := make(map[string]struct{}, 100)
	for i := 0; i < n; i++ {
		k.SetDomain(ctx, types.Domain{
			Name:         fmt.Sprintf("%d", i),
			Admin:        nil,
			ValidUntil:   1000,
			Type:         types.ClosedDomain,
			AccountRenew: 1000000,
			Broker:       nil,
		})
		domainSet[fmt.Sprintf("%d", i)] = struct{}{}
	}
	domains := k.IterateAllDomains(ctx)
	if len(domains) != n {
		t.Fatalf("IterateAllDomains() expected: %d domains, got %d", n, len(domains))
	}
	// check if all domain names are there
	for _, domain := range domains {
		if _, ok := domainSet[domain.Name]; !ok {
			t.Fatalf("IterateAllDomains() unwanted domain: %s", domain.Name)
		}
	}
}

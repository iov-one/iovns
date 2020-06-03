package keeper

import (
	"fmt"
	"reflect"
	"testing"
	"time"

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

func TestKeeper_CreateDomain(t *testing.T) {
	k, ctx, _ := NewTestKeeper(t, true)
	ctx.WithBlockTime(time.Unix(0, 0))
	// create mock domains
	closedDomain := types.Domain{
		Name:         "closed",
		Admin:        AliceKey,
		ValidUntil:   1,
		Type:         types.ClosedDomain,
		AccountRenew: 2,
		Broker:       nil,
	}
	openDomain := types.Domain{
		Name:         "open",
		Admin:        AliceKey,
		ValidUntil:   1,
		Type:         types.OpenDomain,
		AccountRenew: 2,
		Broker:       nil,
	}
	k.CreateDomain(ctx, closedDomain)
	k.CreateDomain(ctx, openDomain)
	t.Run("closed domain", func(t *testing.T) {
		dom, ok := k.GetDomain(ctx, "closed")
		if !ok {
			t.Fatalf("domain not found")
		}
		if !reflect.DeepEqual(dom, closedDomain) {
			t.Fatalf("expected: %+v, got: %+v", closedDomain, dom)
		}
		// check empty account
		acc, ok := k.GetAccount(ctx, "closed", "")
		if !ok {
			t.Fatalf("empty account not found")
		}
		if acc.ValidUntil != types.MaxValidUntil {
			t.Fatalf("unexpected valid until: %d", acc.ValidUntil)
		}
	})
	t.Run("closed domain", func(t *testing.T) {
		dom, _ := k.GetDomain(ctx, "open")
		// check empty account
		acc, ok := k.GetAccount(ctx, "open", "")
		if !ok {
			t.Fatalf("empty account not found")
		}
		if acc.ValidUntil != dom.ValidUntil {
			t.Fatalf("unexpected valid until: %d", acc.ValidUntil)
		}
	})
}

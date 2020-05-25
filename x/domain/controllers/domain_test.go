package controllers

import (
	"errors"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration"
	dt "github.com/iov-one/iovns/x/domain/testing"

	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
	"github.com/stretchr/testify/assert"
)

func TestDomain_requireDomain(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		k, ctx, _ := keeper.NewTestKeeper(t, true)
		k.CreateDomain(ctx, types.Domain{
			Name:         "test",
			Admin:        dt.AliceKey,
			HasSuperuser: false,
		})
		ctrl := NewDomainController(ctx, k, "test")
		err := ctrl.requireDomain()
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
	t.Run("does not exist", func(t *testing.T) {
		k, ctx, _ := keeper.NewTestKeeper(t, true)
		ctrl := NewDomainController(ctx, k, "test")
		err := ctrl.requireDomain()
		if !errors.Is(err, types.ErrDomainDoesNotExist) {
			t.Fatalf("want: %s, got: %s", types.ErrAccountDoesNotExist, err)
		}
	})
}

func TestDomain_domainExpired(t *testing.T) {
	t.Run("domain expired", func(t *testing.T) {
		k, ctx, _ := keeper.NewTestKeeper(t, true)
		k.CreateDomain(ctx, types.Domain{
			Name:         "test",
			Admin:        dt.AliceKey,
			HasSuperuser: false,
			ValidUntil:   0,
		})
		ctrl := NewDomainController(ctx, k, "test")
		expired := ctrl.domainExpired()
		if !expired {
			t.Fatal("validation failed: domain has not expired")
		}
	})
	t.Run("domain not expired", func(t *testing.T) {
		k, ctx, _ := keeper.NewTestKeeper(t, true)
		now := time.Now()
		k.CreateDomain(ctx, types.Domain{
			Name:       "test",
			Admin:      dt.AliceKey,
			ValidUntil: now.Unix() + 10000,
		})
		ctrl := NewDomainController(ctx, k, "test")
		expired := ctrl.domainExpired()
		if expired {
			t.Fatal("validation failed domain has expired")
		}
	})
	t.Run("domain does not exist", func(t *testing.T) {
		k, ctx, _ := keeper.NewTestKeeper(t, true)
		ctrl := NewDomainController(ctx, k, "test")
		assert.Panics(t, func() { _ = ctrl.domainExpired() }, "domain does not exists")
	})
}

func TestDomain_gracePeriodFinished(t *testing.T) {
	cases := map[string]dt.SubTest{
		"grace period finished": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 1,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      dt.AliceKey,
					ValidUntil: 0,
				})
			},
			TestBlockTime: 2,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ctrl := NewDomainController(ctx, k, "test")
				gpf := ctrl.gracePeriodFinished()
				if gpf != true {
					t.Fatal("validation failed: grace period has not expired")
				}
			},
		},
	}
	dt.RunTests(t, cases)
}

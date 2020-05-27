package domain

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
			Name:  "test",
			Admin: dt.AliceKey,
			Type:  types.OpenDomain,
		})
		ctrl := NewController(ctx, k, "test")
		err := ctrl.requireDomain()
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
	t.Run("does not exist", func(t *testing.T) {
		k, ctx, _ := keeper.NewTestKeeper(t, true)
		ctrl := NewController(ctx, k, "test")
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
			Name:       "test",
			Admin:      dt.AliceKey,
			Type:       types.OpenDomain,
			ValidUntil: 0,
		})
		ctrl := NewController(ctx, k, "test")
		err := ctrl.expired()
		if err != nil {
			t.Fatalf("unexpected err: %s", err)
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
		ctrl := NewController(ctx, k, "test")
		err := ctrl.expired()
		if !errors.Is(err, types.ErrDomainNotExpired) {
			t.Fatalf("expected error: %s, got: %s", types.ErrDomainNotExpired, err)
		}
	})
	t.Run("domain does not exist", func(t *testing.T) {
		k, ctx, _ := keeper.NewTestKeeper(t, true)
		ctrl := NewController(ctx, k, "test")
		assert.Panics(t, func() { _ = ctrl.expired() }, "domain does not exists")
	})
}

func TestDomain_gracePeriodFinished(t *testing.T) {
	cases := map[string]dt.SubTest{
		"grace period finished": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 1 * time.Second,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      dt.AliceKey,
					ValidUntil: 0,
				})
			},
			TestBlockTime: 10,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ctrl := NewController(ctx, k, "test")
				err := ctrl.gracePeriodFinished()
				if err != nil {
					t.Fatal("validation failed: grace period has not expired")
				}
			},
		},
		"grace period not finished": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 15 * time.Second,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      dt.AliceKey,
					ValidUntil: 1,
				})
			},
			TestBlockTime: 3,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ctrl := NewController(ctx, k, "test")
				err := ctrl.gracePeriodFinished()
				if !errors.Is(err, types.ErrGracePeriodNotFinished) {
					t.Fatalf("expected error: %s, got: %s", types.ErrGracePeriodNotFinished, err)
				}
			},
		},
	}
	dt.RunTests(t, cases)
}

func TestDomain_ownedBy(t *testing.T) {
	cases := map[string]dt.SubTest{
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      dt.AliceKey,
					ValidUntil: 0,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ctrl := NewController(ctx, k, "test")
				err := ctrl.ownedBy(dt.AliceKey)
				if err != nil {
					t.Fatalf("got error: %s", err)
				}
			},
		},
		"unauthorized": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      dt.AliceKey,
					ValidUntil: 0,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ctrl := NewController(ctx, k, "test")
				err := ctrl.ownedBy(dt.BobKey)
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("want err: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
		},
	}
	dt.RunTests(t, cases)
}

func TestDomain_notExpired(t *testing.T) {
	cases := map[string]dt.SubTest{
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      dt.AliceKey,
					ValidUntil: 2,
				})
			},
			TestBlockTime: 1,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ctrl := NewController(ctx, k, "test")
				err := ctrl.notExpired()
				if err != nil {
					t.Fatalf("got error: %s", err)
				}
			},
		},
		"expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      dt.AliceKey,
					ValidUntil: 1,
				})
			},
			TestBlockTime: 2,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ctrl := NewController(ctx, k, "test")
				err := ctrl.notExpired()
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("want err: %s, got: %s", types.ErrDomainExpired, err)
				}
			},
		},
	}
	dt.RunTests(t, cases)
}

func TestDomain_type(t *testing.T) {
	cases := map[string]dt.SubTest{
		"saved": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:  "test",
					Admin: dt.AliceKey,
					Type:  types.ClosedDomain,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ctrl := NewController(ctx, k, "test")
				err := ctrl.dType(types.ClosedDomain)
				if err != nil {
					t.Fatalf("got error: %s", err)
				}
			},
		},
		"fail want type close domain": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:  "test",
					Admin: dt.AliceKey,
					Type:  types.ClosedDomain,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ctrl := NewController(ctx, k, "test")
				err := ctrl.dType(types.OpenDomain)
				if !errors.Is(err, types.ErrInvalidDomainType) {
					t.Fatalf("want err: %s, got: %s", types.ErrInvalidDomainType, err)
				}
			},
		},
		"fail want open domain": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:  "test",
					Admin: dt.AliceKey,
					Type:  types.OpenDomain,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ctrl := NewController(ctx, k, "test")
				err := ctrl.dType(types.ClosedDomain)
				if !errors.Is(err, types.ErrInvalidDomainType) {
					t.Fatalf("want err: %s, got: %s", types.ErrInvalidDomainType, err)
				}
			},
		},
	}
	dt.RunTests(t, cases)
}

func TestDomain_validName(t *testing.T) {
	cases := map[string]dt.SubTest{
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidDomain: dt.RegexMatchAll,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      dt.AliceKey,
					ValidUntil: 0,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ctrl := NewController(ctx, k, "test")
				err := ctrl.validName()
				if err != nil {
					t.Fatalf("got error: %s", err)
				}
			},
		},
		"invalid name": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidDomain: dt.RegexMatchNothing,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      dt.AliceKey,
					ValidUntil: 0,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ctrl := NewController(ctx, k, "test")
				err := ctrl.validName()
				if !errors.Is(err, types.ErrInvalidDomainName) {
					t.Fatalf("want err: %s, got: %s", types.ErrInvalidDomainName, err)
				}
			},
		},
	}
	dt.RunTests(t, cases)
}

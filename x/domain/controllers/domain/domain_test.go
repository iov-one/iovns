package domain

import (
	"errors"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
	"github.com/stretchr/testify/assert"
)

func TestDomain_requireDomain(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		k, ctx, _ := keeper.NewTestKeeper(t, true)
		ds := k.DomainStore(ctx)
		ds.Create(&types.Domain{
			Name:  "test",
			Admin: keeper.AliceKey,
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
		ds := k.DomainStore(ctx)
		ds.Create(&types.Domain{
			Name:       "test",
			Admin:      keeper.AliceKey,
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
		ds := k.DomainStore(ctx)
		now := time.Now()
		ds.Create(&types.Domain{
			Name:       "test",
			Admin:      keeper.AliceKey,
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
	cases := map[string]keeper.SubTest{
		"grace period finished": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 1 * time.Second,
				})
				ds := k.DomainStore(ctx)
				ds.Create(&types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
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
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 15 * time.Second,
				})
				ds := k.DomainStore(ctx)
				ds.Create(&types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					ValidUntil: 1,
				})
			},
			TestBlockTime: 3,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ctrl := NewController(ctx, k, "test")
				err := ctrl.gracePeriodFinished()
				if !errors.Is(err, types.ErrDomainGracePeriodNotFinished) {
					t.Fatalf("expected error: %s, got: %s", types.ErrDomainGracePeriodNotFinished, err)
				}
			},
		},
	}
	keeper.RunTests(t, cases)
}

func TestDomain_ownedBy(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ds := k.DomainStore(ctx)
				ds.Create(&types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					ValidUntil: 0,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ctrl := NewController(ctx, k, "test")
				err := ctrl.isAdmin(keeper.AliceKey)
				if err != nil {
					t.Fatalf("got error: %s", err)
				}
			},
		},
		"unauthorized": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ds := k.DomainStore(ctx)
				ds.Create(&types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					ValidUntil: 0,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ctrl := NewController(ctx, k, "test")
				err := ctrl.isAdmin(keeper.BobKey)
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("want err: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
		},
	}
	keeper.RunTests(t, cases)
}

func TestDomain_notExpired(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ds := k.DomainStore(ctx)
				ds.Create(&types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
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
				ds := k.DomainStore(ctx)
				ds.Create(&types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
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
	keeper.RunTests(t, cases)
}

func TestDomain_type(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"saved": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ds := k.DomainStore(ctx)
				ds.Create(&types.Domain{
					Name:  "test",
					Admin: keeper.AliceKey,
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
				ds := k.DomainStore(ctx)
				ds.Create(&types.Domain{
					Name:  "test",
					Admin: keeper.AliceKey,
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
				ds := k.DomainStore(ctx)
				ds.Create(&types.Domain{
					Name:  "test",
					Admin: keeper.AliceKey,
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
	keeper.RunTests(t, cases)
}

func TestDomain_validName(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidDomainName: keeper.RegexMatchAll,
				})
				ds := k.DomainStore(ctx)
				ds.Create(&types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
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
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidDomainName: keeper.RegexMatchNothing,
				})
				ds := k.DomainStore(ctx)
				ds.Create(&types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
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
	keeper.RunTests(t, cases)
}

func TestDomain_Renewable(t *testing.T) {
	k, ctx, _ := keeper.NewTestKeeper(t, true)
	ctx = ctx.WithBlockTime(time.Unix(1, 0))
	setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
	setConfig(ctx, configuration.Config{
		DomainRenewalCountMax: 1, // increased by one inside controller
		DomainRenewalPeriod:   10 * time.Second,
	})
	ds := k.DomainStore(ctx)
	ds.Create(&types.Domain{
		Name:       "open",
		Admin:      keeper.AliceKey,
		ValidUntil: time.Unix(18, 0).Unix(),
		Type:       types.OpenDomain,
	})
	ds.Create(&types.Domain{
		Name:       "closed",
		Admin:      keeper.AliceKey,
		ValidUntil: time.Unix(18, 0).Unix(),
		Type:       types.ClosedDomain,
	})

	// 18(DomainValidUntil) + 10 (DomainRP) = 28 newValidUntil
	t.Run("open domain", func(t *testing.T) {
		// 7(time) + 2(DomainRCM) * 10(DomainRP) = 27 maxValidUntil
		d := NewController(ctx.WithBlockTime(time.Unix(7, 0)), k, "open")
		err := d.Validate(Renewable)
		if !errors.Is(err, types.ErrUnauthorized) {
			t.Fatalf("want: %s, got: %s", types.ErrUnauthorized, err)
		}
		// 100(time) + 2(DomainRCM) * 10(DomainRP) = 120 maxValidUntil
		d = NewController(ctx.WithBlockTime(time.Unix(100, 0)), k, "open")
		if err := d.Validate(Renewable); err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
	// 18(DomainValidUntil) + 10 (DomainRP) = 28 newValidUntil
	t.Run("closed domain", func(t *testing.T) {
		// 7(time) + 2(DomainRCM) * 10(DomainRP) = 27 maxValidUntil
		d := NewController(ctx.WithBlockTime(time.Unix(7, 0)), k, "closed")
		err := d.Validate(Renewable)
		if !errors.Is(err, types.ErrUnauthorized) {
			t.Fatalf("want: %s, got: %s", types.ErrUnauthorized, err)
		}
		// 100(time) + 2(DomainRCM) * 10(DomainRP) = 120 maxValidUntil
		d = NewController(ctx.WithBlockTime(time.Unix(100, 0)), k, "closed")
		if err := d.Validate(Renewable); err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
}

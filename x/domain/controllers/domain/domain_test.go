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
			Name:         "test",
			Admin:        dt.AliceKey,
			HasSuperuser: false,
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
			Name:         "test",
			Admin:        dt.AliceKey,
			HasSuperuser: false,
			ValidUntil:   0,
		})
		ctrl := NewController(ctx, k, "test")
		expired := ctrl.expired()
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
		ctrl := NewController(ctx, k, "test")
		expired := ctrl.expired()
		if expired {
			t.Fatal("validation failed domain has expired")
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
				ctrl := NewController(ctx, k, "test")
				gpf := ctrl.gracePeriodFinished()
				if gpf != true {
					t.Fatal("validation failed: grace period has not expired")
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

func TestDomain_superuser(t *testing.T) {
	cases := map[string]dt.SubTest{
		"success true": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        dt.AliceKey,
					HasSuperuser: true,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ctrl := NewController(ctx, k, "test")
				err := ctrl.superuser(true)
				if err != nil {
					t.Fatalf("got error: %s", err)
				}
			},
		},
		"success false": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        dt.AliceKey,
					HasSuperuser: false,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ctrl := NewController(ctx, k, "test")
				err := ctrl.superuser(false)
				if err != nil {
					t.Fatalf("got error: %s", err)
				}
			},
		},
		"fail superuser want true": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        dt.AliceKey,
					HasSuperuser: true,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ctrl := NewController(ctx, k, "test")
				err := ctrl.superuser(false)
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("want err: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
		},
		"fail superuser want false": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        dt.AliceKey,
					HasSuperuser: false,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				ctrl := NewController(ctx, k, "test")
				err := ctrl.superuser(true)
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("want err: %s, got: %s", types.ErrUnauthorized, err)
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

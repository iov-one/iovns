package account

import (
	"errors"
	"github.com/iov-one/iovns/x/domain/controllers/domain"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/mock"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

func TestAccount_transferable(t *testing.T) {
	k, ctx, _ := keeper.NewTestKeeper(t, true)
	// create mock domains and accounts
	// create open domain
	ds := k.DomainStore(ctx)
	as := k.AccountStore(ctx)
	ds.Create(&types.Domain{
		Name:       "open",
		Admin:      keeper.AliceKey,
		ValidUntil: time.Now().Add(100 * time.Hour).Unix(),
		Type:       types.OpenDomain,
	})
	// creat open domain account
	as.Create(&types.Account{
		Domain: "open",
		Name:   "test",
		Owner:  keeper.BobKey,
	})
	// create closed domain
	ds.Create(&types.Domain{
		Name:       "closed",
		Admin:      keeper.AliceKey,
		ValidUntil: time.Now().Add(100 * time.Hour).Unix(),
		Type:       types.ClosedDomain,
	})
	// create closed domain account
	as.Create(&types.Account{
		Domain: "closed",
		Name:   "test",
		Owner:  keeper.BobKey,
	})
	// run tests
	t.Run("closed domain", func(t *testing.T) {
		acc := NewController(ctx, k, "closed", "test")
		// test success
		err := acc.Validate(TransferableBy(keeper.AliceKey))
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
		// test failure
		err = acc.Validate(TransferableBy(keeper.BobKey))
		if !errors.Is(err, types.ErrUnauthorized) {
			t.Fatalf("want: %s, got: %s", types.ErrUnauthorized, err)
		}
	})
	t.Run("open domain", func(t *testing.T) {
		acc := NewController(ctx, k, "open", "test")
		err := acc.Validate(TransferableBy(keeper.BobKey))
		// test success
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
		// test failure
		err = acc.Validate(TransferableBy(keeper.AliceKey))
		if !errors.Is(err, types.ErrUnauthorized) {
			t.Fatalf("want: %s, got: %s", types.ErrUnauthorized, err)
		}
	})
}

func TestAccount_existence(t *testing.T) {
	k, ctx, _ := keeper.NewTestKeeper(t, true)
	as := k.AccountStore(ctx)
	// insert mock account
	as.Create(&types.Account{
		Domain:     "test",
		Name:       "test",
		Owner:      keeper.AliceKey,
		ValidUntil: time.Now().Add(100 * time.Hour).Unix(),
	})
	// run MustExist test
	t.Run("must exist success", func(t *testing.T) {
		acc := NewController(ctx, k, "test", "test")
		err := acc.Validate(MustExist)
		if err != nil {
			t.Errorf("got error: %s", err)
		}
	})
	t.Run("must exist fail", func(t *testing.T) {
		acc := NewController(ctx, k, "test", "does not exist")
		err := acc.Validate(MustExist)
		if !errors.Is(err, types.ErrAccountDoesNotExist) {
			t.Fatalf("want: %s, got: %s", types.ErrAccountDoesNotExist, err)
		}
	})
	// run MustNotExist test
	t.Run("must not exist success", func(t *testing.T) {
		acc := NewController(ctx, k, "test", "does not exist")
		err := acc.Validate(MustNotExist)
		if err != nil {
			t.Errorf("got error: %s", err)
		}
	})
	t.Run("must not exist fail", func(t *testing.T) {
		acc := NewController(ctx, k, "test", "test")
		err := acc.Validate(MustNotExist)
		if !errors.Is(err, types.ErrAccountExists) {
			t.Fatalf("want: %s, got: %s", types.ErrAccountExists, err)
		}
	})
}

func TestAccount_requireAccount(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		k, ctx, _ := keeper.NewTestKeeper(t, true)
		as := k.AccountStore(ctx)
		alice, _ := mock.Addresses()
		as.Create(&types.Account{
			Domain: "test",
			Name:   "test",
			Owner:  alice,
		})
		ctrl := NewController(ctx, k, "test", "test")
		err := ctrl.requireAccount()
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
	t.Run("does not exist", func(t *testing.T) {
		k, ctx, _ := keeper.NewTestKeeper(t, true)
		ctrl := NewController(ctx, k, "test", "test")
		err := ctrl.requireAccount()
		if !errors.Is(err, types.ErrAccountDoesNotExist) {
			t.Fatalf("want: %s, got: %s", types.ErrAccountDoesNotExist, err)
		}
	})
}

func TestAccount_certNotExist(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		acc := &Account{
			account: &types.Account{
				Certificates: []types.Certificate{[]byte("test-cert")},
			},
		}
		err := acc.Validate(CertificateNotExist([]byte("does not exist")))
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
	t.Run("cert exists", func(t *testing.T) {
		acc := &Account{
			account: &types.Account{
				Certificates: []types.Certificate{[]byte("test-cert"), []byte("exists")},
			},
		}
		i := new(int)
		err := acc.Validate(CertificateExists([]byte("exists"), i))
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
		if *i != 1 {
			t.Fatalf("unexpected index pointer: %d", *i)
		}
	})
}

func TestAccount_notExpired(t *testing.T) {
	closedDomain := (&domain.Domain{}).WithDomain(types.Domain{
		Type: types.ClosedDomain,
	})
	openDomain := (&domain.Domain{}).WithDomain(types.Domain{
		Type: types.OpenDomain,
	})
	t.Run("success", func(t *testing.T) {
		acc := (&Account{
			account: &types.Account{
				ValidUntil: 10,
			},
			ctx: sdk.Context{}.WithBlockTime(time.Unix(0, 0)),
		}).WithDomainController(openDomain)
		err := acc.Validate(NotExpired)
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
	t.Run("expired", func(t *testing.T) {
		acc := (&Account{
			account: &types.Account{
				ValidUntil: 10,
			},
			ctx: sdk.Context{}.WithBlockTime(time.Unix(11, 0)),
		}).WithDomainController(openDomain)
		err := acc.Validate(NotExpired)
		if !errors.Is(err, types.ErrAccountExpired) {
			t.Fatalf("want error: %s, got: %s", types.ErrAccountExpired, err)
		}
	})
	t.Run("success account expired but in closed domain", func(t *testing.T) {
		acc := (&Account{
			account: &types.Account{
				ValidUntil: 1,
			},
			ctx: sdk.Context{}.WithBlockTime(time.Unix(20, 0)),
		}).WithDomainController(closedDomain)
		err := acc.Validate(NotExpired)
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
}

func TestAccount_ownedBy(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		alice, _ := mock.Addresses()
		acc := &Account{
			account: &types.Account{Owner: alice},
		}
		err := acc.Validate(Owner(alice))
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
	t.Run("bad owner", func(t *testing.T) {
		alice, bob := mock.Addresses()
		acc := &Account{
			account: &types.Account{Owner: alice},
		}
		err := acc.Validate(Owner(bob))
		if !errors.Is(err, types.ErrUnauthorized) {
			t.Fatalf("unexpected error: %s, wanted: %s", err, types.ErrUnauthorized)
		}
	})
}

func TestAccount_validName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		acc := &Account{
			account: &types.Account{Name: "valid"},
			conf:    &configuration.Config{ValidAccountName: "^(.*?)?"},
		}
		err := acc.Validate(ValidName)
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
	t.Run("success", func(t *testing.T) {
		acc := &Account{
			name: "not valid",
			conf: &configuration.Config{ValidAccountName: "$^"},
		}
		err := acc.Validate(ValidName)
		if !errors.Is(err, types.ErrInvalidAccountName) {
			t.Fatalf("unexpected error: %s, wanted: %s", err, types.ErrInvalidAccountName)
		}
	})
}

func TestAccountRegistrableBy(t *testing.T) {
	closedDomain := (&domain.Domain{}).WithDomain(types.Domain{
		Type:  types.ClosedDomain,
		Admin: keeper.AliceKey,
	})
	openDomain := (&domain.Domain{}).WithDomain(types.Domain{
		Type: types.OpenDomain,
	})
	t.Run("success in closed domain", func(t *testing.T) {
		acc := (&Account{}).WithDomainController(closedDomain)
		err := acc.Validate(RegistrableBy(keeper.AliceKey))
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
	t.Run("fail in closed domain", func(t *testing.T) {
		acc := (&Account{}).WithDomainController(closedDomain)
		err := acc.Validate(RegistrableBy(keeper.BobKey))
		if !errors.Is(err, types.ErrUnauthorized) {
			t.Fatalf("want: %s, got: %s", types.ErrUnauthorized, err)
		}
	})
	t.Run("success other domain type", func(t *testing.T) {
		acc := (&Account{}).WithDomainController(openDomain)
		err := acc.Validate(RegistrableBy(keeper.AliceKey))
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
}

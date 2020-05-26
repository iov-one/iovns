package account

import (
	"errors"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/mock"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

func TestAccount_requireAccount(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		k, ctx, _ := keeper.NewTestKeeper(t, true)
		alice, _ := mock.Addresses()
		k.SetAccount(ctx, types.Account{
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
		err := acc.certNotExist([]byte("does not exist"), nil)
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
		err := acc.certNotExist([]byte("exists"), i)
		if !errors.Is(err, types.ErrCertificateExists) {
			t.Fatalf("unexpected error: %s, wanted: %s", err, types.ErrCertificateExists)
		}
		if *i != 1 {
			t.Fatalf("unexpected index pointer: %d", *i)
		}
	})
}

func TestAccount_notExpired(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		acc := &Account{
			account: &types.Account{
				ValidUntil: 10,
			},
			ctx: sdk.Context{}.WithBlockTime(time.Unix(0, 0)),
		}
		err := acc.notExpired()
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
	t.Run("expired", func(t *testing.T) {
		acc := &Account{
			account: &types.Account{
				ValidUntil: 10,
			},
			ctx: sdk.Context{}.WithBlockTime(time.Unix(11, 0)),
		}
		err := acc.notExpired()
		if !errors.Is(err, types.ErrAccountExpired) {
			t.Fatalf("want error: %s, got: %s", types.ErrAccountExpired, err)
		}
	})
}

func TestAccount_ownedBy(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		alice, _ := mock.Addresses()
		acc := &Account{
			account: &types.Account{Owner: alice},
		}
		err := acc.ownedBy(alice)
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
	t.Run("bad owner", func(t *testing.T) {
		alice, bob := mock.Addresses()
		acc := &Account{
			account: &types.Account{Owner: alice},
		}
		err := acc.ownedBy(bob)
		if !errors.Is(err, types.ErrUnauthorized) {
			t.Fatalf("unexpected error: %s, wanted: %s", err, types.ErrUnauthorized)
		}
	})
}

func TestAccount_validName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		acc := &Account{
			account: &types.Account{Name: "valid"},
			conf:    &configuration.Config{ValidName: "^(.*?)?"},
		}
		err := acc.validName()
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
	t.Run("success", func(t *testing.T) {
		acc := &Account{
			name: "not valid",
			conf: &configuration.Config{ValidName: "$^"},
		}
		err := acc.validName()
		if !errors.Is(err, types.ErrInvalidAccountName) {
			t.Fatalf("unexpected error: %s, wanted: %s", err, types.ErrInvalidAccountName)
		}
	})
}

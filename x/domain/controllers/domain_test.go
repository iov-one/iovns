package controllers

import (
	"errors"
	"math"
	"testing"

	"github.com/iov-one/iovns/mock"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
	"github.com/stretchr/testify/assert"
)

func TestDomain_requireDomain(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		k, ctx, _ := keeper.NewTestKeeper(t, true)
		alice, _ := mock.Addresses()
		k.CreateDomain(ctx, types.Domain{
			Name:         "test",
			Admin:        alice,
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
		alice, _ := mock.Addresses()
		k.CreateDomain(ctx, types.Domain{
			Name:         "test",
			Admin:        alice,
			HasSuperuser: false,
			ValidUntil:   0,
		})
		ctrl := NewDomainController(ctx, k, "test")
		expired := ctrl.domainExpired()
		if expired != true {
			t.Fatal("validation failed: domain has not expired")
		}
	})
	t.Run("domain not expired", func(t *testing.T) {
		k, ctx, _ := keeper.NewTestKeeper(t, true)
		alice, _ := mock.Addresses()
		k.CreateDomain(ctx, types.Domain{
			Name:         "test",
			Admin:        alice,
			HasSuperuser: false,
			ValidUntil:   math.MaxInt64,
		})
		ctrl := NewDomainController(ctx, k, "test")
		expired := ctrl.domainExpired()
		if expired != false {
			t.Fatal("validation failed domain has expired")
		}
	})
	t.Run("domain does not exist", func(t *testing.T) {
		k, ctx, _ := keeper.NewTestKeeper(t, true)
		ctrl := NewDomainController(ctx, k, "test")
		assert.Panics(t, func() { _ = ctrl.domainExpired() }, "domain does not exists")
	})
}

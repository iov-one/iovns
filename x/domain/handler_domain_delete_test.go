package domain

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovnsd"
	"github.com/iov-one/iovnsd/x/domain/keeper"
	"github.com/iov-one/iovnsd/x/domain/types"
	"testing"
)

func Test_handleMsgDomainDelete(t *testing.T) {
	cases := map[string]subTest{
		"fail domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// don't do anything
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handleMsgDomainDelete(ctx, k, types.MsgDeleteDomain{
					Domain: "this does not exist",
					Owner:  bobKey.GetAddress(),
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handleMsgDomainDelete() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
		},
		"fail domain has no superuser": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// set domain with no superuser
				k.SetDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        nil,
					ValidUntil:   0,
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handleMsgDomainDelete(ctx, k, types.MsgDeleteDomain{
					Domain: "test",
					Owner:  bobKey.GetAddress(),
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handleMsgDomainDelete() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"fail domain admin does not match msg owner": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				k.SetDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        bobKey.GetAddress(),
					ValidUntil:   0,
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handleMsgDomainDelete(ctx, k, types.MsgDeleteDomain{
					Domain: "test",
					Owner:  aliceKey.GetAddress(),
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handleMsgDomainDelete() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// set domain
				k.SetDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        aliceKey.GetAddress(),
					ValidUntil:   0,
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
				// add two accounts
				k.SetAccount(ctx, types.Account{
					Domain: "test",
					Name:   "1",
				})
				// add two accounts
				k.SetAccount(ctx, types.Account{
					Domain: "test",
					Name:   "2",
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handleMsgDomainDelete(ctx, k, types.MsgDeleteDomain{
					Domain: "test",
					Owner:  aliceKey.GetAddress(),
				})
				if err != nil {
					t.Fatalf("handleMsgDomainDelete() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, exists := k.GetDomain(ctx, "test")
				if exists {
					t.Fatalf("handleMsgDomainDelete() domain should not exist")
				}
				_, exists = k.GetAccount(ctx, iovnsd.GetAccountKey("test", "1"))
				if exists {
					t.Fatalf("handleMsgDomainDelete() account 1 should not exist")
				}
				_, exists = k.GetAccount(ctx, iovnsd.GetAccountKey("test", "2"))
				if exists {
					t.Fatalf("handleMsgDomainDelete() account 2 should not exist")
				}
			},
		},
	}
	runTests(t, cases)
}

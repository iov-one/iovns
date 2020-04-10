package domain

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovnsd"
	"github.com/iov-one/iovnsd/x/domain/keeper"
	"github.com/iov-one/iovnsd/x/domain/types"
	"testing"
)

func Test_handlerMsgDeleteAccount(t *testing.T) {
	cases := map[string]subTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {

			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgDeleteAccount(ctx, k, types.MsgDeleteAccount{
					Domain: "",
					Name:   "",
					Owner:  nil,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgDeleteAccount() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				k.SetDomain(ctx, types.Domain{
					Name: "test",
				})

			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgDeleteAccount(ctx, k, types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  nil,
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgDeleteAccount() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"msg owner does not own domain or account": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				k.SetDomain(ctx, types.Domain{
					Name:  "test",
					Admin: aliceKey.GetAddress(),
				})
				k.SetAccount(ctx, types.Account{
					Domain: "test",
					Name:   "test",
					Owner:  aliceKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgDeleteAccount(ctx, k, types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  bobKey.GetAddress(),
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgDeleteAccount() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success domain owner": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				k.SetDomain(ctx, types.Domain{
					Name:  "test",
					Admin: aliceKey.GetAddress(),
				})
				k.SetAccount(ctx, types.Account{
					Domain: "test",
					Name:   "test",
					Owner:  bobKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgDeleteAccount(ctx, k, types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  aliceKey.GetAddress(),
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, exists := k.GetAccount(ctx, iovnsd.GetAccountKey("test", "test"))
				if exists {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
		"success account owner": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				k.SetDomain(ctx, types.Domain{
					Name:  "test",
					Admin: aliceKey.GetAddress(),
				})
				k.SetAccount(ctx, types.Account{
					Domain: "test",
					Name:   "test",
					Owner:  bobKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgDeleteAccount(ctx, k, types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  bobKey.GetAddress(),
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, exists := k.GetAccount(ctx, iovnsd.GetAccountKey("test", "test"))
				if exists {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
	}

	// run tests
	runTests(t, cases)
}

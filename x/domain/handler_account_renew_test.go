package domain

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
	"testing"
)

func Test_handlerMsgRenewAccount(t *testing.T) {
	cases := map[string]subTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {

			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgRenewAccount(ctx, k, types.MsgRenewAccount{
					Domain: "does not exist",
					Name:   "",
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgRenewAccount() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// set mock domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					AccountRenew: 100,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgRenewAccount(ctx, k, types.MsgRenewAccount{
					Domain: "test",
					Name:   "does not exist",
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgRenewAccount() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// set mock domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					AccountRenew: 100,
				})
				// set mock account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: 1000,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgRenewAccount(ctx, k, types.MsgRenewAccount{
					Domain: "test",
					Name:   "test",
				})
				if err != nil {
					t.Fatalf("handlerMsgRenewAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				account, _ := k.GetAccount(ctx, "test", "test")
				if account.ValidUntil != 1100 {
					t.Fatalf("handlerMsgRenewAccount() expected 1100, got: %d", account.ValidUntil)
				}
			},
		},
	}

	runTests(t, cases)
}

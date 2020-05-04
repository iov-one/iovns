package domain

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
	"testing"
)

func Test_handlerMsgFlushDomain(t *testing.T) {
	cases := map[string]subTest{
		"domain does not exist": {
			BeforeTest: nil,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgFlushDomain(ctx, k, &types.MsgFlushDomain{
					Domain: "does not exist",
					Owner:  nil,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgFlushDomain() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain has superuser": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        nil,
					ValidUntil:   0,
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgFlushDomain(ctx, k, &types.MsgFlushDomain{
					Domain: "test",
					Owner:  nil,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgFlushDomain() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"msg owner is not domain admin": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        aliceKey.GetAddress(),
					ValidUntil:   0,
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgFlushDomain(ctx, k, &types.MsgFlushDomain{
					Domain: "test",
					Owner:  bobKey.GetAddress(),
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgFlushDomain() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        aliceKey.GetAddress(),
					ValidUntil:   0,
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
				// add empty account 1
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "",
				})
				// add account 2
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "1",
				})
				// add account 2
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "2",
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgFlushDomain(ctx, k, &types.MsgFlushDomain{
					Domain: "test",
					Owner:  aliceKey.GetAddress(),
				})
				if err != nil {
					t.Fatalf("handlerMsgFlushDomain() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				var exists bool
				_, exists = k.GetAccount(ctx, "test", "")
				if !exists {
					t.Fatalf("handlerMsgFlushDomain() empty account was deleted")
				}
				_, exists = k.GetAccount(ctx, "test", "1")
				if exists {
					t.Fatalf("handlerMsgFlushDomain() account 1 was not deleted")
				}
				_, exists = k.GetAccount(ctx, "test", "2")
				if exists {
					t.Fatalf("handlerMsgFlushDomain() account 2 was not deleted")
				}
			},
		},
	}
	runTests(t, cases)
}

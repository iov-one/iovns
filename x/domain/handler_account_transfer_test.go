package domain

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
	"testing"
	"time"
)

func Test_handlerAccountTransfer(t *testing.T) {
	testCases := map[string]subTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// do nothing
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, types.MsgTransferAccount{
					Domain:   "does not exist",
					Name:     "does not exist",
					Owner:    nil,
					NewOwner: nil,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain has expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:         "expired domain",
					Admin:        nil,
					ValidUntil:   0,
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, types.MsgTransferAccount{
					Domain:   "expired domain",
					Name:     "",
					Owner:    nil,
					NewOwner: nil,
				})
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        nil,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, types.MsgTransferAccount{
					Domain:   "test",
					Name:     "this account does not exist",
					Owner:    nil,
					NewOwner: nil,
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"account expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        nil,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
				k.CreateAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        nil,
					ValidUntil:   0,
					Targets:      nil,
					Certificates: nil,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    nil,
					NewOwner: nil,
				})
				if !errors.Is(err, types.ErrAccountExpired) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types.ErrAccountExpired, err)
				}
			},
			AfterTest: nil,
		},
		"if domain has super user only domain admin can transfer accounts": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        aliceKey.GetAddress(),
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
				k.CreateAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        nil,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Targets:      nil,
					Certificates: nil,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    bobKey.GetAddress(),
					NewOwner: nil,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"if domain has no super user then only account owner can transfer accounts": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        nil,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
				k.CreateAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        aliceKey.GetAddress(),
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Targets:      nil,
					Certificates: nil,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    bobKey.GetAddress(),
					NewOwner: nil,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        nil,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
				k.CreateAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        aliceKey.GetAddress(),
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Targets:      nil,
					Certificates: nil,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    aliceKey.GetAddress(),
					NewOwner: nil,
				})
				if err != nil {
					t.Fatalf("handlerMsgTransferAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				account, exists := k.GetAccount(ctx, "test", "test")
				if !exists {
					panic("unexpected account deletion")
				}
				if account.Targets != nil {
					t.Fatalf("handlerAccountTransfer() account targets were not deleted")
				}
				if account.Certificates != nil {
					t.Fatalf("handlerAccountTransfer() account certificates were not deleted")
				}
				if !account.Owner.Equals(aliceKey.GetAddress()) {
					t.Fatalf("handlerAccounTransfer() expected new owner: %s, got: %s", aliceKey.GetAddress(), account.Owner)
				}
			},
		},
	}
	runTests(t, testCases)
}

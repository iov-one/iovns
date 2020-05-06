package domain

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
	"reflect"
	"testing"
	"time"
)

func Test_handlerMsgSetAccountMetadata(t *testing.T) {
	cases := map[string]subTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgSetAccountMetadata(ctx, k, &types.MsgSetAccountMetadata{
					Domain: "does not exist",
					Name:   "",
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgSetAccountMetadata() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// create domain
				k.CreateDomain(ctx, types.Domain{
					Name: "test",
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgSetAccountMetadata(ctx, k, &types.MsgSetAccountMetadata{
					Domain:         "test",
					Name:           "",
					NewMetadataURI: "",
					Owner:          nil,
				})
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("handlerMsgSetAccountMetadata() expected error: %s, got: %s", types.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// create domain
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgSetAccountMetadata(ctx, k, &types.MsgSetAccountMetadata{
					Domain: "test",
					Name:   "does not exist",
					Owner:  nil,
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgSetAccountMetadata() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"account expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// create domain
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
				})
				// create account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: 0,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgSetAccountMetadata(ctx, k, &types.MsgSetAccountMetadata{
					Domain: "test",
					Name:   "test",
					Owner:  nil,
				})
				if !errors.Is(err, types.ErrAccountExpired) {
					t.Fatalf("handlerMsgSetAccountMetadata() expected error: %s, got: %s", types.ErrAccountExpired, err)
				}
			},
			AfterTest: nil,
		},
		"signer is not owner of account": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// create domain
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
				})
				// create account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      aliceKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgSetAccountMetadata(ctx, k, &types.MsgSetAccountMetadata{
					Domain: "test",
					Name:   "test",
					Owner:  bobKey.GetAddress(),
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgSetAccountMetadata() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// create domain
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
				})
				// create account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      aliceKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgSetAccountMetadata(ctx, k, &types.MsgSetAccountMetadata{
					Domain:         "test",
					Name:           "test",
					NewMetadataURI: "https://test.com",
					Owner:          aliceKey.GetAddress(),
				})
				if err != nil {
					t.Fatalf("handlerMsgSetAccountMetadata() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				expected := "https://test.com"
				account, _ := k.GetAccount(ctx, "test", "test")
				if !reflect.DeepEqual(expected, account.MetadataURI) {
					t.Fatalf("handlerMsgSetMetadataURI expected: %+v, got %+v", expected, account.MetadataURI)
				}
			},
		},
	}
	// run tests
	runTests(t, cases)
}

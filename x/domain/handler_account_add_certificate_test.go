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

func Test_handlerMsgAddAccountCertificates(t *testing.T) {
	cases := map[string]subTest{
		"domain does not exist": {
			BeforeTest: nil,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "does not exist",
					Name:           "",
					Owner:          nil,
					NewCertificate: nil,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// add expired domain
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: 0,
				})
				//
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "",
					Owner:          nil,
					NewCertificate: nil,
				})
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				})
				//
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "does not exist",
					Owner:          nil,
					NewCertificate: nil,
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"account expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				})
				// add mock account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: 0,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          nil,
					NewCertificate: nil,
				})
				if !errors.Is(err, types.ErrAccountExpired) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types.ErrAccountExpired, err)
				}
			},
			AfterTest: nil,
		},
		"msg owner is not account owner": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				})
				// add mock account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      aliceKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          bobKey.GetAddress(),
					NewCertificate: nil,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"certificate exists": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				})
				// add mock account
				k.CreateAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "test",
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:        aliceKey.GetAddress(),
					Certificates: [][]byte{[]byte("test")},
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          aliceKey.GetAddress(),
					NewCertificate: []byte("test"),
				})
				if !errors.Is(err, types.ErrCertificateExists) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types.ErrCertificateExists, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				})
				// add mock account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      aliceKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          aliceKey.GetAddress(),
					NewCertificate: []byte("test"),
				})
				if err != nil {
					t.Fatalf("handlerMsgAddAccountCertificates() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				expected := [][]byte{[]byte("test")}
				account, _ := k.GetAccount(ctx, "test", "test")
				if !reflect.DeepEqual(account.Certificates, expected) {
					t.Fatalf("handlerMsgAddAccountCertificates: got: %#v, expected: %#v", account.Certificates, expected)
				}
			},
		},
	}
	runTests(t, cases)
}

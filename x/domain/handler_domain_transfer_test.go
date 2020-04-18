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

func Test_handlerMsgTransferDomain(t *testing.T) {
	cases := map[string]subTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {

			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgTransferDomain(ctx, k, types.MsgTransferDomain{
					Domain:   "does not exist",
					Owner:    nil,
					NewAdmin: nil,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgTransferDomain() expected error: %s, got error: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain has no superuser": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				k.SetDomain(ctx, types.Domain{
					Name:         "test",
					HasSuperuser: false,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgTransferDomain(ctx, k, types.MsgTransferDomain{
					Domain:   "test",
					Owner:    nil,
					NewAdmin: nil,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgTransferDomain() expected error: %s, got error: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"domain has expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				k.SetDomain(ctx, types.Domain{
					Name:         "test",
					HasSuperuser: true,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgTransferDomain(ctx, k, types.MsgTransferDomain{
					Domain:   "test",
					Owner:    nil,
					NewAdmin: nil,
				})
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("handlerMsgTransferDomain() expected error: %s, got error: %s", types.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"msg signer is not domain admin": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				k.SetDomain(ctx, types.Domain{
					Name:         "test",
					HasSuperuser: true,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:        aliceKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgTransferDomain(ctx, k, types.MsgTransferDomain{
					Domain:   "test",
					Owner:    bobKey.GetAddress(),
					NewAdmin: nil,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgTransferDomain() expected error: %s, got error: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// create domain
				k.SetDomain(ctx, types.Domain{
					Name:         "test",
					HasSuperuser: true,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:        aliceKey.GetAddress(),
				})
				// add empty account
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "",
				})
				// add account 1
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "1",
					Owner:      aliceKey.GetAddress(),
					ValidUntil: 0,
					Targets: []iovns.BlockchainAddress{{
						ID:      "test",
						Address: "test",
					}},
					Certificates: [][]byte{[]byte("cert")},
					Broker:       nil,
				})
				// add account 2
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "2",
					Owner:      aliceKey.GetAddress(),
					ValidUntil: 0,
					Targets: []iovns.BlockchainAddress{{
						ID:      "test",
						Address: "test",
					}},
					Certificates: [][]byte{[]byte("cert")},
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgTransferDomain(ctx, k, types.MsgTransferDomain{
					Domain:   "test",
					Owner:    aliceKey.GetAddress(),
					NewAdmin: bobKey.GetAddress(),
				})
				if err != nil {
					t.Fatalf("handlerMsgTransferDomain() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// check domain new owner
				domain, _ := k.GetDomain(ctx, "test")
				if !bobKey.GetAddress().Equals(domain.Admin) {
					t.Fatalf("handlerMsgTransferDomain() expected domain owner: %s, got: %s", bobKey.GetAddress(), domain.Admin)
				}
				// check if account new owner has changed
				account, _ := k.GetAccount(ctx, "test", "1")
				if !account.Owner.Equals(bobKey.GetAddress()) {
					t.Fatalf("handlerMsgTransferDomain() expected account owner: %s, got: %s", bobKey.GetAddress(), account.Owner)
				}
				// check if targets deleted
				if account.Targets != nil {
					t.Fatalf("handlerMsgTransferDomain expected account targets: <nil>, got: %#v", account.Targets)
				}
				// check if certs deleted
				if account.Certificates != nil {
					t.Fatalf("handlerMsgTransferDomain expected account certificates: <nil>, got: %#v", account.Certificates)
				}
				// check no changes in empty account
				if emptyAcc, _ := k.GetAccount(ctx, "test", ""); !reflect.DeepEqual(emptyAcc, types.Account{Domain: "test", Name: ""}) {
					t.Fatalf("handlerMsgTransferdomain() empty account mismatch, expected: %+v, got: %+v", types.Account{Domain: "test", Name: ""}, emptyAcc)
				}
			},
		},
	}

	runTests(t, cases)
}

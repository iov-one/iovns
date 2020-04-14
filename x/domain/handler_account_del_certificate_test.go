package domain

import (
	"bytes"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
	"testing"
)

func Test_handlerMsgDeleteAccountCertificate(t *testing.T) {
	cases := map[string]subTest{
		"account does not exist": {
			BeforeTest: nil,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "does not exist",
					DeleteCertificate: nil,
					Owner:             nil,
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgDeleteAccountCertificate() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"msg signer is not account owner": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				k.SetAccount(ctx, types.Account{
					Domain: "test",
					Name:   "test",
					Owner:  aliceKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: nil,
					Owner:             bobKey.GetAddress(),
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgDeleteAccountCertificate() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"certificate does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				k.SetAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        aliceKey.GetAddress(),
					Certificates: nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: []byte("does not exist"),
					Owner:             aliceKey.GetAddress(),
				})
				if !errors.Is(err, types.ErrCertificateDoesNotExist) {
					t.Fatalf("handlerMsgDeleteAccountCertificate() expected error: %s, got: %s", types.ErrCertificateDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				k.SetAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        aliceKey.GetAddress(),
					Certificates: [][]byte{[]byte("test")},
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: []byte("test"),
					Owner:             aliceKey.GetAddress(),
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccountCertificates() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// check if certificate is still present
				account, _ := k.GetAccount(ctx, iovns.GetAccountKey("test", "test"))
				for _, cert := range account.Certificates {
					if bytes.Equal(cert, []byte("test")) {
						t.Fatalf("handlerMsgDeleteAccountCertificates() certificate not deleted")
					}
				}
				// success
			},
		},
	}

	runTests(t, cases)
}

package domain

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

func Test_handlerMsgAddAccountCertificates(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"domain does not exist": {
			BeforeTest: nil,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "does not exist",
					Name:           "",
					Owner:          keeper.BobKey,
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
					Admin:      keeper.BobKey,
				})
				//
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "",
					Owner:          keeper.BobKey,
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
					Admin:      keeper.BobKey,
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
					Admin:      keeper.BobKey,
				})
				// add mock account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: 0,
					Owner:      keeper.AliceKey,
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
					Admin:      keeper.AliceKey,
				})
				// add mock account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          keeper.BobKey,
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
					Admin:      keeper.BobKey,
				})
				// add mock account
				k.CreateAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "test",
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:        keeper.AliceKey,
					Certificates: []types.Certificate{[]byte("test")},
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          keeper.AliceKey,
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
					Admin:      keeper.AliceKey,
				})
				// add mock account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          keeper.AliceKey,
					NewCertificate: []byte("test"),
				})
				if err != nil {
					t.Fatalf("handlerMsgAddAccountCertificates() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				expected := []types.Certificate{[]byte("test")}
				account, _ := k.GetAccount(ctx, "test", "test")
				if !reflect.DeepEqual(account.Certificates, expected) {
					t.Fatalf("handlerMsgAddAccountCertificates: got: %#v, expected: %#v", account.Certificates, expected)
				}
			},
		},
	}
	keeper.RunTests(t, cases)
}

func Test_handlerMsgDeleteAccountCertificate(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"account does not exist": {
			BeforeTest: nil,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "does not exist",
					DeleteCertificate: nil,
					Owner:             keeper.BobKey,
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgDeleteAccountCertificate() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"msg signer is not account owner": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: nil,
					Owner:             keeper.BobKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgDeleteAccountCertificate() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"certificate does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        keeper.AliceKey,
					Certificates: nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: []byte("does not exist"),
					Owner:             keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrCertificateDoesNotExist) {
					t.Fatalf("handlerMsgDeleteAccountCertificate() expected error: %s, got: %s", types.ErrCertificateDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        keeper.AliceKey,
					Certificates: []types.Certificate{[]byte("test")},
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: []byte("test"),
					Owner:             keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccountCertificates() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// check if certificate is still present
				account, _ := k.GetAccount(ctx, "test", "test")
				for _, cert := range account.Certificates {
					if bytes.Equal(cert, []byte("test")) {
						t.Fatalf("handlerMsgDeleteAccountCertificates() certificate not deleted")
					}
				}
				// success
			},
		},
	}

	keeper.RunTests(t, cases)
}

func Test_Closed_handlerMsgDeleteAccount(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"domain admin can delete": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					Type:       types.ClosedDomain,
					ValidUntil: types.MaxValidUntil,
				})
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, exists := k.GetAccount(ctx, "test", "test")
				if exists {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
		"domain expired": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					Type:       types.ClosedDomain,
					ValidUntil: 2,
				})
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
				})
			},
			TestBlockTime: 2,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
		},
		"account owner cannot delete": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					Type:       types.ClosedDomain,
					ValidUntil: types.MaxValidUntil,
				})
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("unexpected error: %s", err)
				}
			},
		},
	}
	keeper.RunTests(t, cases)
}

func Test_Open_handlerMsgDeleteAccount(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"domain admin cannot can delete before grace period": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: keeper.RegexMatchNothing,
					ValidBlockchainID:      keeper.RegexMatchAll,
					AccountGracePeriod:     1000 * time.Second,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					Type:       types.OpenDomain,
					ValidUntil: types.MaxValidUntil,
				})
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
				})
			},
			TestBlockTime: 3,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
		},
		"no domain valid until check": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: keeper.RegexMatchNothing,
					ValidBlockchainID:      keeper.RegexMatchAll,
					DomainRenewalPeriod:    10,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					Type:       types.OpenDomain,
					ValidUntil: 2,
				})
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
				})
			},
			TestBlockTime: 100,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, exists := k.GetAccount(ctx, "test", "test")
				if exists {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
		"only account owner can delete before grace period": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: keeper.RegexMatchNothing,
					ValidBlockchainID:      keeper.RegexMatchAll,
					AccountGracePeriod:     10 * time.Second,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					Type:       types.OpenDomain,
					ValidUntil: types.MaxValidUntil,
				})
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey,
				})
			},
			TestBlockTime: 5,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// admin test
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("unexpected error: %v", err)
				}
				// anyone test
				_, err = handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.CharlieKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("unexpected error: %v", err)
				}
				// account owner test
				_, err = handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, exists := k.GetAccount(ctx, "test", "test")
				if exists {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
		"domain admin can delete after grace": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: keeper.RegexMatchNothing,
					ValidBlockchainID:      keeper.RegexMatchAll,
					AccountGracePeriod:     10,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					Type:       types.OpenDomain,
					ValidUntil: types.MaxValidUntil,
				})
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey,
				})
			},
			TestBlockTime: 100,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// admin test
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
				})
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, exists := k.GetAccount(ctx, "test", "test")
				if exists {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
		"anyone can delete after grace": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: keeper.RegexMatchNothing,
					ValidBlockchainID:      keeper.RegexMatchAll,
					AccountGracePeriod:     10,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					Type:       types.OpenDomain,
					ValidUntil: types.MaxValidUntil,
				})
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey,
				})
			},
			TestBlockTime: 100,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// admin test
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.CharlieKey,
				})
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, exists := k.GetAccount(ctx, "test", "test")
				if exists {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
	}
	keeper.RunTests(t, cases)
}

func Test_Common_handlerMsgDeleteAccount(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {

			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "does not exist",
					Name:   "does not exist",
					Owner:  keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgDeleteAccount() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:  "test",
					Admin: keeper.BobKey,
				})

			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
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
		"success domain owner": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:  "test",
					Admin: keeper.AliceKey,
				})
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, exists := k.GetAccount(ctx, "test", "test")
				if exists {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
		"success account owner": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:  "test",
					Admin: keeper.AliceKey,
				})
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, exists := k.GetAccount(ctx, "test", "test")
				if exists {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
	}

	// run tests
	keeper.RunTests(t, cases)
}

func Test_ClosedDomain_handlerMsgRegisterAccount(t *testing.T) {
	testCases := map[string]keeper.SubTest{
		"only domain admin can register account": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: keeper.RegexMatchNothing,
					ValidBlockchainID:      keeper.RegexMatchAll,
					DomainRenewalPeriod:    10,
					AccountRenewalPeriod:   10,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        keeper.BobKey,
					ValidUntil:   time.Now().Add(100000 * time.Hour).Unix(),
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test",
					Owner:      keeper.BobKey,
					Registerer: keeper.BobKey,
				})
				if err != nil {
					t.Fatalf("handlerRegisterAccount() got error: %s", err)
				}
				_, err = handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test2",
					Owner:      keeper.BobKey,
					Registerer: keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerRegisterAccount() got error: %s", err)
				}
			},
		},
		"account valid until is set to max": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in blockchain targets
				// get set config function
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: keeper.RegexMatchNothing, // don't match anything
					ValidBlockchainID:      keeper.RegexMatchAll,     // match all
					DomainRenewalPeriod:    10,
				})
				// add a closed domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        keeper.BobKey,
					ValidUntil:   time.Now().Add(100000 * time.Hour).Unix(),
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test",
					Owner:      keeper.BobKey,
					Registerer: keeper.BobKey,
				})
				if err != nil {
					t.Fatalf("handlerRegisterAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				acc, ok := k.GetAccount(ctx, "test", "test")
				if !ok {
					t.Fatal("account test not found")
				}
				if acc.ValidUntil != types.MaxValidUntil {
					t.Fatalf("unexpected account valid until %d", acc.ValidUntil)
				}
			},
		},
		"account owner can be different than domain admin": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in blockchain targets
				// get set config function
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: keeper.RegexMatchNothing, // don't match anything
					ValidBlockchainID:      keeper.RegexMatchAll,     // match all
					DomainRenewalPeriod:    10,
				})
				// add a closed domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        keeper.BobKey,
					ValidUntil:   time.Now().Add(100000 * time.Hour).Unix(),
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test",
					Registerer: keeper.BobKey,
					Owner:      keeper.BobKey,
				})
				if err != nil {
					t.Fatalf("handlerRegisterAccount() got error: %s", err)
				}
				_, err = handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test2",
					Registerer: keeper.BobKey,
					Owner:      keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerRegisterAccount() got error: %s", err)
				}
			},
		},
	}
	// run tests
	keeper.RunTests(t, testCases)
}

func Test_OpenDomain_handleMsgRegisterAccount(t *testing.T) {
	testCases := map[string]keeper.SubTest{
		"account valid until is now plus config account renew": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in blockchain targets
				// get set config function
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: keeper.RegexMatchNothing, // don't match anything
					ValidBlockchainID:      keeper.RegexMatchAll,     // match all
					DomainRenewalPeriod:    10 * time.Second,
					AccountRenewalPeriod:   10 * time.Second,
				})
				// add a closed domain
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					ValidUntil: time.Now().Add(100000 * time.Hour).Unix(),
					Type:       types.OpenDomain,
					Broker:     nil,
				})
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test",
					Owner:      keeper.BobKey,
					Registerer: keeper.BobKey,
				})
				if err != nil {
					t.Fatalf("handlerRegisterAccount() got error: %s", err)
				}
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				acc, ok := k.GetAccount(ctx, "test", "test")
				if !ok {
					t.Fatal("account test not found")
				}
				expected := iovns.TimeToSeconds(time.Unix(11, 0))
				if acc.ValidUntil != expected {
					t.Fatalf("unexpected account valid until %d, expected %d", acc.ValidUntil, expected)
				}
			},
		},
	}
	keeper.RunTests(t, testCases)
}

func Test_Common_handleMsgRegisterAccount(t *testing.T) {
	testCases := map[string]keeper.SubTest{
		"fail invalid blockchain targets address": {
			TestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in blockchain targets
				// get set config function
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: keeper.RegexMatchNothing, // don't match anything
					ValidBlockchainID:      keeper.RegexMatchAll,     // match all
					DomainRenewalPeriod:    10,
				})
				// add a domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        keeper.BobKey,
					ValidUntil:   2,
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
					Targets: []types.BlockchainAddress{
						{
							ID:      "works",
							Address: "won't work",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrInvalidBlockchainTarget) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrInvalidBlockchainTarget, err)
				}
			},
			AfterTest: nil,
		},
		// TODO cleanup comments
		"fail invalid blockchain targets id": {
			TestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in blockchain targets
				// get set config function
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      keeper.RegexMatchNothing, // don't match anything
					ValidBlockchainAddress: keeper.RegexMatchAll,     // match all
					DomainRenewalPeriod:    10,
				})
				// add a domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        keeper.BobKey,
					ValidUntil:   2,
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
					Targets: []types.BlockchainAddress{
						{
							ID:      "invalid blockchain id",
							Address: "valid blockchain address",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrInvalidBlockchainTarget) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrInvalidBlockchainTarget, err)
				}
			},
			AfterTest: nil,
		},
		"fail invalid account name": {
			TestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set config
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: keeper.RegexMatchAll,     // match all
					ValidBlockchainID:      keeper.RegexMatchAll,     // match all
					ValidAccountName:       keeper.RegexMatchNothing, // match nothing
					DomainRenewalPeriod:    10,
				})
				// add a domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        keeper.AliceKey,
					ValidUntil:   2,
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "this won't match",
					Owner:  keeper.AliceKey,
					Targets: []types.BlockchainAddress{
						{
							ID:      "works",
							Address: "works",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrInvalidAccountName) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrInvalidAccountName, err)
				}
			},
			AfterTest: nil,
		},
		"fail domain name does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in blockchain targets
				// get set config function
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: keeper.RegexMatchAll, // match all
					ValidBlockchainID:      keeper.RegexMatchAll, // match all
					ValidAccountName:       keeper.RegexMatchAll, // match nothing
					DomainRenewalPeriod:    10,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "this does not exist",
					Name:   "works",
					Owner:  keeper.AliceKey,
					Targets: []types.BlockchainAddress{
						{
							ID:      "works",
							Address: "works",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"fail only owner of domain with superuser can register accounts": {
			TestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in blockchain targets
				// get set config function
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: keeper.RegexMatchAll, // match all
					ValidBlockchainID:      keeper.RegexMatchAll, // match all
					ValidAccountName:       keeper.RegexMatchAll, // match nothing
					DomainRenewalPeriod:    10,
				})
				// add a domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        keeper.BobKey,
					ValidUntil:   2,
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey, // invalid owner
					Targets: []types.BlockchainAddress{
						{
							ID:      "works",
							Address: "works",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"fail domain has expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in blockchain targets
				// get set config function
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: keeper.RegexMatchAll, // match all
					ValidBlockchainID:      keeper.RegexMatchAll, // match all
					ValidAccountName:       keeper.RegexMatchAll, // match nothing
					DomainRenewalPeriod:    10,
				})
				// add a domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        keeper.BobKey,
					ValidUntil:   0, // domain is expired
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
					Targets: []types.BlockchainAddress{
						{
							ID:      "works",
							Address: "works",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"fail account exists": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in blockchain targets
				// get set config function
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: keeper.RegexMatchAll, // match all
					ValidBlockchainID:      keeper.RegexMatchAll, // match all
					ValidAccountName:       keeper.RegexMatchAll, // match nothing
					DomainRenewalPeriod:    10,
				})
				// add a domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        keeper.BobKey,
					ValidUntil:   time.Now().Add(100000 * time.Hour).Unix(),
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
				// add an account that we are gonna try to overwrite
				k.CreateAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "exists",
					Owner:        keeper.AliceKey,
					ValidUntil:   0,
					Targets:      nil,
					Certificates: nil,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "exists",
					Owner:  keeper.BobKey,
					Targets: []types.BlockchainAddress{
						{
							ID:      "works",
							Address: "works",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrAccountExists) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrAccountExists, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in blockchain targets
				// get set config function
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: keeper.RegexMatchAll, // match all
					ValidBlockchainID:      keeper.RegexMatchAll, // match all
					ValidAccountName:       keeper.RegexMatchAll, // match nothing
					DomainRenewalPeriod:    10,
				})
				// add a domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        keeper.BobKey,
					ValidUntil:   time.Now().Add(100000 * time.Hour).Unix(),
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test",
					Owner:      keeper.BobKey,
					Registerer: keeper.BobKey,
					Targets: []types.BlockchainAddress{
						{
							ID:      "works",
							Address: "works",
						},
					},
					Broker: nil,
				})
				if err != nil {
					t.Fatalf("handleMsgRegisterAccount() got error: %s", err)
				}
			},
			AfterTest: nil, // TODO fill with matching data
		},
	}
	// run tests
	keeper.RunTests(t, testCases)
}

func Test_handlerMsgRenewAccount(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {

			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgRenewAccount(ctx, k, &types.MsgRenewAccount{
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
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set mock domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					AccountRenew: 100,
					Admin:        keeper.BobKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgRenewAccount(ctx, k, &types.MsgRenewAccount{
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
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set mock domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					AccountRenew: 1000 * time.Second,
					Admin:        keeper.BobKey,
				})
				// set mock account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Unix(1000, 0)),
					Owner:      keeper.AliceKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgRenewAccount(ctx, k, &types.MsgRenewAccount{
					Domain: "test",
					Name:   "test",
				})
				if err != nil {
					t.Fatalf("handlerMsgRenewAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				account, _ := k.GetAccount(ctx, "test", "test")
				want := iovns.TimeToSeconds(time.Unix(1000, 0).Add(1000 * time.Second))
				if account.ValidUntil != want {
					t.Fatalf("handlerMsgRenewAccount() want: %d, got: %d", want, account.ValidUntil)
				}
			},
		},
	}

	keeper.RunTests(t, cases)
}

func Test_handlerMsgReplaceAccountTargets(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"invalid blockchain target": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      keeper.RegexMatchNothing,
					ValidBlockchainAddress: keeper.RegexMatchNothing,
				})
				// create domain
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
				})
				// create account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, &types.MsgReplaceAccountTargets{
					Domain: "test",
					Name:   "test",
					NewTargets: []types.BlockchainAddress{
						{
							ID:      "invalid",
							Address: "invalid",
						},
					},
					Owner: keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrInvalidBlockchainTarget) {
					t.Fatalf("handlerMsgReplaceAccountTargets() expected error: %s, got: %s", types.ErrInvalidBlockchainTarget, err)
				}
			},
		},
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      keeper.RegexMatchAll,
					ValidBlockchainAddress: keeper.RegexMatchAll,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, &types.MsgReplaceAccountTargets{
					Domain: "does not exist",
					Name:   "",
					NewTargets: []types.BlockchainAddress{
						{
							ID:      "valid",
							Address: "valid",
						},
					},
					Owner: nil,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgReplaceAccountTargets() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      keeper.RegexMatchAll,
					ValidBlockchainAddress: keeper.RegexMatchAll,
				})
				// create domain
				k.CreateDomain(ctx, types.Domain{
					Name:  "test",
					Admin: keeper.BobKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, &types.MsgReplaceAccountTargets{
					Domain: "test",
					NewTargets: []types.BlockchainAddress{
						{
							ID:      "valid",
							Address: "valid",
						},
					},
					Owner: nil,
				})
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("handlerMsgReplaceAccountTargets() expected error: %s, got: %s", types.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      keeper.RegexMatchAll,
					ValidBlockchainAddress: keeper.RegexMatchAll,
				})
				// create domain
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, &types.MsgReplaceAccountTargets{
					Domain: "test",
					Name:   "does not exist",
					NewTargets: []types.BlockchainAddress{
						{
							ID:      "valid",
							Address: "valid",
						},
					},
					Owner: nil,
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgReplaceAccountTargets() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"account expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      keeper.RegexMatchAll,
					ValidBlockchainAddress: keeper.RegexMatchAll,
				})
				// create domain
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
					Type:       types.OpenDomain,
				})
				// create account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: 0,
					Owner:      keeper.AliceKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, &types.MsgReplaceAccountTargets{
					Domain: "test",
					Name:   "test",
					NewTargets: []types.BlockchainAddress{
						{
							ID:      "valid",
							Address: "valid",
						},
					},
					Owner: nil,
				})
				if !errors.Is(err, types.ErrAccountExpired) {
					t.Fatalf("handlerMsgReplaceAccountTargets() expected error: %s, got: %s", types.ErrAccountExpired, err)
				}
			},
			AfterTest: nil,
		},
		"signer is not owner of account": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      keeper.RegexMatchAll,
					ValidBlockchainAddress: keeper.RegexMatchAll,
				})
				// create domain
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
				})
				// create account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, &types.MsgReplaceAccountTargets{
					Domain: "test",
					Name:   "test",
					NewTargets: []types.BlockchainAddress{
						{
							ID:      "valid",
							Address: "valid",
						},
					},
					Owner: keeper.BobKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgReplaceAccountTargets() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      keeper.RegexMatchAll,
					ValidBlockchainAddress: keeper.RegexMatchAll,
					BlockchainTargetMax:    5,
				})
				// create domain
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
				})
				// create account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, &types.MsgReplaceAccountTargets{
					Domain: "test",
					Name:   "test",
					NewTargets: []types.BlockchainAddress{
						{
							ID:      "valid",
							Address: "valid",
						},
					},
					Owner: keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgReplaceAccountTargets() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				expected := []types.BlockchainAddress{{
					ID:      "valid",
					Address: "valid",
				}}
				account, _ := k.GetAccount(ctx, "test", "test")
				if !reflect.DeepEqual(expected, account.Targets) {
					t.Fatalf("handlerMsgReplaceAccountTargets() expected: %+v, got %+v", expected, account.Targets)
				}
			},
		},
	}
	// run tests
	keeper.RunTests(t, cases)
}

func Test_handlerMsgSetAccountMetadata(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgSetAccountMetadata(ctx, k, &types.MsgSetAccountMetadata{
					Domain: "does not exist",
					Name:   "",
					Owner:  keeper.AliceKey,
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
					Name:  "test",
					Admin: keeper.BobKey,
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
					Admin:      keeper.BobKey,
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
					Admin:      keeper.BobKey,
				})
				// create account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: 0,
					Owner:      keeper.AliceKey,
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
					Admin:      keeper.BobKey,
				})
				// create account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgSetAccountMetadata(ctx, k, &types.MsgSetAccountMetadata{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
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
					Admin:      keeper.BobKey,
				})
				// create account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgSetAccountMetadata(ctx, k, &types.MsgSetAccountMetadata{
					Domain:         "test",
					Name:           "test",
					NewMetadataURI: "https://test.com",
					Owner:          keeper.AliceKey,
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
	keeper.RunTests(t, cases)
}

func Test_Closed_handlerAccountTransfer(t *testing.T) {
	testCases := map[string]keeper.SubTest{
		"only domain admin can transfer": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// domain owned by alice
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        keeper.AliceKey,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
				// account owned by bob
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					Owner:      keeper.BobKey,
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// alice is domain owner and should transfer account owned by bob to alice
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    keeper.AliceKey,
					NewOwner: keeper.CharlieKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgTransferAccount() got error: %s", err)
				}
				_, err = handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    keeper.BobKey,
					NewOwner: keeper.CharlieKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgTransferAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				account, exists := k.GetAccount(ctx, "test", "test")
				if !exists {
					panic("unexpected account deletion")
				}
				if !account.Owner.Equals(keeper.CharlieKey) {
					t.Fatalf("handlerAccounTransfer() expected new owner: %s, got: %s", keeper.CharlieKey, account.Owner)
				}
			},
		},
		"domain admin can reset account content": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// domain owned by alice
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        keeper.AliceKey,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
				// account owned by bob
				k.CreateAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        keeper.BobKey,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					MetadataURI:  "lol",
					Certificates: []types.Certificate{[]byte("test")},
					Targets: []types.BlockchainAddress{
						{
							ID:      "works",
							Address: "works",
						},
					},
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// alice is domain owner and should transfer account owned by bob to alice
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    keeper.AliceKey,
					NewOwner: keeper.CharlieKey,
					Reset:    true,
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
					panic("targets not deleted")
				}
				if account.Certificates != nil {
					panic("certificates not deleted")
				}
				if account.MetadataURI != "" {
					panic("metadata not deleted")
				}
			},
		},
	}

	keeper.RunTests(t, testCases)
}

func Test_Open_handlerAccountTransfer(t *testing.T) {
	testCases := map[string]keeper.SubTest{
		"domain admin cannot transfer": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// domain owned by alice
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        keeper.AliceKey,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:         types.OpenDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
				// account owned by bob
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					Owner:      keeper.BobKey,
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    keeper.AliceKey,
					NewOwner: keeper.CharlieKey,
					Reset:    false,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgTransferAccount() got error: %s", err)
				}

				_, err = handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    keeper.BobKey,
					NewOwner: keeper.CharlieKey,
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
				if !account.Owner.Equals(keeper.CharlieKey) {
					t.Fatalf("handlerAccounTransfer() expected new owner: %s, got: %s", keeper.CharlieKey, account.Owner)
				}
			},
		},
		"domain admin cannot reset account content": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// domain owned by alice
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        keeper.AliceKey,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:         types.OpenDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
				// account owned by bob
				k.CreateAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        keeper.BobKey,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					MetadataURI:  "lol",
					Certificates: []types.Certificate{[]byte("test")},
					Targets: []types.BlockchainAddress{
						{
							ID:      "works",
							Address: "works",
						},
					},
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// alice is domain owner and should transfer account owned by bob to alice
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    keeper.AliceKey,
					NewOwner: keeper.CharlieKey,
					Reset:    true,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgTransferAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				account, exists := k.GetAccount(ctx, "test", "test")
				if !exists {
					t.Fatal("unexpected account deletion")
				}
				if account.Targets == nil {
					t.Fatal("targets deleted")
				}
				if account.Certificates == nil {
					t.Fatal("certificates deleted")
				}
				if account.MetadataURI == "" {
					t.Fatal("metadata not deleted")
				}
			},
		},
	}
	keeper.RunTests(t, testCases)
}

func Test_Common_handlerAccountTransfer(t *testing.T) {
	testCases := map[string]keeper.SubTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// do nothing
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
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
					Admin:        keeper.BobKey,
					ValidUntil:   0,
					Type:         types.OpenDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
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
					Admin:        iovns.ZeroAddress,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:         types.OpenDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
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
					Admin:        iovns.ZeroAddress,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:         types.OpenDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
				k.CreateAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        keeper.BobKey,
					ValidUntil:   0,
					Targets:      nil,
					Certificates: nil,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
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
					Admin:        keeper.AliceKey,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
				k.CreateAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        keeper.BobKey,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Targets:      nil,
					Certificates: nil,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    keeper.BobKey,
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
					Admin:        iovns.ZeroAddress,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:         types.OpenDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
				k.CreateAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        keeper.AliceKey,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Targets:      nil,
					Certificates: nil,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    keeper.BobKey,
					NewOwner: nil,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success domain without superuser": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        iovns.ZeroAddress,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:         types.OpenDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					Owner:      keeper.AliceKey,
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    keeper.AliceKey,
					NewOwner: keeper.BobKey,
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
				if !account.Owner.Equals(keeper.BobKey) {
					t.Fatalf("handlerAccounTransfer() expected new owner: %s, got: %s", keeper.BobKey, account.Owner)
				}
			},
		},
	}
	keeper.RunTests(t, testCases)
}

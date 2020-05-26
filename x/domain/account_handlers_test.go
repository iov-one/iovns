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
	dt "github.com/iov-one/iovns/x/domain/testing"
	"github.com/iov-one/iovns/x/domain/types"
)

func Test_handlerMsgAddAccountCertificates(t *testing.T) {
	cases := map[string]dt.SubTest{
		"domain does not exist": {
			BeforeTest: nil,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "does not exist",
					Name:           "",
					Owner:          dt.BobKey,
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
					Admin:      dt.BobKey,
				})
				//
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "",
					Owner:          dt.BobKey,
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
					Admin:      dt.BobKey,
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
					Admin:      dt.BobKey,
				})
				// add mock account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: 0,
					Owner:      dt.AliceKey,
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
					Admin:      dt.AliceKey,
				})
				// add mock account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      dt.AliceKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          dt.BobKey,
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
					Admin:      dt.BobKey,
				})
				// add mock account
				k.CreateAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "test",
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:        dt.AliceKey,
					Certificates: []types.Certificate{[]byte("test")},
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          dt.AliceKey,
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
					Admin:      dt.AliceKey,
				})
				// add mock account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      dt.AliceKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          dt.AliceKey,
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
	dt.RunTests(t, cases)
}

func Test_handlerMsgDeleteAccountCertificate(t *testing.T) {
	cases := map[string]dt.SubTest{
		"account does not exist": {
			BeforeTest: nil,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "does not exist",
					DeleteCertificate: nil,
					Owner:             dt.BobKey,
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
					Owner:  dt.AliceKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: nil,
					Owner:             dt.BobKey,
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
					Owner:        dt.AliceKey,
					Certificates: nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: []byte("does not exist"),
					Owner:             dt.AliceKey,
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
					Owner:        dt.AliceKey,
					Certificates: []types.Certificate{[]byte("test")},
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: []byte("test"),
					Owner:             dt.AliceKey,
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

	dt.RunTests(t, cases)
}

func Test_handlerMsgDeleteAccount(t *testing.T) {
	cases := map[string]dt.SubTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {

			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "does not exist",
					Name:   "does not exist",
					Owner:  dt.AliceKey,
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
					Admin: dt.BobKey,
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
		"msg owner does not own domain or account": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:  "test",
					Admin: dt.AliceKey,
				})
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "test",
					Owner:  dt.AliceKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  dt.BobKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgDeleteAccount() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success domain owner": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:  "test",
					Admin: dt.AliceKey,
				})
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "test",
					Owner:  dt.BobKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  dt.AliceKey,
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
					Admin: dt.AliceKey,
				})
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "test",
					Owner:  dt.BobKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  dt.BobKey,
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
	dt.RunTests(t, cases)
}

func Test_handleMsgRegisterAccount(t *testing.T) {
	testCases := map[string]dt.SubTest{
		"fail invalid blockchain targets address": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in blockchain targets
				// get set config function
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: dt.RegexMatchNothing, // don't match anything
					ValidBlockchainID:      dt.RegexMatchAll,     // match all
					DomainRenew:            10,
				})
				// add a domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        dt.BobKey,
					ValidUntil:   0,
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  dt.AliceKey,
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
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in blockchain targets
				// get set config function
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      dt.RegexMatchNothing, // don't match anything
					ValidBlockchainAddress: dt.RegexMatchAll,     // match all
					DomainRenew:            10,
				})
				// add a domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        dt.BobKey,
					ValidUntil:   0,
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  dt.AliceKey,
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
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: dt.RegexMatchAll,     // match all
					ValidBlockchainID:      dt.RegexMatchAll,     // match all
					ValidName:              dt.RegexMatchNothing, // match nothing
					DomainRenew:            10,
				})
				// add a domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        dt.AliceKey,
					ValidUntil:   2,
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "this won't match",
					Owner:  dt.AliceKey,
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
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: dt.RegexMatchAll, // match all
					ValidBlockchainID:      dt.RegexMatchAll, // match all
					ValidName:              dt.RegexMatchAll, // match nothing
					DomainRenew:            10,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "this does not exist",
					Name:   "works",
					Owner:  dt.AliceKey,
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
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: dt.RegexMatchAll, // match all
					ValidBlockchainID:      dt.RegexMatchAll, // match all
					ValidName:              dt.RegexMatchAll, // match nothing
					DomainRenew:            10,
				})
				// add a domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        dt.BobKey,
					ValidUntil:   2,
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  dt.AliceKey, // invalid owner
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
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: dt.RegexMatchAll, // match all
					ValidBlockchainID:      dt.RegexMatchAll, // match all
					ValidName:              dt.RegexMatchAll, // match nothing
					DomainRenew:            10,
				})
				// add a domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        dt.BobKey,
					ValidUntil:   0, // domain is expired
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  dt.BobKey,
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
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: dt.RegexMatchAll, // match all
					ValidBlockchainID:      dt.RegexMatchAll, // match all
					ValidName:              dt.RegexMatchAll, // match nothing
					DomainRenew:            10,
				})
				// add a domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        dt.BobKey,
					ValidUntil:   time.Now().Add(100000 * time.Hour).Unix(),
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
				// add an account that we are gonna try to overwrite
				k.CreateAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "exists",
					Owner:        dt.AliceKey,
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
					Owner:  dt.BobKey,
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
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: dt.RegexMatchAll, // match all
					ValidBlockchainID:      dt.RegexMatchAll, // match all
					ValidName:              dt.RegexMatchAll, // match nothing
					DomainRenew:            10,
				})
				// add a domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        dt.BobKey,
					ValidUntil:   time.Now().Add(100000 * time.Hour).Unix(),
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  dt.BobKey,
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
	dt.RunTests(t, testCases)
}

func Test_handlerMsgRenewAccount(t *testing.T) {
	cases := map[string]dt.SubTest{
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
					Admin:        dt.BobKey,
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
					Admin:        dt.BobKey,
				})
				// set mock account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Unix(1000, 0)),
					Owner:      dt.AliceKey,
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

	dt.RunTests(t, cases)
}

func Test_handlerMsgReplaceAccountTargets(t *testing.T) {
	cases := map[string]dt.SubTest{
		"invalid blockchain target": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set config to match nothing
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      dt.RegexMatchNothing,
					ValidBlockchainAddress: dt.RegexMatchNothing,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, &types.MsgReplaceAccountTargets{
					Domain: "",
					Name:   "",
					NewTargets: []types.BlockchainAddress{
						{
							ID:      "invalid",
							Address: "invalid",
						},
					},
					Owner: nil,
				})
				if !errors.Is(err, types.ErrInvalidBlockchainTarget) {
					t.Fatalf("handlerMsgReplaceAccountTargets() expected error: %s, got: %s", types.ErrInvalidBlockchainTarget, err)
				}
			},
			AfterTest: nil,
		},
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      dt.RegexMatchAll,
					ValidBlockchainAddress: dt.RegexMatchAll,
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
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      dt.RegexMatchAll,
					ValidBlockchainAddress: dt.RegexMatchAll,
				})
				// create domain
				k.CreateDomain(ctx, types.Domain{
					Name:  "test",
					Admin: dt.BobKey,
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
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      dt.RegexMatchAll,
					ValidBlockchainAddress: dt.RegexMatchAll,
				})
				// create domain
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      dt.BobKey,
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
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      dt.RegexMatchAll,
					ValidBlockchainAddress: dt.RegexMatchAll,
				})
				// create domain
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      dt.BobKey,
				})
				// create account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: 0,
					Owner:      dt.AliceKey,
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
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      dt.RegexMatchAll,
					ValidBlockchainAddress: dt.RegexMatchAll,
				})
				// create domain
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      dt.BobKey,
				})
				// create account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      dt.AliceKey,
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
					Owner: dt.BobKey,
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
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      dt.RegexMatchAll,
					ValidBlockchainAddress: dt.RegexMatchAll,
				})
				// create domain
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      dt.BobKey,
				})
				// create account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      dt.AliceKey,
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
					Owner: dt.AliceKey,
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
	dt.RunTests(t, cases)
}

func Test_handlerMsgSetAccountMetadata(t *testing.T) {
	cases := map[string]dt.SubTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgSetAccountMetadata(ctx, k, &types.MsgSetAccountMetadata{
					Domain: "does not exist",
					Name:   "",
					Owner:  dt.AliceKey,
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
					Admin: dt.BobKey,
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
					Admin:      dt.BobKey,
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
					Admin:      dt.BobKey,
				})
				// create account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: 0,
					Owner:      dt.AliceKey,
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
					Admin:      dt.BobKey,
				})
				// create account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      dt.AliceKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgSetAccountMetadata(ctx, k, &types.MsgSetAccountMetadata{
					Domain: "test",
					Name:   "test",
					Owner:  dt.BobKey,
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
					Admin:      dt.BobKey,
				})
				// create account
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      dt.AliceKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgSetAccountMetadata(ctx, k, &types.MsgSetAccountMetadata{
					Domain:         "test",
					Name:           "test",
					NewMetadataURI: "https://test.com",
					Owner:          dt.AliceKey,
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
	dt.RunTests(t, cases)
}

func Test_handlerAccountTransfer(t *testing.T) {
	testCases := map[string]dt.SubTest{
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
					Admin:        dt.BobKey,
					ValidUntil:   0,
					HasSuperuser: false,
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
					HasSuperuser: false,
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
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
				k.CreateAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        dt.BobKey,
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
					Admin:        dt.AliceKey,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
				k.CreateAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        dt.BobKey,
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
					Owner:    dt.BobKey,
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
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
				k.CreateAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        dt.AliceKey,
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
					Owner:    dt.BobKey,
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
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					Owner:      dt.AliceKey,
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    dt.AliceKey,
					NewOwner: dt.BobKey,
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
				if !account.Owner.Equals(dt.BobKey) {
					t.Fatalf("handlerAccounTransfer() expected new owner: %s, got: %s", dt.BobKey, account.Owner)
				}
			},
		},
		"success domain with superuser": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// domain owned by alice
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        dt.AliceKey,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
				// account owned by bob
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					Owner:      dt.BobKey,
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// alice is domain owner and should transfer account owned by bob to alice
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    dt.AliceKey,
					NewOwner: dt.CharlieKey,
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
				if !account.Owner.Equals(dt.CharlieKey) {
					t.Fatalf("handlerAccounTransfer() expected new owner: %s, got: %s", dt.CharlieKey, account.Owner)
				}
			},
		},
	}
	dt.RunTests(t, testCases)
}

func Test_validateBlockchainTargets(t *testing.T) {
	type args struct {
		targets []types.BlockchainAddress
		conf    configuration.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "duplicate blockchain target",
			args: args{
				targets: []types.BlockchainAddress{
					{
						ID:      "duplicate",
						Address: "does not matter",
					},
					{
						ID:      "duplicate",
						Address: "does not matter",
					},
				},
				conf: configuration.Config{
					ValidBlockchainID:      dt.RegexMatchAll,
					ValidBlockchainAddress: dt.RegexMatchAll,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateBlockchainTargets(tt.args.targets, tt.args.conf); (err != nil) != tt.wantErr {
				t.Errorf("validateBlockchainTargets() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

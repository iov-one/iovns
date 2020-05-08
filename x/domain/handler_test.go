package domain

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/keeper"
	types2 "github.com/iov-one/iovns/x/domain/types"
)

func Test_handlerMsgAddAccountCertificates(t *testing.T) {
	cases := map[string]subTest{
		"domain does not exist": {
			BeforeTest: nil,
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types2.MsgAddAccountCertificates{
					Domain:         "does not exist",
					Name:           "",
					Owner:          nil,
					NewCertificate: nil,
				})
				if !errors.Is(err, types2.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types2.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// add expired domain
				k.CreateDomain(ctx, types2.Domain{
					Name:       "test",
					ValidUntil: 0,
				})
				//
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types2.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "",
					Owner:          nil,
					NewCertificate: nil,
				})
				if !errors.Is(err, types2.ErrDomainExpired) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types2.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				})
				//
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types2.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "does not exist",
					Owner:          nil,
					NewCertificate: nil,
				})
				if !errors.Is(err, types2.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types2.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"account expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				})
				// add mock account
				k.CreateAccount(ctx, types2.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: 0,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types2.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          nil,
					NewCertificate: nil,
				})
				if !errors.Is(err, types2.ErrAccountExpired) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types2.ErrAccountExpired, err)
				}
			},
			AfterTest: nil,
		},
		"msg owner is not account owner": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				})
				// add mock account
				k.CreateAccount(ctx, types2.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      aliceKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types2.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          bobKey.GetAddress(),
					NewCertificate: nil,
				})
				if !errors.Is(err, types2.ErrUnauthorized) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types2.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"certificate exists": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				})
				// add mock account
				k.CreateAccount(ctx, types2.Account{
					Domain:       "test",
					Name:         "test",
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:        aliceKey.GetAddress(),
					Certificates: [][]byte{[]byte("test")},
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types2.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          aliceKey.GetAddress(),
					NewCertificate: []byte("test"),
				})
				if !errors.Is(err, types2.ErrCertificateExists) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types2.ErrCertificateExists, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				})
				// add mock account
				k.CreateAccount(ctx, types2.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      aliceKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types2.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          aliceKey.GetAddress(),
					NewCertificate: []byte("test"),
				})
				if err != nil {
					t.Fatalf("handlerMsgAddAccountCertificates() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
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

func Test_handlerMsgDeleteAccountCertificate(t *testing.T) {
	cases := map[string]subTest{
		"account does not exist": {
			BeforeTest: nil,
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types2.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "does not exist",
					DeleteCertificate: nil,
					Owner:             nil,
				})
				if !errors.Is(err, types2.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgDeleteAccountCertificate() expected error: %s, got: %s", types2.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"msg signer is not account owner": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateAccount(ctx, types2.Account{
					Domain: "test",
					Name:   "test",
					Owner:  aliceKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types2.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: nil,
					Owner:             bobKey.GetAddress(),
				})
				if !errors.Is(err, types2.ErrUnauthorized) {
					t.Fatalf("handlerMsgDeleteAccountCertificate() expected error: %s, got: %s", types2.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"certificate does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateAccount(ctx, types2.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        aliceKey.GetAddress(),
					Certificates: nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types2.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: []byte("does not exist"),
					Owner:             aliceKey.GetAddress(),
				})
				if !errors.Is(err, types2.ErrCertificateDoesNotExist) {
					t.Fatalf("handlerMsgDeleteAccountCertificate() expected error: %s, got: %s", types2.ErrCertificateDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateAccount(ctx, types2.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        aliceKey.GetAddress(),
					Certificates: [][]byte{[]byte("test")},
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types2.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: []byte("test"),
					Owner:             aliceKey.GetAddress(),
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccountCertificates() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
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

	runTests(t, cases)
}

func Test_handlerMsgDeleteAccount(t *testing.T) {
	cases := map[string]subTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {

			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types2.MsgDeleteAccount{
					Domain: "does not exist",
					Name:   "does not exist",
					Owner:  nil,
				})
				if !errors.Is(err, types2.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgDeleteAccount() expected error: %s, got: %s", types2.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name: "test",
				})

			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types2.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  nil,
				})
				if !errors.Is(err, types2.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgDeleteAccount() expected error: %s, got: %s", types2.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"msg owner does not own domain or account": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name:  "test",
					Admin: aliceKey.GetAddress(),
				})
				k.CreateAccount(ctx, types2.Account{
					Domain: "test",
					Name:   "test",
					Owner:  aliceKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types2.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  bobKey.GetAddress(),
				})
				if !errors.Is(err, types2.ErrUnauthorized) {
					t.Fatalf("handlerMsgDeleteAccount() expected error: %s, got: %s", types2.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success domain owner": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name:  "test",
					Admin: aliceKey.GetAddress(),
				})
				k.CreateAccount(ctx, types2.Account{
					Domain: "test",
					Name:   "test",
					Owner:  bobKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types2.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  aliceKey.GetAddress(),
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, exists := k.GetAccount(ctx, "test", "test")
				if exists {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
		"success account owner": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name:  "test",
					Admin: aliceKey.GetAddress(),
				})
				k.CreateAccount(ctx, types2.Account{
					Domain: "test",
					Name:   "test",
					Owner:  bobKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types2.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  bobKey.GetAddress(),
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, exists := k.GetAccount(ctx, "test", "test")
				if exists {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
	}

	// run tests
	runTests(t, cases)
}

func Test_handleMsgRegisterAccount(t *testing.T) {
	testCases := map[string]subTest{
		"fail invalid blockchain targets address": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in blockchain targets
				// get set config function
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: regexMatchNothing, // don't match anything
					ValidBlockchainID:      regexMatchAll,     // match all
					DomainRenew:            10,
				})
				// add a domain
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					Admin:        nil,
					ValidUntil:   0,
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types2.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  aliceKey.GetAddress(),
					Targets: []iovns.BlockchainAddress{
						{
							ID:      "works",
							Address: "won't work",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types2.ErrInvalidBlockchainTarget) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types2.ErrInvalidBlockchainTarget, err)
				}
			},
			AfterTest: nil,
		},
		// TODO cleanup comments
		"fail invalid blockchain targets id": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in blockchain targets
				// get set config function
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      regexMatchNothing, // don't match anything
					ValidBlockchainAddress: regexMatchAll,     // match all
					DomainRenew:            10,
				})
				// add a domain
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					Admin:        nil,
					ValidUntil:   0,
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types2.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  aliceKey.GetAddress(),
					Targets: []iovns.BlockchainAddress{
						{
							ID:      "invalid blockchain id",
							Address: "valid blockchain address",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types2.ErrInvalidBlockchainTarget) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types2.ErrInvalidBlockchainTarget, err)
				}
			},
			AfterTest: nil,
		},
		"fail invalid account name": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// set config
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: regexMatchAll,     // match all
					ValidBlockchainID:      regexMatchAll,     // match all
					ValidName:              regexMatchNothing, // match nothing
					DomainRenew:            10,
				})
				// add a domain
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					Admin:        nil,
					ValidUntil:   0,
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types2.MsgRegisterAccount{
					Domain: "test",
					Name:   "this won't match",
					Owner:  aliceKey.GetAddress(),
					Targets: []iovns.BlockchainAddress{
						{
							ID:      "works",
							Address: "works",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types2.ErrInvalidAccountName) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types2.ErrInvalidAccountName, err)
				}
			},
			AfterTest: nil,
		},
		"fail domain name does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in blockchain targets
				// get set config function
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: regexMatchAll, // match all
					ValidBlockchainID:      regexMatchAll, // match all
					ValidName:              regexMatchAll, // match nothing
					DomainRenew:            10,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types2.MsgRegisterAccount{
					Domain: "this does not exist",
					Name:   "works",
					Owner:  aliceKey.GetAddress(),
					Targets: []iovns.BlockchainAddress{
						{
							ID:      "works",
							Address: "works",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types2.ErrDomainDoesNotExist) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types2.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"fail only owner of domain with superuser can register accounts": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in blockchain targets
				// get set config function
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: regexMatchAll, // match all
					ValidBlockchainID:      regexMatchAll, // match all
					ValidName:              regexMatchAll, // match nothing
					DomainRenew:            10,
				})
				// add a domain
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					Admin:        bobKey.GetAddress(),
					ValidUntil:   0,
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types2.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  aliceKey.GetAddress(), // invalid owner
					Targets: []iovns.BlockchainAddress{
						{
							ID:      "works",
							Address: "works",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types2.ErrUnauthorized) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types2.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"fail domain has expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in blockchain targets
				// get set config function
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: regexMatchAll, // match all
					ValidBlockchainID:      regexMatchAll, // match all
					ValidName:              regexMatchAll, // match nothing
					DomainRenew:            10,
				})
				// add a domain
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					Admin:        bobKey.GetAddress(),
					ValidUntil:   0, // domain is expired
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types2.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  bobKey.GetAddress(),
					Targets: []iovns.BlockchainAddress{
						{
							ID:      "works",
							Address: "works",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types2.ErrDomainExpired) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types2.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"fail account exists": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in blockchain targets
				// get set config function
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: regexMatchAll, // match all
					ValidBlockchainID:      regexMatchAll, // match all
					ValidName:              regexMatchAll, // match nothing
					DomainRenew:            10,
				})
				// add a domain
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					Admin:        bobKey.GetAddress(),
					ValidUntil:   time.Now().Add(100000 * time.Hour).Unix(),
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
				// add an account that we are gonna try to overwrite
				k.CreateAccount(ctx, types2.Account{
					Domain:       "test",
					Name:         "exists",
					Owner:        nil,
					ValidUntil:   0,
					Targets:      nil,
					Certificates: nil,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types2.MsgRegisterAccount{
					Domain: "test",
					Name:   "exists",
					Owner:  bobKey.GetAddress(),
					Targets: []iovns.BlockchainAddress{
						{
							ID:      "works",
							Address: "works",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types2.ErrAccountExists) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types2.ErrAccountExists, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in blockchain targets
				// get set config function
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidBlockchainAddress: regexMatchAll, // match all
					ValidBlockchainID:      regexMatchAll, // match all
					ValidName:              regexMatchAll, // match nothing
					DomainRenew:            10,
				})
				// add a domain
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					Admin:        bobKey.GetAddress(),
					ValidUntil:   time.Now().Add(100000 * time.Hour).Unix(),
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types2.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  bobKey.GetAddress(),
					Targets: []iovns.BlockchainAddress{
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
	runTests(t, testCases)
}

func Test_handlerMsgRenewAccount(t *testing.T) {
	cases := map[string]subTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {

			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgRenewAccount(ctx, k, &types2.MsgRenewAccount{
					Domain: "does not exist",
					Name:   "",
				})
				if !errors.Is(err, types2.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgRenewAccount() expected error: %s, got: %s", types2.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// set mock domain
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					AccountRenew: 100,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgRenewAccount(ctx, k, &types2.MsgRenewAccount{
					Domain: "test",
					Name:   "does not exist",
				})
				if !errors.Is(err, types2.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgRenewAccount() expected error: %s, got: %s", types2.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// set mock domain
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					AccountRenew: 100,
				})
				// set mock account
				k.CreateAccount(ctx, types2.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: 1000,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgRenewAccount(ctx, k, &types2.MsgRenewAccount{
					Domain: "test",
					Name:   "test",
				})
				if err != nil {
					t.Fatalf("handlerMsgRenewAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				account, _ := k.GetAccount(ctx, "test", "test")
				if account.ValidUntil != 1100 {
					t.Fatalf("handlerMsgRenewAccount() expected 1100, got: %d", account.ValidUntil)
				}
			},
		},
	}

	runTests(t, cases)
}

func Test_handlerMsgReplaceAccountTargets(t *testing.T) {
	cases := map[string]subTest{
		"invalid blockchain target": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// set config to match nothing
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      regexMatchNothing,
					ValidBlockchainAddress: regexMatchNothing,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, &types2.MsgReplaceAccountTargets{
					Domain: "",
					Name:   "",
					NewTargets: []iovns.BlockchainAddress{
						{
							ID:      "invalid",
							Address: "invalid",
						},
					},
					Owner: nil,
				})
				if !errors.Is(err, types2.ErrInvalidBlockchainTarget) {
					t.Fatalf("handlerMsgReplaceAccountTargets() expected error: %s, got: %s", types2.ErrInvalidBlockchainTarget, err)
				}
			},
			AfterTest: nil,
		},
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      regexMatchAll,
					ValidBlockchainAddress: regexMatchAll,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, &types2.MsgReplaceAccountTargets{
					Domain: "does not exist",
					Name:   "",
					NewTargets: []iovns.BlockchainAddress{
						{
							ID:      "valid",
							Address: "valid",
						},
					},
					Owner: nil,
				})
				if !errors.Is(err, types2.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgReplaceAccountTargets() expected error: %s, got: %s", types2.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      regexMatchAll,
					ValidBlockchainAddress: regexMatchAll,
				})
				// create domain
				k.CreateDomain(ctx, types2.Domain{
					Name: "test",
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, &types2.MsgReplaceAccountTargets{
					Domain: "test",
					NewTargets: []iovns.BlockchainAddress{
						{
							ID:      "valid",
							Address: "valid",
						},
					},
					Owner: nil,
				})
				if !errors.Is(err, types2.ErrDomainExpired) {
					t.Fatalf("handlerMsgReplaceAccountTargets() expected error: %s, got: %s", types2.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      regexMatchAll,
					ValidBlockchainAddress: regexMatchAll,
				})
				// create domain
				k.CreateDomain(ctx, types2.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, &types2.MsgReplaceAccountTargets{
					Domain: "test",
					Name:   "does not exist",
					NewTargets: []iovns.BlockchainAddress{
						{
							ID:      "valid",
							Address: "valid",
						},
					},
					Owner: nil,
				})
				if !errors.Is(err, types2.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgReplaceAccountTargets() expected error: %s, got: %s", types2.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"account expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      regexMatchAll,
					ValidBlockchainAddress: regexMatchAll,
				})
				// create domain
				k.CreateDomain(ctx, types2.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
				})
				// create account
				k.CreateAccount(ctx, types2.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: 0,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, &types2.MsgReplaceAccountTargets{
					Domain: "test",
					Name:   "test",
					NewTargets: []iovns.BlockchainAddress{
						{
							ID:      "valid",
							Address: "valid",
						},
					},
					Owner: nil,
				})
				if !errors.Is(err, types2.ErrAccountExpired) {
					t.Fatalf("handlerMsgReplaceAccountTargets() expected error: %s, got: %s", types2.ErrAccountExpired, err)
				}
			},
			AfterTest: nil,
		},
		"signer is not owner of account": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      regexMatchAll,
					ValidBlockchainAddress: regexMatchAll,
				})
				// create domain
				k.CreateDomain(ctx, types2.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
				})
				// create account
				k.CreateAccount(ctx, types2.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      aliceKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, &types2.MsgReplaceAccountTargets{
					Domain: "test",
					Name:   "test",
					NewTargets: []iovns.BlockchainAddress{
						{
							ID:      "valid",
							Address: "valid",
						},
					},
					Owner: bobKey.GetAddress(),
				})
				if !errors.Is(err, types2.ErrUnauthorized) {
					t.Fatalf("handlerMsgReplaceAccountTargets() expected error: %s, got: %s", types2.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      regexMatchAll,
					ValidBlockchainAddress: regexMatchAll,
				})
				// create domain
				k.CreateDomain(ctx, types2.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
				})
				// create account
				k.CreateAccount(ctx, types2.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      aliceKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, &types2.MsgReplaceAccountTargets{
					Domain: "test",
					Name:   "test",
					NewTargets: []iovns.BlockchainAddress{
						{
							ID:      "valid",
							Address: "valid",
						},
					},
					Owner: aliceKey.GetAddress(),
				})
				if err != nil {
					t.Fatalf("handlerMsgReplaceAccountTargets() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				expected := []iovns.BlockchainAddress{{
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
	runTests(t, cases)
}

func Test_handlerMsgSetAccountMetadata(t *testing.T) {
	cases := map[string]subTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgSetAccountMetadata(ctx, k, &types2.MsgSetAccountMetadata{
					Domain: "does not exist",
					Name:   "",
				})
				if !errors.Is(err, types2.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgSetAccountMetadata() expected error: %s, got: %s", types2.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// create domain
				k.CreateDomain(ctx, types2.Domain{
					Name: "test",
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgSetAccountMetadata(ctx, k, &types2.MsgSetAccountMetadata{
					Domain:         "test",
					Name:           "",
					NewMetadataURI: "",
					Owner:          nil,
				})
				if !errors.Is(err, types2.ErrDomainExpired) {
					t.Fatalf("handlerMsgSetAccountMetadata() expected error: %s, got: %s", types2.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// create domain
				k.CreateDomain(ctx, types2.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgSetAccountMetadata(ctx, k, &types2.MsgSetAccountMetadata{
					Domain: "test",
					Name:   "does not exist",
					Owner:  nil,
				})
				if !errors.Is(err, types2.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgSetAccountMetadata() expected error: %s, got: %s", types2.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"account expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// create domain
				k.CreateDomain(ctx, types2.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
				})
				// create account
				k.CreateAccount(ctx, types2.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: 0,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgSetAccountMetadata(ctx, k, &types2.MsgSetAccountMetadata{
					Domain: "test",
					Name:   "test",
					Owner:  nil,
				})
				if !errors.Is(err, types2.ErrAccountExpired) {
					t.Fatalf("handlerMsgSetAccountMetadata() expected error: %s, got: %s", types2.ErrAccountExpired, err)
				}
			},
			AfterTest: nil,
		},
		"signer is not owner of account": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// create domain
				k.CreateDomain(ctx, types2.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
				})
				// create account
				k.CreateAccount(ctx, types2.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      aliceKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgSetAccountMetadata(ctx, k, &types2.MsgSetAccountMetadata{
					Domain: "test",
					Name:   "test",
					Owner:  bobKey.GetAddress(),
				})
				if !errors.Is(err, types2.ErrUnauthorized) {
					t.Fatalf("handlerMsgSetAccountMetadata() expected error: %s, got: %s", types2.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// create domain
				k.CreateDomain(ctx, types2.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
				})
				// create account
				k.CreateAccount(ctx, types2.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      aliceKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgSetAccountMetadata(ctx, k, &types2.MsgSetAccountMetadata{
					Domain:         "test",
					Name:           "test",
					NewMetadataURI: "https://test.com",
					Owner:          aliceKey.GetAddress(),
				})
				if err != nil {
					t.Fatalf("handlerMsgSetAccountMetadata() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
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

func Test_handlerAccountTransfer(t *testing.T) {
	testCases := map[string]subTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// do nothing
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types2.MsgTransferAccount{
					Domain:   "does not exist",
					Name:     "does not exist",
					Owner:    nil,
					NewOwner: nil,
				})
				if !errors.Is(err, types2.ErrDomainDoesNotExist) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types2.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain has expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name:         "expired domain",
					Admin:        nil,
					ValidUntil:   0,
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types2.MsgTransferAccount{
					Domain:   "expired domain",
					Name:     "",
					Owner:    nil,
					NewOwner: nil,
				})
				if !errors.Is(err, types2.ErrDomainExpired) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types2.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					Admin:        nil,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types2.MsgTransferAccount{
					Domain:   "test",
					Name:     "this account does not exist",
					Owner:    nil,
					NewOwner: nil,
				})
				if !errors.Is(err, types2.ErrAccountDoesNotExist) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types2.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"account expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					Admin:        nil,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
				k.CreateAccount(ctx, types2.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        nil,
					ValidUntil:   0,
					Targets:      nil,
					Certificates: nil,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types2.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    nil,
					NewOwner: nil,
				})
				if !errors.Is(err, types2.ErrAccountExpired) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types2.ErrAccountExpired, err)
				}
			},
			AfterTest: nil,
		},
		"if domain has super user only domain admin can transfer accounts": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					Admin:        aliceKey.GetAddress(),
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
				k.CreateAccount(ctx, types2.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        nil,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Targets:      nil,
					Certificates: nil,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types2.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    bobKey.GetAddress(),
					NewOwner: nil,
				})
				if !errors.Is(err, types2.ErrUnauthorized) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types2.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"if domain has no super user then only account owner can transfer accounts": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					Admin:        nil,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
				k.CreateAccount(ctx, types2.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        aliceKey.GetAddress(),
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Targets:      nil,
					Certificates: nil,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types2.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    bobKey.GetAddress(),
					NewOwner: nil,
				})
				if !errors.Is(err, types2.ErrUnauthorized) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types2.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					Admin:        nil,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
				k.CreateAccount(ctx, types2.Account{
					Domain:       "test",
					Name:         "test",
					Owner:        aliceKey.GetAddress(),
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Targets:      nil,
					Certificates: nil,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types2.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    aliceKey.GetAddress(),
					NewOwner: nil,
				})
				if err != nil {
					t.Fatalf("handlerMsgTransferAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
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

func Test_handleMsgDomainDelete(t *testing.T) {
	cases := map[string]subTest{
		"fail domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// don't do anything
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types2.MsgDeleteDomain{
					Domain: "this does not exist",
					Owner:  bobKey.GetAddress(),
				})
				if !errors.Is(err, types2.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgDeleteDomain() expected error: %s, got: %s", types2.ErrDomainDoesNotExist, err)
				}
			},
		},
		"fail domain has no superuser": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// set domain with no superuser
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					Admin:        nil,
					ValidUntil:   0,
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types2.MsgDeleteDomain{
					Domain: "test",
					Owner:  bobKey.GetAddress(),
				})
				if !errors.Is(err, types2.ErrUnauthorized) {
					t.Fatalf("handlerMsgDeleteDomain() expected error: %s, got: %s", types2.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"fail domain admin does not match msg owner": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 1000000000000000,
				})
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					Admin:        bobKey.GetAddress(),
					ValidUntil:   0,
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types2.MsgDeleteDomain{
					Domain: "test",
					Owner:  aliceKey.GetAddress(),
				})
				if !errors.Is(err, types2.ErrUnauthorized) {
					t.Fatalf("handlerMsgDeleteDomain() expected error: %s, got: %s", types2.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 1000000000000000, // unexpired domain
				})
				// set domain
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					Admin:        aliceKey.GetAddress(),
					ValidUntil:   0,
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
				// add two accounts
				k.CreateAccount(ctx, types2.Account{
					Domain: "test",
					Name:   "1",
				})
				// add two accounts
				k.CreateAccount(ctx, types2.Account{
					Domain: "test",
					Name:   "2",
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types2.MsgDeleteDomain{
					Domain: "test",
					Owner:  aliceKey.GetAddress(),
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteDomain() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, exists := k.GetDomain(ctx, "test")
				if exists {
					t.Fatalf("handlerMsgDeleteDomain() domain should not exist")
				}
				_, exists = k.GetAccount(ctx, "test", "1")
				if exists {
					t.Fatalf("handlerMsgDeleteDomain() account 1 should not exist")
				}
				_, exists = k.GetAccount(ctx, "test", "2")
				if exists {
					t.Fatalf("handlerMsgDeleteDomain() account 2 should not exist")
				}
			},
		},
		"success claim expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 1,
				})
				// set domain
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					Admin:        aliceKey.GetAddress(),
					ValidUntil:   0,
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
				// add two accounts
				k.CreateAccount(ctx, types2.Account{
					Domain: "test",
					Name:   "1",
				})
				// add two accounts
				k.CreateAccount(ctx, types2.Account{
					Domain: "test",
					Name:   "2",
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types2.MsgDeleteDomain{
					Domain: "test",
					Owner:  aliceKey.GetAddress(),
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteDomain() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, exists := k.GetDomain(ctx, "test")
				if exists {
					t.Fatalf("handlerMsgDeleteDomain() domain should not exist")
				}
				_, exists = k.GetAccount(ctx, "test", "1")
				if exists {
					t.Fatalf("handlerMsgDeleteDomain() account 1 should not exist")
				}
				_, exists = k.GetAccount(ctx, "test", "2")
				if exists {
					t.Fatalf("handlerMsgDeleteDomain() account 2 should not exist")
				}
			},
		},
	}
	runTests(t, cases)
}

func Test_handlerMsgFlushDomain(t *testing.T) {
	cases := map[string]subTest{
		"domain does not exist": {
			BeforeTest: nil,
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgFlushDomain(ctx, k, &types2.MsgFlushDomain{
					Domain: "does not exist",
					Owner:  nil,
				})
				if !errors.Is(err, types2.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgFlushDomain() expected error: %s, got: %s", types2.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain has superuser": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					Admin:        nil,
					ValidUntil:   0,
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgFlushDomain(ctx, k, &types2.MsgFlushDomain{
					Domain: "test",
					Owner:  nil,
				})
				if !errors.Is(err, types2.ErrUnauthorized) {
					t.Fatalf("handlerMsgFlushDomain() expected error: %s, got: %s", types2.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"msg owner is not domain admin": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					Admin:        aliceKey.GetAddress(),
					ValidUntil:   0,
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgFlushDomain(ctx, k, &types2.MsgFlushDomain{
					Domain: "test",
					Owner:  bobKey.GetAddress(),
				})
				if !errors.Is(err, types2.ErrUnauthorized) {
					t.Fatalf("handlerMsgFlushDomain() expected error: %s, got: %s", types2.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// set domain
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					Admin:        aliceKey.GetAddress(),
					ValidUntil:   0,
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
				// add empty account 1
				k.CreateAccount(ctx, types2.Account{
					Domain: "test",
					Name:   "",
				})
				// add account 2
				k.CreateAccount(ctx, types2.Account{
					Domain: "test",
					Name:   "1",
				})
				// add account 2
				k.CreateAccount(ctx, types2.Account{
					Domain: "test",
					Name:   "2",
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgFlushDomain(ctx, k, &types2.MsgFlushDomain{
					Domain: "test",
					Owner:  aliceKey.GetAddress(),
				})
				if err != nil {
					t.Fatalf("handlerMsgFlushDomain() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
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

func TestHandleMsgRegisterDomain(t *testing.T) {
	testCases := map[string]subTest{
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				configSetter := getConfigSetter(k.ConfigurationKeeper)
				// set config
				configSetter.SetConfig(ctx, configuration.Config{
					Owners:      []types.AccAddress{aliceKey.GetAddress()},
					ValidDomain: "^(.*?)?",
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// register domain with superuser
				_, err := handleMsgRegisterDomain(ctx, k, &types2.MsgRegisterDomain{
					Name:         "domain",
					HasSuperuser: true,
					AccountRenew: 10,
				})
				if err != nil {
					t.Fatalf("handleMsgRegisterDomain() with superuser, got error: %s", err)
				}
				// register domain without super user
				_, err = handleMsgRegisterDomain(ctx, k, &types2.MsgRegisterDomain{
					Name:         "domain-without-superuser",
					Admin:        aliceKey.GetAddress(),
					HasSuperuser: false,
					Broker:       nil,
					AccountRenew: 20,
				})
				if err != nil {
					t.Fatalf("handleMsgRegisterDomain() without superuser, got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// TODO do reflect.DeepEqual checks on expected results vs results returned
				_, ok := k.GetDomain(ctx, "domain")
				if !ok {
					t.Fatalf("handleMsgRegisterDomain() could not find 'domain'")
				}
				_, ok = k.GetDomain(ctx, "domain-without-superuser")
				if !ok {
					t.Fatalf("handleMsgRegisterDomain() could not find 'domain-without-superuser'")
				}
			},
		},
		"fail domain name exists": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name:         "exists",
					Admin:        nil,
					ValidUntil:   0,
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterDomain(ctx, k, &types2.MsgRegisterDomain{
					Name:         "exists",
					Admin:        nil,
					HasSuperuser: false,
					Broker:       nil,
					AccountRenew: 0,
				})
				if !errors.Is(err, types2.ErrDomainAlreadyExists) {
					t.Fatalf("handleMsgRegisterDomain() expected: %s got: %s", types2.ErrDomainAlreadyExists, err)
				}
			},
			AfterTest: nil,
		},
		"fail domain does not match valid domain regexp": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// get set config function
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidDomain: "$^",
					DomainRenew: 0,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterDomain(ctx, k, &types2.MsgRegisterDomain{
					Name:         "invalid-name",
					Admin:        nil,
					HasSuperuser: false,
					Broker:       nil,
					AccountRenew: 0,
				})
				if !errors.Is(err, types2.ErrInvalidDomainName) {
					t.Fatalf("handleMsgRegisterDomain() expected error: %s, got: %s", types2.ErrInvalidDomainName, err)
				}
			},
			AfterTest: nil,
		},
		"fail domain with no super user must be registered by configuration owner": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// add config with owner
				config := configuration.Config{
					Owners:                 []types.AccAddress{aliceKey.GetAddress()},
					ValidDomain:            "^(.*?)?",
					ValidName:              "",
					ValidBlockchainID:      "",
					ValidBlockchainAddress: "",
					DomainRenew:            0,
				}
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, config)
			},
			Test: func(t *testing.T, k Keeper, ctx types.Context, mock *keeper.Mocks) {
				// try to register domain with no super user
				_, err := handleMsgRegisterDomain(ctx, k, &types2.MsgRegisterDomain{
					Name:         "some-domain",
					Admin:        bobKey.GetAddress(),
					HasSuperuser: false,
					Broker:       nil,
					AccountRenew: 10,
				})
				if !errors.Is(err, types2.ErrUnauthorized) {
					t.Fatalf("handleMsgRegisterDomain() expecter error: %s, got: %s", types2.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
	}
	// run all test cases
	runTests(t, testCases)
}

func Test_handlerDomainRenew(t *testing.T) {
	cases := map[string]subTest{
		"domain not found": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {

			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgRenewDomain(ctx, k, &types2.MsgRenewDomain{Domain: "does not exist"})
				if !errors.Is(err, types2.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgRenewDomain() expected error: %s, got: %s", types2.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// add config
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainRenew: 1,
				})
				// add domain
				k.CreateDomain(ctx, types2.Domain{
					Name:       "test",
					ValidUntil: 1000,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgRenewDomain(ctx, k, &types2.MsgRenewDomain{Domain: "test"})
				if err != nil {
					t.Fatalf("handlerMsgRenewDomain() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// get domain
				domain, _ := k.GetDomain(ctx, "test")
				if domain.ValidUntil != 1001 {
					t.Fatalf("handlerMsgRenewDomain() expected 1001, got: %d", domain.ValidUntil)
				}
			},
		},
	}
	// run tests
	runTests(t, cases)
}

func Test_handlerMsgTransferDomain(t *testing.T) {
	cases := map[string]subTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {

			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferDomain(ctx, k, &types2.MsgTransferDomain{
					Domain:   "does not exist",
					Owner:    nil,
					NewAdmin: nil,
				})
				if !errors.Is(err, types2.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgTransferDomain() expected error: %s, got error: %s", types2.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain has no superuser": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					HasSuperuser: false,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferDomain(ctx, k, &types2.MsgTransferDomain{
					Domain:   "test",
					Owner:    nil,
					NewAdmin: nil,
				})
				if !errors.Is(err, types2.ErrUnauthorized) {
					t.Fatalf("handlerMsgTransferDomain() expected error: %s, got error: %s", types2.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"domain has expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					HasSuperuser: true,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferDomain(ctx, k, &types2.MsgTransferDomain{
					Domain:   "test",
					Owner:    nil,
					NewAdmin: nil,
				})
				if !errors.Is(err, types2.ErrDomainExpired) {
					t.Fatalf("handlerMsgTransferDomain() expected error: %s, got error: %s", types2.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"msg signer is not domain admin": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					HasSuperuser: true,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:        aliceKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferDomain(ctx, k, &types2.MsgTransferDomain{
					Domain:   "test",
					Owner:    bobKey.GetAddress(),
					NewAdmin: nil,
				})
				if !errors.Is(err, types2.ErrUnauthorized) {
					t.Fatalf("handlerMsgTransferDomain() expected error: %s, got error: %s", types2.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				// create domain
				k.CreateDomain(ctx, types2.Domain{
					Name:         "test",
					HasSuperuser: true,
					ValidUntil:   iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:        aliceKey.GetAddress(),
				})
				// add empty account
				k.CreateAccount(ctx, types2.Account{
					Domain: "test",
					Name:   "",
				})
				// add account 1
				k.CreateAccount(ctx, types2.Account{
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
				k.CreateAccount(ctx, types2.Account{
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
			Test: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferDomain(ctx, k, &types2.MsgTransferDomain{
					Domain:   "test",
					Owner:    aliceKey.GetAddress(),
					NewAdmin: bobKey.GetAddress(),
				})
				if err != nil {
					t.Fatalf("handlerMsgTransferDomain() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks) {
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
				if emptyAcc, _ := k.GetAccount(ctx, "test", ""); !reflect.DeepEqual(emptyAcc, types2.Account{Domain: "test", Name: ""}) {
					t.Fatalf("handlerMsgTransferdomain() empty account mismatch, expected: %+v, got: %+v", types2.Account{Domain: "test", Name: ""}, emptyAcc)
				}
			},
		},
	}

	runTests(t, cases)
}

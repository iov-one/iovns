package domain

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovnsd"
	"github.com/iov-one/iovnsd/x/configuration"
	"github.com/iov-one/iovnsd/x/domain/keeper"
	"github.com/iov-one/iovnsd/x/domain/types"
	"testing"
	"time"
)

func Test_handleMsgRegisterAccount(t *testing.T) {
	testCases := map[string]subTest{
		"fail invalid blockchain targets address": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
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
				k.SetDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        nil,
					ValidUntil:   0,
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handleMsgRegisterAccount(ctx, k, types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  aliceKey.GetAddress(),
					Targets: []iovnsd.BlockchainAddress{
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
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
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
				k.SetDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        nil,
					ValidUntil:   0,
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handleMsgRegisterAccount(ctx, k, types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  aliceKey.GetAddress(),
					Targets: []iovnsd.BlockchainAddress{
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
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
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
				k.SetDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        nil,
					ValidUntil:   0,
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handleMsgRegisterAccount(ctx, k, types.MsgRegisterAccount{
					Domain: "test",
					Name:   "this won't match",
					Owner:  aliceKey.GetAddress(),
					Targets: []iovnsd.BlockchainAddress{
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
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
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
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handleMsgRegisterAccount(ctx, k, types.MsgRegisterAccount{
					Domain: "this does not exist",
					Name:   "works",
					Owner:  aliceKey.GetAddress(),
					Targets: []iovnsd.BlockchainAddress{
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
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
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
				k.SetDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        bobKey.GetAddress(),
					ValidUntil:   0,
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handleMsgRegisterAccount(ctx, k, types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  aliceKey.GetAddress(), // invalid owner
					Targets: []iovnsd.BlockchainAddress{
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
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
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
				k.SetDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        bobKey.GetAddress(),
					ValidUntil:   0, // domain is expired
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handleMsgRegisterAccount(ctx, k, types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  bobKey.GetAddress(),
					Targets: []iovnsd.BlockchainAddress{
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
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
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
				k.SetDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        bobKey.GetAddress(),
					ValidUntil:   time.Now().Add(100000 * time.Hour).Unix(),
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
				// add an account that we are gonna try to overwrite
				k.SetAccount(ctx, types.Account{
					Domain:       "test",
					Name:         "exists",
					Owner:        nil,
					ValidUntil:   0,
					Targets:      nil,
					Certificates: nil,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handleMsgRegisterAccount(ctx, k, types.MsgRegisterAccount{
					Domain: "test",
					Name:   "exists",
					Owner:  bobKey.GetAddress(),
					Targets: []iovnsd.BlockchainAddress{
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
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
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
				k.SetDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        bobKey.GetAddress(),
					ValidUntil:   time.Now().Add(100000 * time.Hour).Unix(),
					HasSuperuser: true,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handleMsgRegisterAccount(ctx, k, types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  bobKey.GetAddress(),
					Targets: []iovnsd.BlockchainAddress{
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

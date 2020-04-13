package domain

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/types"
	"testing"
)

func TestHandleMsgRegisterDomain(t *testing.T) {
	testCases := map[string]subTest{
		"success": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context) {
				configSetter := getConfigSetter(k.ConfigurationKeeper)
				// set config
				configSetter.SetConfig(ctx, configuration.Config{
					Owner:       aliceKey.GetAddress(),
					ValidDomain: "^(.*?)?",
				})
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context) {
				// register domain with superuser
				_, err := handleMsgRegisterDomain(ctx, k, types.MsgRegisterDomain{
					Name:         "domain",
					HasSuperuser: true,
					AccountRenew: 10,
				})
				if err != nil {
					t.Fatalf("handleMsgRegisterDomain() with superuser, got error: %s", err)
				}
				// register domain without super user
				_, err = handleMsgRegisterDomain(ctx, k, types.MsgRegisterDomain{
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
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context) {
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
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context) {
				k.SetDomain(ctx, types.Domain{
					Name:         "exists",
					Admin:        nil,
					ValidUntil:   0,
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context) {
				_, err := handleMsgRegisterDomain(ctx, k, types.MsgRegisterDomain{
					Name:         "exists",
					Admin:        nil,
					HasSuperuser: false,
					Broker:       nil,
					AccountRenew: 0,
				})
				if !errors.Is(err, types.ErrDomainAlreadyExists) {
					t.Fatalf("handleMsgRegisterDomain() expected: %s got: %s", types.ErrDomainAlreadyExists, err)
				}
			},
			AfterTest: nil,
		},
		"fail domain does not match valid domain regexp": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context) {
				// get set config function
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidDomain: "$^",
					DomainRenew: 0,
				})
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context) {
				_, err := handleMsgRegisterDomain(ctx, k, types.MsgRegisterDomain{
					Name:         "invalid-name",
					Admin:        nil,
					HasSuperuser: false,
					Broker:       nil,
					AccountRenew: 0,
				})
				if !errors.Is(err, types.ErrInvalidDomainName) {
					t.Fatalf("handleMsgRegisterDomain() expected error: %s, got: %s", types.ErrInvalidDomainName, err)
				}
			},
			AfterTest: nil,
		},
		"fail domain with no super user must be registered by configuration owner": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context) {
				// add config with owner
				config := configuration.Config{
					Owner:                  aliceKey.GetAddress(),
					ValidDomain:            "^(.*?)?",
					ValidName:              "",
					ValidBlockchainID:      "",
					ValidBlockchainAddress: "",
					DomainRenew:            0,
				}
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, config)
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context) {
				// try to register domain with no super user
				_, err := handleMsgRegisterDomain(ctx, k, types.MsgRegisterDomain{
					Name:         "some-domain",
					Admin:        bobKey.GetAddress(),
					HasSuperuser: false,
					Broker:       nil,
					AccountRenew: 10,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handleMsgRegisterDomain() expecter error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
	}
	// run all test cases
	runTests(t, testCases)
}

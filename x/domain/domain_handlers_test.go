package domain

import (
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

func Test_handleMsgDomainDelete(t *testing.T) {
	cases := map[string]dt.SubTest{
		"fail domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// don't do anything
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "this does not exist",
					Owner:  dt.BobKey,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgDeleteDomain() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
		},
		"fail domain has no superuser": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set domain with no superuser
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        iovns.ZeroAddress,
					ValidUntil:   0,
					Type:         types.OpenDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test",
					Owner:  dt.BobKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgDeleteDomain() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"fail domain admin does not match msg owner": {
			BeforeTestBlockTime: 0,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 1000000000000000,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        dt.BobKey,
					ValidUntil:   0,
					Type:         types.CloseDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			TestBlockTime: 1,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test",
					Owner:  dt.AliceKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgDeleteDomain() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"fail domain grace period not over": {
			BeforeTestBlockTime: 0,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 5,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        dt.BobKey,
					ValidUntil:   3,
					Type:         types.CloseDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			TestBlockTime: 3,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test",
					Owner:  dt.AliceKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgDeleteDomain() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success domain grace period over": {
			BeforeTestBlockTime: 0,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 5,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        dt.BobKey,
					ValidUntil:   4,
					Type:         types.CloseDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			TestBlockTime: 10,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test",
					Owner:  dt.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteDomain() got error: %s", err)
				}
			},
			AfterTest: nil,
		},
		"success owner can delete one of the domains after one expires and deleted": {
			BeforeTestBlockTime: 1589826438,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 1,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:         "test1",
					Admin:        dt.BobKey,
					ValidUntil:   1589826439,
					Type:         types.CloseDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:         "test2",
					Admin:        dt.BobKey,
					ValidUntil:   1589828251,
					Type:         types.CloseDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			TestBlockTime: 1589826441,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// another user can delete expired domain
				_, err := handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test1",
					Owner:  dt.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteDomain() got error: %s", err)
				}
				_, err = handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test2",
					Owner:  dt.AliceKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgDeleteDomain() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
				_, err = handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test2",
					Owner:  dt.BobKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteDomain() got error: %s", err)
				}
			},
			AfterTest: nil,
		},
		"success owner can delete their domain before grace period": {
			BeforeTestBlockTime: 0,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 1000000000000000, // unexpired domain
				})
				// set domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        dt.AliceKey,
					ValidUntil:   0,
					Type:         types.CloseDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			TestBlockTime: 4,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test",
					Owner:  dt.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteDomain() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
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
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 1,
				})
				// set domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        dt.AliceKey,
					ValidUntil:   0,
					Type:         types.CloseDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
				// add two accounts
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "1",
					Owner:  dt.BobKey,
				})
				// add two accounts
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "2",
					Owner:  dt.BobKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test",
					Owner:  dt.BobKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteDomain() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
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
	dt.RunTests(t, cases)
}

func TestHandleMsgRegisterDomain(t *testing.T) {
	testCases := map[string]dt.SubTest{
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				configSetter := dt.GetConfigSetter(k.ConfigurationKeeper)
				// set config
				configSetter.SetConfig(ctx, configuration.Config{
					Configurer:  dt.AliceKey,
					ValidDomain: "^(.*?)?",
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// register domain with superuser
				_, err := handleMsgRegisterDomain(ctx, k, &types.MsgRegisterDomain{
					Name:         "domain",
					DomainType:   types.CloseDomain,
					AccountRenew: 10,
					Admin:        dt.BobKey,
				})
				if err != nil {
					t.Fatalf("handleMsgRegisterDomain() with superuser, got error: %s", err)
				}
				// register domain without super user
				_, err = handleMsgRegisterDomain(ctx, k, &types.MsgRegisterDomain{
					Name:         "domain-without-superuser",
					Admin:        dt.AliceKey,
					DomainType:   types.OpenDomain,
					Broker:       nil,
					AccountRenew: 20,
				})
				if err != nil {
					t.Fatalf("handleMsgRegisterDomain() without superuser, got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
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
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:         "exists",
					Admin:        dt.BobKey,
					ValidUntil:   0,
					Type:         types.CloseDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterDomain(ctx, k, &types.MsgRegisterDomain{
					Name:         "exists",
					Admin:        dt.AliceKey,
					DomainType:   types.CloseDomain,
					AccountRenew: 0,
				})
				if !errors.Is(err, types.ErrDomainAlreadyExists) {
					t.Fatalf("handleMsgRegisterDomain() expected: %s got: %s", types.ErrDomainAlreadyExists, err)
				}
			},
			AfterTest: nil,
		},
		"fail domain does not match valid domain regexp": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// get set config function
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidDomain: "$^",
					DomainRenew: 0,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterDomain(ctx, k, &types.MsgRegisterDomain{
					Name:         "invalid-name",
					Admin:        nil,
					DomainType:   types.OpenDomain,
					Broker:       nil,
					AccountRenew: 0,
				})
				if !errors.Is(err, types.ErrInvalidDomainName) {
					t.Fatalf("handleMsgRegisterDomain() expected error: %s, got: %s", types.ErrInvalidDomainName, err)
				}
			},
			AfterTest: nil,
		},
		/* TODO REMOVE THIS TEST CASE AS IT DOES NOT APPLY WITH NEW SPEC.
		"fail domain with no super user must be registered by configuration owner": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// add config with owner
				config := configuration.Config{
					Configurer:             aliceKey,
					ValidDomain:            "^(.*?)?",
					ValidName:              "",
					ValidBlockchainID:      "",
					ValidBlockchainAddress: "",
					DomainRenew:            0,
				}
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, config)
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mock *keeper.Mocks) {
				// try to register domain with no super user
				_, err := handleMsgRegisterDomain(ctx, k, &sdk.MsgRegisterDomain{
					Name:         "some-domain",
					Admin:        bobKey,
					Type: types.OpenDomain,
					Broker:       nil,
					AccountRenew: 10,
				})
				if !errors.Is(err, sdk.ErrUnauthorized) {
					t.Fatalf("handleMsgRegisterDomain() expecter error: %s, got: %s", sdk.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		*/
	}
	// run all test cases
	dt.RunTests(t, testCases)
}

func Test_handlerDomainRenew(t *testing.T) {
	cases := map[string]dt.SubTest{
		"domain not found": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {

			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgRenewDomain(ctx, k, &types.MsgRenewDomain{Domain: "does not exist"})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgRenewDomain() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// add config
				setConfig := dt.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainRenew: 1 * time.Second,
				})
				// add domain
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: 1000,
					Admin:      dt.BobKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgRenewDomain(ctx, k, &types.MsgRenewDomain{Domain: "test"})
				if err != nil {
					t.Fatalf("handlerMsgRenewDomain() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// get domain
				domain, _ := k.GetDomain(ctx, "test")
				if domain.ValidUntil != 1001 {
					t.Fatalf("handlerMsgRenewDomain() expected 1001, got: %d", domain.ValidUntil)
				}
			},
		},
	}
	// run tests
	dt.RunTests(t, cases)
}

func Test_handlerMsgTransferDomain(t *testing.T) {
	cases := map[string]dt.SubTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {

			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferDomain(ctx, k, &types.MsgTransferDomain{
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
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:  "test",
					Type:  types.OpenDomain,
					Admin: iovns.ZeroAddress,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferDomain(ctx, k, &types.MsgTransferDomain{
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
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:  "test",
					Type:  types.CloseDomain,
					Admin: dt.BobKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferDomain(ctx, k, &types.MsgTransferDomain{
					Domain:   "test",
					Owner:    dt.BobKey,
					NewAdmin: nil,
				})
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("handlerMsgTransferDomain() expected error: %s, got error: %s", types.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"msg signer is not domain admin": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					Type:       types.CloseDomain,
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      dt.AliceKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferDomain(ctx, k, &types.MsgTransferDomain{
					Domain:   "test",
					Owner:    dt.BobKey,
					NewAdmin: nil,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgTransferDomain() expected error: %s, got error: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// create domain
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					Type:       types.CloseDomain,
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      dt.AliceKey,
				})
				// add empty account
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "",
					Owner:  dt.AliceKey,
				})
				// add account 1
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "1",
					Owner:      dt.AliceKey,
					ValidUntil: 0,
					Targets: []types.BlockchainAddress{{
						ID:      "test",
						Address: "test",
					}},
					Certificates: []types.Certificate{[]byte("cert")},
					Broker:       nil,
				})
				// add account 2
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "2",
					Owner:      dt.AliceKey,
					ValidUntil: 0,
					Targets: []types.BlockchainAddress{{
						ID:      "test",
						Address: "test",
					}},
					Certificates: []types.Certificate{[]byte("cert")},
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferDomain(ctx, k, &types.MsgTransferDomain{
					Domain:   "test",
					Owner:    dt.AliceKey,
					NewAdmin: dt.BobKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgTransferDomain() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// check domain new owner
				domain, _ := k.GetDomain(ctx, "test")
				if !dt.BobKey.Equals(domain.Admin) {
					t.Fatalf("handlerMsgTransferDomain() expected domain owner: %s, got: %s", dt.BobKey, domain.Admin)
				}
				// check if account new owner has changed
				account, _ := k.GetAccount(ctx, "test", "1")
				if !account.Owner.Equals(dt.BobKey) {
					t.Fatalf("handlerMsgTransferDomain() expected account owner: %s, got: %s", dt.BobKey, account.Owner)
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
				if emptyAcc, _ := k.GetAccount(ctx, "test", ""); !reflect.DeepEqual(emptyAcc, types.Account{Domain: "test", Name: "", Owner: dt.AliceKey}) {
					t.Fatalf("handlerMsgTransferdomain() empty account mismatch, expected: %+v, got: %+v", types.Account{Domain: "test", Name: ""}, emptyAcc)
				}
			},
		},
	}

	dt.RunTests(t, cases)
}

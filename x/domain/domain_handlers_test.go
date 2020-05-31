package domain

import (
	"errors"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

func Test_handleMsgDomainDelete(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"fail domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// don't do anything
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "this does not exist",
					Owner:  keeper.BobKey,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgDeleteDomain() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
		},
		"fail domain open": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        keeper.AliceKey,
					ValidUntil:   0,
					Type:         types.OpenDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test",
					Owner:  keeper.BobKey,
				})
				if !errors.Is(err, types.ErrInvalidDomainType) {
					t.Fatalf("handlerMsgDeleteDomain() expected error: %s, got: %s", types.ErrInvalidDomainType, err)
				}
			},
			AfterTest: nil,
		},
		"fail domain admin does not match msg owner": {
			BeforeTestBlockTime: 0,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 1000000000000000,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        keeper.BobKey,
					ValidUntil:   0,
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			TestBlockTime: 1,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test",
					Owner:  keeper.AliceKey,
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
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 5,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        keeper.BobKey,
					ValidUntil:   3,
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			TestBlockTime: 3,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test",
					Owner:  keeper.AliceKey,
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
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 5,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        keeper.BobKey,
					ValidUntil:   4,
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			TestBlockTime: 10,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test",
					Owner:  keeper.AliceKey,
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
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 1,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:         "test1",
					Admin:        keeper.BobKey,
					ValidUntil:   1589826439,
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
				k.CreateDomain(ctx, types.Domain{
					Name:         "test2",
					Admin:        keeper.BobKey,
					ValidUntil:   1589828251,
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			TestBlockTime: 1589826441,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// another user can delete expired domain
				_, err := handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test1",
					Owner:  keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteDomain() got error: %s", err)
				}
				_, err = handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test2",
					Owner:  keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgDeleteDomain() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
				_, err = handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test2",
					Owner:  keeper.BobKey,
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
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 1000000000000000, // unexpired domain
				})
				// set domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        keeper.AliceKey,
					ValidUntil:   0,
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			TestBlockTime: 4,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test",
					Owner:  keeper.AliceKey,
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
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 1,
				})
				// set domain
				k.CreateDomain(ctx, types.Domain{
					Name:         "test",
					Admin:        keeper.AliceKey,
					ValidUntil:   0,
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
				// add two accounts
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "1",
					Owner:  keeper.BobKey,
				})
				// add two accounts
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "2",
					Owner:  keeper.BobKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test",
					Owner:  keeper.BobKey,
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
	keeper.RunTests(t, cases)
}

func TestHandleMsgRegisterDomain(t *testing.T) {
	testCases := map[string]keeper.SubTest{
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				configSetter := keeper.GetConfigSetter(k.ConfigurationKeeper)
				// set config
				configSetter.SetConfig(ctx, configuration.Config{
					Configurer:      keeper.AliceKey,
					ValidDomainName: "^(.*?)?",
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterDomain(ctx, k, &types.MsgRegisterDomain{
					Name:         "domain-closed",
					DomainType:   types.ClosedDomain,
					AccountRenew: 10,
					Admin:        keeper.BobKey,
				})
				if err != nil {
					t.Fatalf("handleMsgRegisterDomain() with close domain, got error: %s", err)
				}
				_, err = handleMsgRegisterDomain(ctx, k, &types.MsgRegisterDomain{
					Name:         "domain-open",
					Admin:        keeper.AliceKey,
					DomainType:   types.OpenDomain,
					Broker:       nil,
					AccountRenew: 20,
				})
				if err != nil {
					t.Fatalf("handleMsgRegisterDomain() with open domain, got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// TODO do reflect.DeepEqual checks on expected results vs results returned
				_, ok := k.GetDomain(ctx, "domain-closed")
				if !ok {
					t.Fatalf("handleMsgRegisterDomain() could not find 'domain-closed'")
				}
				_, ok = k.GetDomain(ctx, "domain-open")
				if !ok {
					t.Fatalf("handleMsgRegisterDomain() could not find 'domain-open'")
				}
			},
		},
		"fail domain name exists": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:         "exists",
					Admin:        keeper.BobKey,
					ValidUntil:   0,
					Type:         types.ClosedDomain,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterDomain(ctx, k, &types.MsgRegisterDomain{
					Name:         "exists",
					Admin:        keeper.AliceKey,
					DomainType:   types.ClosedDomain,
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
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidDomainName:     "$^",
					DomainRenewalPeriod: 0,
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
	}
	// run all test cases
	keeper.RunTests(t, testCases)
}

func Test_handlerDomainRenew(t *testing.T) {
	cases := map[string]keeper.SubTest{
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
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainRenewalPeriod: 1 * time.Second,
				})
				// add domain
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: 1000,
					Admin:      keeper.BobKey,
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
	keeper.RunTests(t, cases)
}

func Test_handlerMsgTransferDomain(t *testing.T) {
	cases := map[string]keeper.SubTest{
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
		"domain type open": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:  "test",
					Type:  types.OpenDomain,
					Admin: keeper.AliceKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferDomain(ctx, k, &types.MsgTransferDomain{
					Domain:   "test",
					Owner:    nil,
					NewAdmin: nil,
				})
				if !errors.Is(err, types.ErrInvalidDomainType) {
					t.Fatalf("handlerMsgTransferDomain() expected error: %s, got error: %s", types.ErrInvalidDomainType, err)
				}
			},
			AfterTest: nil,
		},
		"domain type closed": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				k.CreateDomain(ctx, types.Domain{
					Name:  "test",
					Type:  types.ClosedDomain,
					Admin: keeper.AliceKey,
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
					Type:  types.ClosedDomain,
					Admin: keeper.BobKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferDomain(ctx, k, &types.MsgTransferDomain{
					Domain:   "test",
					Owner:    keeper.BobKey,
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
					Type:       types.ClosedDomain,
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      keeper.AliceKey,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferDomain(ctx, k, &types.MsgTransferDomain{
					Domain:   "test",
					Owner:    keeper.BobKey,
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
					Type:       types.ClosedDomain,
					ValidUntil: iovns.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      keeper.AliceKey,
				})
				// add empty account
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "",
					Owner:  keeper.AliceKey,
				})
				// add account 1
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "1",
					Owner:      keeper.AliceKey,
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
					Owner:      keeper.AliceKey,
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
					Domain:       "test",
					Owner:        keeper.AliceKey,
					NewAdmin:     keeper.BobKey,
					TransferFlag: types.TransferOwned,
				})
				if err != nil {
					t.Fatalf("handlerMsgTransferDomain() got error: %s", err)
				}
			},
		},
	}

	keeper.RunTests(t, cases)
}

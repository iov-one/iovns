package starname

import (
	"errors"
	"github.com/iov-one/iovns/pkg/utils"
	"github.com/iov-one/iovns/x/starname/keeper/executor"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/starname/keeper"
	"github.com/iov-one/iovns/x/starname/types"
)

func Test_Closed_handleMsgDomainDelete(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"success only admin can delete before grace period": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 10 * time.Second,
				})
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					ValidUntil: 2,
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("1"),
					Owner:  keeper.BobKey,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("2"),
					Owner:  keeper.BobKey,
				}).Create()
			},
			TestBlockTime: 3,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test",
					Owner:  keeper.BobKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("unexpected error: %s", err)
				}
				_, err = handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test",
					Owner:  keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteDomain() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				exists := k.DomainStore(ctx).Read((&types.Domain{Name: "test"}).PrimaryKey(), new(types.Domain))
				if exists {
					t.Fatalf("handlerMsgDeleteDomain() domain should not exist")
				}
				exists = k.AccountStore(ctx).Read((&types.Account{Domain: "test", Name: utils.StrPtr("1")}).PrimaryKey(), new(types.Account))
				if exists {
					t.Fatalf("handlerMsgDeleteDomain() account 1 should not exist")
				}
				exists = k.AccountStore(ctx).Read((&types.Account{Domain: "test", Name: utils.StrPtr("2")}).PrimaryKey(), new(types.Account))
				if exists {
					t.Fatalf("handlerMsgDeleteDomain() account 2 should not exist")
				}
			},
		},
		"success anyone can after grace period": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 10 * time.Second,
				})
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					ValidUntil: 2,
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("1"),
					Owner:  keeper.BobKey,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("2"),
					Owner:  keeper.BobKey,
				}).Create()
			},
			TestBlockTime: 1000,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test",
					Owner:  keeper.CharlieKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteDomain() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				exists := k.DomainStore(ctx).Read((&types.Domain{Name: "test"}).PrimaryKey(), new(types.Domain))
				if exists {
					t.Fatalf("handlerMsgDeleteDomain() domain should not exist")
				}
				exists = k.AccountStore(ctx).Read((&types.Account{Domain: "test", Name: utils.StrPtr("1")}).PrimaryKey(), new(types.Account))
				if exists {
					t.Fatalf("handlerMsgDeleteDomain() account 1 should not exist")
				}
				exists = k.AccountStore(ctx).Read((&types.Account{Domain: "test", Name: utils.StrPtr("2")}).PrimaryKey(), new(types.Account))
				if exists {
					t.Fatalf("handlerMsgDeleteDomain() account 2 should not exist")
				}
			},
		},
	}
	keeper.RunTests(t, cases)
}

func Test_Open_handleMsgDomainDelete(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"success anyone can after grace period": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 10 * time.Second,
				})
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					ValidUntil: 2,
					Type:       types.OpenDomain,
					Broker:     nil,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("1"),
					Owner:  keeper.BobKey,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("2"),
					Owner:  keeper.BobKey,
				}).Create()
			},
			TestBlockTime: 1000,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test",
					Owner:  keeper.CharlieKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteDomain() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				exists := k.DomainStore(ctx).Read((&types.Domain{Name: "test"}).PrimaryKey(), new(types.Domain))
				if exists {
					t.Fatalf("handlerMsgDeleteDomain() domain should not exist")
				}
				exists = k.AccountStore(ctx).Read((&types.Account{Domain: "test", Name: utils.StrPtr("1")}).PrimaryKey(), new(types.Account))
				if exists {
					t.Fatalf("handlerMsgDeleteDomain() account 1 should not exist")
				}
				exists = k.AccountStore(ctx).Read((&types.Account{Domain: "test", Name: utils.StrPtr("2")}).PrimaryKey(), new(types.Account))
				if exists {
					t.Fatalf("handlerMsgDeleteDomain() account 2 should not exist")
				}
			},
		},
		"domain cannot be deleted before grace period even by admin": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 10 * time.Second,
				})
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					ValidUntil: 2,
					Type:       types.OpenDomain,
					Broker:     nil,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("1"),
					Owner:  keeper.BobKey,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("2"),
					Owner:  keeper.BobKey,
				}).Create()
			},
			TestBlockTime: 3,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test",
					Owner:  keeper.CharlieKey,
				})
				if !errors.Is(err, types.ErrDomainGracePeriodNotFinished) {
					t.Fatalf("unexpected error: %s", err)
				}
				_, err = handlerMsgDeleteDomain(ctx, k, &types.MsgDeleteDomain{
					Domain: "test",
					Owner:  keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrDomainGracePeriodNotFinished) {
					t.Fatalf("unexpected error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				exists := k.DomainStore(ctx).Read((&types.Domain{Name: "test"}).PrimaryKey(), new(types.Domain))
				if !exists {
					t.Fatalf("handlerMsgDeleteDomain() domain should exist")
				}
				exists = k.AccountStore(ctx).Read((&types.Account{Domain: "test", Name: utils.StrPtr("1")}).PrimaryKey(), new(types.Account))
				if !exists {
					t.Fatalf("handlerMsgDeleteDomain() account 1 should exist")
				}
				exists = k.AccountStore(ctx).Read((&types.Account{Domain: "test", Name: utils.StrPtr("2")}).PrimaryKey(), new(types.Account))
				if !exists {
					t.Fatalf("handlerMsgDeleteDomain() account 2 should exist")
				}
			},
		},
	}
	keeper.RunTests(t, cases)
}

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
		"fail domain admin does not match msg owner": {
			BeforeTestBlockTime: 0,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 1000000000000000,
				})
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					ValidUntil: 0,
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
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
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					ValidUntil: 3,
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
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
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					ValidUntil: 4,
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
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
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test1",
					Admin:      keeper.BobKey,
					ValidUntil: 1589826439,
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test2",
					Admin:      keeper.BobKey,
					ValidUntil: 1589828251,
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
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
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {

			},
		},
		"success owner can delete their domain before grace period": {
			BeforeTestBlockTime: 0,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainGracePeriod: 1000000000000000, // unexpired domain
				})
				// set domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					ValidUntil: 0,
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
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
				exists := k.DomainStore(ctx).Read((&types.Domain{Name: "test"}).PrimaryKey(), new(types.Domain))
				if exists {
					t.Fatalf("handlerMsgDeleteDomain() domain should not exist")
				}
				exists = k.AccountStore(ctx).Read((&types.Account{Domain: "test", Name: utils.StrPtr("1")}).PrimaryKey(), new(types.Account))
				if exists {
					t.Fatalf("handlerMsgDeleteDomain() account 1 should not exist")
				}
				exists = k.AccountStore(ctx).Read((&types.Account{Domain: "test", Name: utils.StrPtr("2")}).PrimaryKey(), new(types.Account))
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
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					ValidUntil: 0,
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
				// add two accounts
				executor.NewAccount(ctx, k, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("1"),
					Owner:  keeper.BobKey,
				}).Create()
				// add two accounts
				executor.NewAccount(ctx, k, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("2"),
					Owner:  keeper.BobKey,
				}).Create()
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
				exists := k.DomainStore(ctx).Read((&types.Domain{Name: "test"}).PrimaryKey(), new(types.Domain))
				if exists {
					t.Fatalf("handlerMsgDeleteDomain() domain should not exist")
				}
				exists = k.AccountStore(ctx).Read((&types.Account{Domain: "test", Name: utils.StrPtr("1")}).PrimaryKey(), new(types.Account))
				if exists {
					t.Fatalf("handlerMsgDeleteDomain() account 1 should not exist")
				}
				exists = k.AccountStore(ctx).Read((&types.Account{Domain: "test", Name: utils.StrPtr("2")}).PrimaryKey(), new(types.Account))
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
					Name:       "domain-closed",
					DomainType: types.ClosedDomain,
					Admin:      keeper.BobKey,
				})
				if err != nil {
					t.Fatalf("handleMsgRegisterDomain() with close domain, got error: %s", err)
				}
				_, err = handleMsgRegisterDomain(ctx, k, &types.MsgRegisterDomain{
					Name:       "domain-open",
					Admin:      keeper.AliceKey,
					DomainType: types.OpenDomain,
					Broker:     nil,
				})
				if err != nil {
					t.Fatalf("handleMsgRegisterDomain() with open domain, got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// TODO do reflect.DeepEqual checks on expected results vs results returned
				exists := k.DomainStore(ctx).Read((&types.Domain{Name: "domain-closed"}).PrimaryKey(), new(types.Domain))
				if !exists {
					t.Fatalf("handleMsgRegisterDomain() could not find 'domain-closed'")
				}
				exists = k.DomainStore(ctx).Read((&types.Domain{Name: "domain-open"}).PrimaryKey(), new(types.Domain))
				if !exists {
					t.Fatalf("handleMsgRegisterDomain() could not find 'domain-open'")
				}
			},
		},
		"fail domain name exists": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "exists",
					Admin:      keeper.BobKey,
					ValidUntil: 0,
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterDomain(ctx, k, &types.MsgRegisterDomain{
					Name:       "exists",
					Admin:      keeper.AliceKey,
					DomainType: types.ClosedDomain,
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
					Name:       "invalid-name",
					Admin:      nil,
					DomainType: types.OpenDomain,
					Broker:     nil,
				})
				if !errors.Is(err, types.ErrInvalidDomainName) {
					t.Fatalf("handleMsgRegisterDomain() expected error: %s, got: %s", types.ErrInvalidDomainName, err)
				}
			},
			// TODO ADD AFTER TEST
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {

			},
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
			BeforeTestBlockTime: 1000,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// add config
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainRenewalCountMax: 2,
					DomainRenewalPeriod:   1 * time.Second,
					DomainGracePeriod:     10 * time.Second,
				})
				// add domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: 1000,
					Admin:      keeper.BobKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgRenewDomain(ctx, k, &types.MsgRenewDomain{Domain: "test"})
				if err != nil {
					t.Fatalf("handlerMsgRenewDomain() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// get domain
				domain := new(types.Domain)
				_ = k.DomainStore(ctx).Read((&types.Domain{Name: "test"}).PrimaryKey(), domain)
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
				executor.NewDomain(ctx, k, types.Domain{
					Name:  "test",
					Type:  types.OpenDomain,
					Admin: keeper.AliceKey,
				}).Create()
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
				executor.NewDomain(ctx, k, types.Domain{
					Name:  "test",
					Type:  types.ClosedDomain,
					Admin: keeper.AliceKey,
				}).Create()
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
				executor.NewDomain(ctx, k, types.Domain{
					Name:  "test",
					Type:  types.ClosedDomain,
					Admin: keeper.BobKey,
				}).Create()
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
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Type:       types.ClosedDomain,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      keeper.AliceKey,
				}).Create()
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
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Type:       types.ClosedDomain,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      keeper.AliceKey,
				}).Create()
				// add account 1
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("1"),
					Owner:      keeper.AliceKey,
					ValidUntil: 0,
					Resources: []types.Resource{{
						URI:      "test",
						Resource: "test",
					}},
					Certificates: []types.Certificate{[]byte("cert")},
					Broker:       nil,
				}).Create()
				// add account 2
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("2"),
					Owner:      keeper.AliceKey,
					ValidUntil: 0,
					Resources: []types.Resource{{
						URI:      "test",
						Resource: "test",
					}},
					Certificates: []types.Certificate{[]byte("cert")},
					Broker:       nil,
				}).Create()
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
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {

			},
		},
	}

	keeper.RunTests(t, cases)
}

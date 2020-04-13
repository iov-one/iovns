package domain

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
	"reflect"
	"testing"
	"time"
)

func Test_handlerMsgReplaceAccountTargets(t *testing.T) {
	cases := map[string]subTest{
		"invalid blockchain target": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// set config to match nothing
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      regexMatchNothing,
					ValidBlockchainAddress: regexMatchNothing,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, types.MsgReplaceAccountTargets{
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
				if !errors.Is(err, types.ErrInvalidBlockchainTarget) {
					t.Fatalf("handlerMsgReplaceAccountTargets() expected error: %s, got: %s", types.ErrInvalidBlockchainTarget, err)
				}
			},
			AfterTest: nil,
		},
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// set config to match all
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      regexMatchAll,
					ValidBlockchainAddress: regexMatchAll,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, types.MsgReplaceAccountTargets{
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
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgReplaceAccountTargets() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// set config to match all
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      regexMatchAll,
					ValidBlockchainAddress: regexMatchAll,
				})
				// create domain
				k.SetDomain(ctx, types.Domain{
					Name: "test",
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, types.MsgReplaceAccountTargets{
					Domain: "test",
					NewTargets: []iovns.BlockchainAddress{
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
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// set config to match all
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      regexMatchAll,
					ValidBlockchainAddress: regexMatchAll,
				})
				// create domain
				k.SetDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, types.MsgReplaceAccountTargets{
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
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgReplaceAccountTargets() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"account expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// set config to match all
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      regexMatchAll,
					ValidBlockchainAddress: regexMatchAll,
				})
				// create domain
				k.SetDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
				})
				// create account
				k.SetAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: 0,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, types.MsgReplaceAccountTargets{
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
				if !errors.Is(err, types.ErrAccountExpired) {
					t.Fatalf("handlerMsgReplaceAccountTargets() expected error: %s, got: %s", types.ErrAccountExpired, err)
				}
			},
			AfterTest: nil,
		},
		"signer is not owner of account": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// set config to match all
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      regexMatchAll,
					ValidBlockchainAddress: regexMatchAll,
				})
				// create domain
				k.SetDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
				})
				// create account
				k.SetAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      aliceKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, types.MsgReplaceAccountTargets{
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
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgReplaceAccountTargets() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// set config to match all
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidBlockchainID:      regexMatchAll,
					ValidBlockchainAddress: regexMatchAll,
				})
				// create domain
				k.SetDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
				})
				// create account
				k.SetAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "test",
					ValidUntil: iovns.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      aliceKey.GetAddress(),
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgReplaceAccountTargets(ctx, k, types.MsgReplaceAccountTargets{
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
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				expected := []iovns.BlockchainAddress{{
					ID:      "valid",
					Address: "valid",
				}}
				account, _ := k.GetAccount(ctx, iovns.GetAccountKey("test", "test"))
				if !reflect.DeepEqual(expected, account.Targets) {
					t.Fatalf("handlerMsgReplaceAccountTargets() expected: %+v, got %+v", expected, account.Targets)
				}
			},
		},
	}
	// run tests
	runTests(t, cases)
}

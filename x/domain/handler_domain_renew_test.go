package domain

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
	"testing"
)

func Test_handlerDomainRenew(t *testing.T) {
	cases := map[string]subTest{
		"domain not found": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {

			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgRenewDomain(ctx, k, types.MsgRenewDomain{Domain: "does not exist"})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgRenewDomain() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// add config
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					DomainRenew: 1,
				})
				// add domain
				k.CreateDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: 1000,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				_, err := handlerMsgRenewDomain(ctx, k, types.MsgRenewDomain{Domain: "test"})
				if err != nil {
					t.Fatalf("handlerMsgRenewDomain() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
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

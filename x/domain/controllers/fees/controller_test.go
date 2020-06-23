package fees

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

func Test_FeeApplier(t *testing.T) {
	cases := map[string]struct {
		FeeConfig   *configuration.Fees
		Msg         sdk.Msg
		Domain      types.Domain
		ExpectedFee sdk.Int
	}{
		"register open domain length 5 coin price 1": {
			FeeConfig: &configuration.Fees{
				FeeCoinDenom:    "tiov",
				FeeCoinPrice:    sdk.NewDec(1),
				RegisterDomain5: sdk.NewDec(5),
				FeeDefault:      sdk.NewDec(4),
			},
			Msg: &types.MsgRegisterDomain{
				Name:       "test1",
				DomainType: types.OpenDomain,
			},
			ExpectedFee: sdk.NewInt(1),
		},
	}

	for _, c := range cases {
		k, ctx, _ := keeper.NewTestKeeper(t, true)
		k.ConfigurationKeeper.(keeper.ConfigurationSetter).SetFees(ctx, c.FeeConfig)
		ctrl := NewController(ctx, k, c.Domain)
		got := ctrl.GetFee(c.Msg)
		if got.Amount != c.ExpectedFee {
			t.Fatalf("expected fee: %s, got %s", c.ExpectedFee, got.Amount)
		}

	}
}

package fees

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

func Test_FeeApplier(t *testing.T) {
	dec1, err := sdk.NewDecFromStr("0.000001000000000000")
	if err != nil {
		t.Fatal(err)
	}
	dec2, err := sdk.NewDecFromStr("572.332205500000000000")
	if err != nil {
		t.Fatal(err)
	}
	dec3, err := sdk.NewDecFromStr("0.572332205500000100")
	if err != nil {
		t.Fatal(err)
	}
	expect, err := sdk.NewDecFromStr("572332205")
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string]struct {
		FeeConfig   *configuration.Fees
		Msg         sdk.Msg
		Domain      types.Domain
		ExpectedFee sdk.Dec
	}{
		"register open domain length 5 coin price 5": {
			FeeConfig: &configuration.Fees{
				FeeCoinDenom:    "tiov",
				FeeCoinPrice:    dec1,
				RegisterDomain5: dec2,
				FeeDefault:      dec3,
			},
			Msg: &types.MsgRegisterDomain{
				Name:       "test1",
				DomainType: types.ClosedDomain,
			},
			Domain: types.Domain{
				Name: "test1",
			},
			ExpectedFee: expect,
		},
	}

	for _, c := range cases {
		k, ctx, _ := keeper.NewTestKeeper(t, true)
		k.ConfigurationKeeper.(keeper.ConfigurationSetter).SetFees(ctx, c.FeeConfig)
		ctrl := NewController(ctx, k, c.Domain)
		got := ctrl.GetFee(c.Msg)
		if !got.Amount.Equal(expect.RoundInt()) {
			t.Fatalf("expected fee: %s, got %s", c.ExpectedFee, got.Amount)
		}
	}
}

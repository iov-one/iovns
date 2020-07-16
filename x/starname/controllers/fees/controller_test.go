package fees

import (
	"github.com/iov-one/iovns/tutils"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/starname/keeper"
	"github.com/iov-one/iovns/x/starname/types"
)

// decFromStr is a helper to convert string decimals such as 0.12311 easily
func decFromStr(str string) sdk.Dec {
	dec, err := sdk.NewDecFromStr(str)
	if err != nil {
		panic(err)
	}
	return dec
}

func Test_FeeApplier(t *testing.T) {
	fee := configuration.Fees{
		FeeCoinDenom:                 "tiov",
		FeeCoinPrice:                 decFromStr("2"),
		FeeDefault:                   decFromStr("2"),
		RegisterAccountClosed:        decFromStr("4"),
		RegisterAccountOpen:          sdk.NewDec(6),
		TransferAccountClosed:        sdk.NewDec(8),
		TransferAccountOpen:          sdk.NewDec(10),
		ReplaceAccountResources:      sdk.NewDec(12),
		AddAccountCertificate:        sdk.NewDec(14),
		DelAccountCertificate:        sdk.NewDec(16),
		SetAccountMetadata:           sdk.NewDec(18),
		RegisterDomain1:              sdk.NewDec(20),
		RegisterDomain2:              sdk.NewDec(22),
		RegisterDomain3:              sdk.NewDec(24),
		RegisterDomain4:              sdk.NewDec(26),
		RegisterDomain5:              sdk.NewDec(28),
		RegisterDomainDefault:        sdk.NewDec(30),
		RegisterOpenDomainMultiplier: sdk.NewDec(2),
		TransferDomainClosed:         sdk.NewDec(34),
		TransferDomainOpen:           sdk.NewDec(36),
		RenewDomainOpen:              sdk.NewDec(28),
	}
	cases := map[string]struct {
		Msg         sdk.Msg
		Domain      types.Domain
		ExpectedFee sdk.Dec
	}{
		"register closed domain 5": {
			Msg: &types.MsgRegisterDomain{},
			Domain: types.Domain{
				Name: "test1",
			},
			ExpectedFee: sdk.NewDec(14),
		},
		"register closed domain 4": {
			Msg: &types.MsgRegisterDomain{},
			Domain: types.Domain{
				Name: "test",
			},
			ExpectedFee: sdk.NewDec(13),
		},
		"register closed domain 3": {
			Msg: &types.MsgRegisterDomain{},
			Domain: types.Domain{
				Name: "tes",
			},
			ExpectedFee: sdk.NewDec(12),
		},
		"register closed domain 2": {
			Msg: &types.MsgRegisterDomain{},
			Domain: types.Domain{
				Name: "te",
			},
			ExpectedFee: sdk.NewDec(11),
		},
		"register closed domain 1": {
			Msg: &types.MsgRegisterDomain{},
			Domain: types.Domain{
				Name: "t",
			},
			ExpectedFee: sdk.NewDec(10),
		},
		"register open domain 5": {
			Msg: &types.MsgRegisterDomain{},
			Domain: types.Domain{
				Name: "test1",
				Type: types.OpenDomain,
			},
			ExpectedFee: sdk.NewDec(28),
		},
		"register open domain default": {
			Msg: &types.MsgRegisterDomain{},
			Domain: types.Domain{
				Name: "test12",
				Type: types.ClosedDomain,
			},
			ExpectedFee: sdk.NewDec(15),
		},
		"transfer domain open": {
			Msg:         &types.MsgTransferDomain{},
			Domain:      types.Domain{Type: types.OpenDomain},
			ExpectedFee: sdk.NewDec(18),
		},
		"transfer domain closed": {
			Msg:         &types.MsgTransferDomain{},
			Domain:      types.Domain{Type: types.ClosedDomain},
			ExpectedFee: sdk.NewDec(17),
		},
		"set metadata": {
			Msg:         &types.MsgReplaceAccountMetadata{},
			ExpectedFee: sdk.NewDec(9),
		},
		"delete certs": {
			Msg:         &types.MsgDeleteAccountCertificate{},
			ExpectedFee: sdk.NewDec(8),
		},
		"add certs": {
			Msg:         &types.MsgAddAccountCertificates{},
			ExpectedFee: sdk.NewDec(7),
		},
		"replace resources": {
			Msg:         &types.MsgReplaceAccountResources{},
			ExpectedFee: sdk.NewDec(6),
		},
		"transfer account closed": {
			Msg:         &types.MsgTransferAccount{},
			Domain:      types.Domain{Type: types.ClosedDomain},
			ExpectedFee: sdk.NewDec(4),
		},
		"transfer account open": {
			Msg:         &types.MsgTransferAccount{},
			Domain:      types.Domain{Type: types.OpenDomain},
			ExpectedFee: sdk.NewDec(5),
		},
		"register account open": {
			Msg:         &types.MsgRegisterAccount{},
			Domain:      types.Domain{Type: types.OpenDomain},
			ExpectedFee: sdk.NewDec(3),
		},
		"register account closed": {
			Msg:         &types.MsgRegisterAccount{},
			Domain:      types.Domain{Type: types.ClosedDomain},
			ExpectedFee: sdk.NewDec(2),
		},
		"renew account closed": {
			Msg:         &types.MsgRenewAccount{},
			Domain:      types.Domain{Type: types.ClosedDomain},
			ExpectedFee: sdk.NewDec(2),
		},
		"renew account open": {
			Msg:         &types.MsgRenewAccount{},
			Domain:      types.Domain{Type: types.OpenDomain},
			ExpectedFee: sdk.NewDec(3),
		},
		"renew domain open": {
			Msg:         &types.MsgRenewDomain{},
			Domain:      types.Domain{Type: types.OpenDomain},
			ExpectedFee: sdk.NewDec(14),
		},
		"renew domain closed": {
			Msg:         &types.MsgRenewDomain{},
			Domain:      types.Domain{Type: types.ClosedDomain, Name: "renew"},
			ExpectedFee: sdk.NewDec(6), // it's three accounts-> "", "1", "2"; so 4/2*3=6
		},
		"default fee unknown message": {
			Msg:         &keeper.DullMsg{},
			ExpectedFee: sdk.NewDec(1),
		},
		"use default fee": {
			Msg:         &types.MsgRenewDomain{},
			Domain:      types.Domain{Type: types.ClosedDomain, Name: "not exists"}, // since it does not exist, the fee is 0
			ExpectedFee: sdk.NewDec(1),
		},
	}
	k, ctx, _ := keeper.NewTestKeeper(t, true)
	ds := k.DomainStore(ctx)
	as := k.AccountStore(ctx)
	ds.Create(&types.Domain{Name: "renew", Admin: keeper.AliceKey})
	as.Create(&types.Account{Domain: "renew", Name: tutils.StrPtr(types.EmptyAccountName), Owner: keeper.AliceKey}) // TODO in the future this might be removed
	as.Create(&types.Account{Domain: "renew", Name: tutils.StrPtr("1"), Owner: keeper.AliceKey})
	as.Create(&types.Account{Domain: "renew", Name: tutils.StrPtr("2"), Owner: keeper.AliceKey})

	k.ConfigurationKeeper.(keeper.ConfigurationSetter).SetFees(ctx, &fee)
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl := NewController(ctx, k, c.Domain)
			got := ctrl.GetFee(c.Msg)
			if !got.Amount.Equal(c.ExpectedFee.RoundInt()) {
				t.Fatalf("expected fee: %s, got %s", c.ExpectedFee, got.Amount)
			}
		})
	}
}

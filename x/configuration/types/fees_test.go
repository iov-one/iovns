package types

import (
	"github.com/cosmos/cosmos-sdk/types"
	"testing"
)

func TestFees_Validate(t *testing.T) {
	type fields struct {
		IovTokenPrice                types.Dec
		DefaultFee                   types.Coin
		RegisterClosedAccount        types.Coin
		RegisterOpenAccount          types.Coin
		TransferClosedAccount        types.Coin
		TransferOpenAccount          types.Coin
		ReplaceAccountTargets        types.Coin
		AddAccountCertificate        types.Coin
		DelAccountCertificate        types.Coin
		SetAccountMetadata           types.Coin
		RegisterDomain               types.Coin
		RegisterDomain1              types.Coin
		RegisterDomain2              types.Coin
		RegisterDomain3              types.Coin
		RegisterDomain4              types.Coin
		RegisterDomain5              types.Coin
		RegisterDomainDefault        types.Coin
		RegisterOpenDomainMultiplier types.Int
		TransferDomainClosed         types.Coin
		TransferDomainOpen           types.Coin
		RenewOpenDomain              types.Coin
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "fail field not set",
			fields: fields{
				IovTokenPrice: types.NewDec(10),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fees{
				IovTokenPrice:                tt.fields.IovTokenPrice,
				DefaultFee:                   tt.fields.DefaultFee,
				RegisterClosedAccount:        tt.fields.RegisterClosedAccount,
				RegisterOpenAccount:          tt.fields.RegisterOpenAccount,
				TransferClosedAccount:        tt.fields.TransferClosedAccount,
				TransferOpenAccount:          tt.fields.TransferOpenAccount,
				ReplaceAccountTargets:        tt.fields.ReplaceAccountTargets,
				AddAccountCertificate:        tt.fields.AddAccountCertificate,
				DelAccountCertificate:        tt.fields.DelAccountCertificate,
				SetAccountMetadata:           tt.fields.SetAccountMetadata,
				RegisterDomain:               tt.fields.RegisterDomain,
				RegisterDomain1:              tt.fields.RegisterDomain1,
				RegisterDomain2:              tt.fields.RegisterDomain2,
				RegisterDomain3:              tt.fields.RegisterDomain3,
				RegisterDomain4:              tt.fields.RegisterDomain4,
				RegisterDomain5:              tt.fields.RegisterDomain5,
				RegisterDomainDefault:        tt.fields.RegisterDomainDefault,
				RegisterOpenDomainMultiplier: tt.fields.RegisterOpenDomainMultiplier,
				TransferDomainClosed:         tt.fields.TransferDomainClosed,
				TransferDomainOpen:           tt.fields.TransferDomainOpen,
				RenewOpenDomain:              tt.fields.RenewOpenDomain,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

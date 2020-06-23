package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
)

func TestFees_Validate(t *testing.T) {
	type fields struct {
		FeeCoinDenom string
		FeeCoinPrice types.Dec
		FeeDefault   types.Dec
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "success",
			fields: func() fields {
				fees := NewFeeConfiguration()
				fees.SetDefaults("test")
				return fields(*fees)
			}(),
			wantErr: false,
		},
		{
			name:    "fail missing fee",
			fields:  fields{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FeeConfiguration{
				FeeCoinDenom:                 tt.fields.FeeCoinDenom,
				FeeCoinPrice:                 tt.fields.FeeCoinPrice,
				FeeDefault:                   tt.fields.FeeDefault,
				RegisterAccountClosed:        tt.fields.RegisterAccountClosed,
				RegisterAccountOpen:          tt.fields.RegisterAccountOpen,
				TransferAccountClosed:        tt.fields.TransferAccountClosed,
				TransferAccountOpen:          tt.fields.TransferAccountOpen,
				ReplaceAccountTargets:        tt.fields.ReplaceAccountTargets,
				AddAccountCertificate:        tt.fields.AddAccountCertificate,
				DelAccountCertificate:        tt.fields.DelAccountCertificate,
				SetAccountMetadata:           tt.fields.SetAccountMetadata,
				RegisterDomain1:              tt.fields.RegisterDomain1,
				RegisterDomain2:              tt.fields.RegisterDomain2,
				RegisterDomain3:              tt.fields.RegisterDomain3,
				RegisterDomain4:              tt.fields.RegisterDomain4,
				RegisterDomain5:              tt.fields.RegisterDomain5,
				RegisterDomainDefault:        tt.fields.RegisterDomainDefault,
				RegisterOpenDomainMultiplier: tt.fields.RegisterOpenDomainMultiplier,
				TransferDomainClosed:         tt.fields.TransferDomainClosed,
				TransferDomainOpen:           tt.fields.TransferDomainOpen,
				RenewDomainOpen:              tt.fields.RenewDomainOpen,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

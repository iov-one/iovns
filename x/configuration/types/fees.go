package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/fatih/structs"
)

func NewFees() *Fees {
	return &Fees{}
}

// Validate validates the fee object
func (f *Fees) Validate() error {
	if f == nil {
		return fmt.Errorf("fees is nil")
	}
	m := structs.New(f)
	for _, field := range m.Fields() {
		switch fee := field.Value().(type) {
		case sdk.Dec:
			if fee.IsNil() {
				return fmt.Errorf("nil dec in field %s", field.Name())
			}
			if fee.IsZero() {
				return fmt.Errorf("zero dec in field %s", field.Name())
			}
			if fee.IsNegative() {
				return fmt.Errorf("negative dec in field %s", field.Name())
			}
		case string:
			if err := sdk.ValidateDenom(fee); err != nil {
				return fmt.Errorf("invalid coin denom in field %s: %s", field.Name(), fee)
			}
		default:
			panic(fmt.Sprintf("invalid type: %T", fee))
		}
	}
	return nil
}

// SetDefaults sets the default fees, it takes only one parameter which is the coin name that
// will be used by the users who want to access the domain module functionalities
func (f *Fees) SetDefaults(denom string) {
	if err := sdk.ValidateDenom(denom); err != nil {
		panic(fmt.Errorf("invalid coin denom %s: %w", denom, err))
	}
	defaultFeeParameter := sdk.NewDecFromInt(sdk.NewInt(10))
	if f == nil {
		panic("cannot set default fees for nil fees")
	}
	*f = Fees{
		FeeCoinDenom:                 denom,
		FeeCoinPrice:                 sdk.NewDecFromInt(sdk.NewInt(10)),
		FeeDefault:                   defaultFeeParameter,
		RegisterAccountClosed:        defaultFeeParameter,
		RegisterAccountOpen:          defaultFeeParameter,
		TransferAccountClosed:        defaultFeeParameter,
		TransferAccountOpen:          defaultFeeParameter,
		ReplaceAccountResources:      defaultFeeParameter,
		AddAccountCertificate:        defaultFeeParameter,
		DelAccountCertificate:        defaultFeeParameter,
		SetAccountMetadata:           defaultFeeParameter,
		RegisterDomain1:              defaultFeeParameter,
		RegisterDomain2:              defaultFeeParameter,
		RegisterDomain3:              defaultFeeParameter,
		RegisterDomain4:              defaultFeeParameter,
		RegisterDomain5:              defaultFeeParameter,
		RegisterDomainDefault:        defaultFeeParameter,
		RegisterOpenDomainMultiplier: sdk.NewDec(2),
		TransferDomainClosed:         defaultFeeParameter,
		TransferDomainOpen:           defaultFeeParameter,
		RenewDomainOpen:              defaultFeeParameter,
	}
}

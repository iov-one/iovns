package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/fatih/structs"
)

func NewFees() *Fees {
	return &Fees{}
}

// Fees contains different type of fees
// to calculate coins to detract when
// processing different messages
type Fees struct {
	// FeeCoinDenom defines the denominator of the coin used to process fees
	FeeCoinDenom string
	FeeCoinPrice sdk.Dec
	// DefaultFee is the parameter defining the default fee
	DefaultFee sdk.Dec
	// account fees
	RegisterClosedAccount sdk.Dec
	RegisterOpenAccount   sdk.Dec
	TransferClosedAccount sdk.Dec
	TransferOpenAccount   sdk.Dec
	ReplaceAccountTargets sdk.Dec
	AddAccountCertificate sdk.Dec
	DelAccountCertificate sdk.Dec
	SetAccountMetadata    sdk.Dec
	// domain fees
	// Register domain
	RegisterDomain               sdk.Dec
	RegisterDomain1              sdk.Dec
	RegisterDomain2              sdk.Dec
	RegisterDomain3              sdk.Dec
	RegisterDomain4              sdk.Dec
	RegisterDomain5              sdk.Dec
	RegisterDomainDefault        sdk.Dec
	RegisterOpenDomainMultiplier sdk.Dec
	// TransferDomain
	TransferDomainClosed sdk.Dec
	TransferDomainOpen   sdk.Dec
	// RenewDomain
	RenewOpenDomain sdk.Dec
}

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
	defaultFeeParameter := sdk.NewDec(1)
	if f == nil {
		panic("cannot set default fees for nil fees")
	}
	*f = Fees{
		FeeCoinDenom:                 denom,
		FeeCoinPrice:                 sdk.NewDec(10),
		DefaultFee:                   defaultFeeParameter,
		RegisterClosedAccount:        defaultFeeParameter,
		RegisterOpenAccount:          defaultFeeParameter,
		TransferClosedAccount:        defaultFeeParameter,
		TransferOpenAccount:          defaultFeeParameter,
		ReplaceAccountTargets:        defaultFeeParameter,
		AddAccountCertificate:        defaultFeeParameter,
		DelAccountCertificate:        defaultFeeParameter,
		SetAccountMetadata:           defaultFeeParameter,
		RegisterDomain:               defaultFeeParameter,
		RegisterDomain1:              defaultFeeParameter,
		RegisterDomain2:              defaultFeeParameter,
		RegisterDomain3:              defaultFeeParameter,
		RegisterDomain4:              defaultFeeParameter,
		RegisterDomain5:              defaultFeeParameter,
		RegisterDomainDefault:        defaultFeeParameter,
		RegisterOpenDomainMultiplier: sdk.NewDec(2),
		TransferDomainClosed:         defaultFeeParameter,
		TransferDomainOpen:           defaultFeeParameter,
		RenewOpenDomain:              defaultFeeParameter,
	}
}

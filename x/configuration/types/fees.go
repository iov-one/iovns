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
	FeeCoinDenom string  `json:"fee_coin_denom"`
	FeeCoinPrice sdk.Dec `json:"fee_coin_price"`
	// DefaultFee is the parameter defining the default fee
	DefaultFee sdk.Dec `json:"default_fee"`
	// account fees
	RegisterAccountClosed sdk.Dec `json:"register_account_closed"`
	RegisterAccountOpen   sdk.Dec `json:"register_account_open"`
	TransferAccountClosed sdk.Dec `json:"transfer_account_closed"`
	TransferAccountOpen   sdk.Dec `json:"transfer_account_open"`
	ReplaceAccountTargets sdk.Dec `json:"replace_account_targets"`
	AddAccountCertificate sdk.Dec `json:"add_account_certificate"`
	DelAccountCertificate sdk.Dec `json:"del_account_certificate"`
	SetAccountMetadata    sdk.Dec `json:"set_account_metadata"`
	// domain fees
	// Register domain
	RegisterDomain1              sdk.Dec `json:"register_domain_1"`
	RegisterDomain2              sdk.Dec `json:"register_domain_2"`
	RegisterDomain3              sdk.Dec `json:"register_domain_3"`
	RegisterDomain4              sdk.Dec `json:"register_domain_4"`
	RegisterDomain5              sdk.Dec `json:"register_domain_5"`
	RegisterDomainDefault        sdk.Dec `json:"register_domain_default"`
	RegisterOpenDomainMultiplier sdk.Dec `json:"register_open_domain_multiplier"`
	// TransferDomain
	TransferDomainClosed sdk.Dec `json:"transfer_domain_closed"`
	TransferDomainOpen   sdk.Dec `json:"transfer_domain_open"`
	// RenewDomain
	RenewDomainOpen sdk.Dec `json:"renew_domain_open"`
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
		RegisterAccountClosed:        defaultFeeParameter,
		RegisterAccountOpen:          defaultFeeParameter,
		TransferAccountClosed:        defaultFeeParameter,
		TransferAccountOpen:          defaultFeeParameter,
		ReplaceAccountTargets:        defaultFeeParameter,
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

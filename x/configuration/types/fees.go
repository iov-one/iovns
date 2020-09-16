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
	FeeCoinDenom string `json:"fee_coin_denom"`
	// FeeCoinPrice defines the price of the coin
	FeeCoinPrice sdk.Dec `json:"fee_coin_price"`
	// FeeDefault is the parameter defining the default fee
	FeeDefault sdk.Dec `json:"fee_default"`
	// account fees
	// RegisterAccountClosed is the fee to be paid to register an account in a closed domain
	RegisterAccountClosed sdk.Dec `json:"register_account_closed"`
	// RegisterAccountOpen is the fee to be paid to register an account in an open domain
	RegisterAccountOpen sdk.Dec `json:"register_account_open"`
	// TransferAccountClosed is the fee to be paid to register an account in a closed domain
	TransferAccountClosed sdk.Dec `json:"transfer_account_closed"`
	// TransferAccountOpen is the fee to be paid to register an account in an open domain
	TransferAccountOpen sdk.Dec `json:"transfer_account_open"`
	// ReplaceAccountResources is the fee to be paid to replace account's resources
	ReplaceAccountResources sdk.Dec `json:"replace_account_resources"`
	// AddAccountCertificate is the fee to be paid to add a certificate to an account
	AddAccountCertificate sdk.Dec `json:"add_account_certificate"`
	// DelAccountCertificate is the feed to be paid to delete a certificate in an account
	DelAccountCertificate sdk.Dec `json:"del_account_certificate"`
	// SetAccountMetadata is the fee to be paid to set account's metadata
	SetAccountMetadata sdk.Dec `json:"set_account_metadata"`
	// domain fees
	// Register domain
	// RegisterDomain1 is the fee to be paid to register a domain with one character
	RegisterDomain1 sdk.Dec `json:"register_domain_1"`
	// RegisterDomain2 is the fee to be paid to register a domain with two characters
	RegisterDomain2 sdk.Dec `json:"register_domain_2"`
	// RegisterDomain3 is the fee to be paid to register a domain with three characters
	RegisterDomain3 sdk.Dec `json:"register_domain_3"`
	// RegisterDomain4 is the fee to be paid to register a domain with four characters
	RegisterDomain4 sdk.Dec `json:"register_domain_4"`
	// RegisterDomain5 is the fee to be paid to register a domain with five characters
	RegisterDomain5 sdk.Dec `json:"register_domain_5"`
	// RegisterDomainDefault is the fee to be paid to register a domain with more than five characters
	RegisterDomainDefault sdk.Dec `json:"register_domain_default"`
	// RegisterDomainMultiplier is the multiplication applied to fees in register domain operations if they're of open type
	RegisterOpenDomainMultiplier sdk.Dec `json:"register_open_domain_multiplier"`
	// TransferDomain
	// TransferDomainClosed is the fee to be paid to transfer a closed domain
	TransferDomainClosed sdk.Dec `json:"transfer_domain_closed"`
	// TransferDomainOpen is the fee to be paid to transfer open domains
	TransferDomainOpen sdk.Dec `json:"transfer_domain_open"`
	// RenewDomainOpen is the fee to be paid to renew an open domain
	RenewDomainOpen sdk.Dec `json:"renew_domain_open"`
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

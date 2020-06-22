package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Params struct {
	// ID is the fee parameter identifier
	ID string
}

// ControllerFunc is the function signature used by account controllers
type FeeFunc func(ctx sdk.Context, msg ProductMsg) (sdk.Dec, error)

func NewFees() *Fees {
	return &Fees{}
}

// Fees contains different type of fees
// to calculate coins to detract when
// processing different messages
type FeeSeed sdk.Dec

func (f FeeSeed) Validate() error {
	if f.IsNil() {
		return fmt.Errorf("nil dec in id %s", f.ID)
	}
	if f.Amount.IsZero() {
		return fmt.Errorf("zero dec in id %s", f.ID)
	}
	if f.Amount.IsNegative() {
		return fmt.Errorf("negative dec in id %s", f.ID)
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
		FeeDefault:                   defaultFeeParameter,
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

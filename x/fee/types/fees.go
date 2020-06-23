package types

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/fee/contracts"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type FeeSeed struct {
	ID     string  `json:"id"`
	Amount sdk.Dec `json:"amount"`
}

func ValidateAmount(amount sdk.Dec) error {
	if amount.IsNil() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "nil dec")
	}
	if amount.IsZero() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "zero dec")
	}
	if amount.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "negative dec")
	}
	return nil
}

func (f FeeSeed) Validate() error {
	if len(f.ID) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "nil id")
	}
	if err := ValidateAmount(f.Amount); err != nil {
		return sdkerrors.Wrapf(err, "seed id %s", f.ID)
	}
	return nil
}

type FeeParamaters struct {
	FeeCoinDenom string  `json:"fee_coin_denomination"`
	FeeCoinPrice sdk.Dec `json:"fee_coin_price"`
	FeeDefault   sdk.Dec `json:"fee_default"`
}

func (fp FeeParamaters) Validate() error {
	if err := sdk.ValidateDenom(fp.FeeCoinDenom); err != nil {
		return err
	}
	if err := ValidateAmount(fp.FeeCoinPrice); err != nil {
		return sdkerrors.Wrapf(err, "param id %s", "fee_coin_price")
	}
	if err := ValidateAmount(fp.FeeDefault); err != nil {
		return sdkerrors.Wrapf(err, "param id %s", "fee_default")
	}
	return nil
}

// FeeConfiguration is the modules representation in the genesis
type FeeConfiguration struct {
	FeeConfigurer sdk.AccAddress `json:"fee_configurer"`
	FeeParamaters FeeParamaters  `json:"fee_parameters"`
	FeeSeeds      []FeeSeed      `json:"fee_seeds"`
}

func (f FeeConfiguration) Validate() error {
	if f.FeeConfigurer.Empty() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "no configurer specified")
	}
	if err := f.FeeParamaters.Validate(); err != nil {
		return err
	}
	for _, fee := range f.FeeSeeds {
		if err := fee.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func NewFeeConfiguration() *FeeConfiguration {
	return &FeeConfiguration{}
}

// SetDefaults sets the default fees, it takes only one parameter which is the coin name that
// will be used by the users who want to access the domain module functionalities
func (f *FeeConfiguration) SetDefaults(denom string) {
	if err := sdk.ValidateDenom(denom); err != nil {
		panic(fmt.Errorf("invalid coin denom %s: %w", denom, err))
	}
	if f == nil {
		panic("cannot set default fees for nil fees")
	}
	*f = FeeConfiguration{
		FeeParamaters: FeeParamaters{
			FeeCoinDenom: denom,
			FeeCoinPrice: sdk.NewDec(1),
			FeeDefault:   sdk.NewDec(1),
		},
		FeeSeeds: contracts.ContractFeeSeeds,
	}
}

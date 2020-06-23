package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgUpdateConfiguration struct {
	Fees       *FeeConfiguration
	Configurer sdk.AccAddress
}

func (m MsgUpdateConfiguration) Route() string {
	return RouterKey
}

func (m MsgUpdateConfiguration) Type() string {
	return "update_fees"
}

func (m MsgUpdateConfiguration) ValidateBasic() error {
	if m.Configurer.Empty() {
		return errors.Wrapf(errors.ErrInvalidRequest, "signer is missing")
	}
	// check if fees are valid
	if err := m.Fees.Validate(); err != nil {
		return errors.Wrap(errors.ErrInvalidRequest, err.Error())
	}
	return nil
}

func (m MsgUpdateConfiguration) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgUpdateConfiguration) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{m.Configurer} }

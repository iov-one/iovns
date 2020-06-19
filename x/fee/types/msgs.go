package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgUpdateFees struct {
	Fees       *Fees
	Configurer sdk.AccAddress
}

func (m MsgUpdateFees) Route() string {
	return RouterKey
}

func (m MsgUpdateFees) Type() string {
	return "update_fees"
}

func (m MsgUpdateFees) ValidateBasic() error {
	if m.Configurer.Empty() {
		return errors.Wrapf(errors.ErrInvalidRequest, "signer is missing")
	}
	// check if fees are valid
	if err := m.Fees.Validate(); err != nil {
		return errors.Wrap(errors.ErrInvalidRequest, err.Error())
	}
	return nil
}

func (m MsgUpdateFees) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgUpdateFees) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{m.Configurer} }

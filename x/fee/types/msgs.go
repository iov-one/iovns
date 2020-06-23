package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgUpdateFeeConfiguration struct {
	NewFeeConfiguration *FeeConfiguration `json:"new_fee_configuration"`
	Signer              sdk.AccAddress
}

func (m MsgUpdateFeeConfiguration) Route() string {
	return RouterKey
}

func (m MsgUpdateFeeConfiguration) Type() string {
	return "update_fees"
}

func (m MsgUpdateFeeConfiguration) ValidateBasic() error {
	if m.Signer.Empty() {
		return errors.Wrapf(errors.ErrInvalidRequest, "signer is missing")
	}
	// check if fees are valid
	if err := m.NewFeeConfiguration.Validate(); err != nil {
		return errors.Wrap(errors.ErrInvalidRequest, err.Error())
	}
	return nil
}

func (m MsgUpdateFeeConfiguration) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgUpdateFeeConfiguration) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{m.Signer} }

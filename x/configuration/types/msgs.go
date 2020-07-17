package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgUpdateConfig is used to update
// configuration using a multisig strategy
type MsgUpdateConfig struct {
	Signer           sdk.AccAddress
	NewConfiguration Config
}

var _ sdk.Msg = (*MsgUpdateConfig)(nil)

func (m MsgUpdateConfig) Route() string { return RouterKey }

func (m MsgUpdateConfig) Type() string { return "update_config" }

func (m MsgUpdateConfig) ValidateBasic() error {
	if m.Signer.Empty() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "no signer specified")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m MsgUpdateConfig) GetSignBytes() []byte { return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m)) }

// GetSigners implements sdk.Msg
func (m MsgUpdateConfig) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{m.Signer} }

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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "signer is missing")
	}
	// check if fees are valid
	if err := m.Fees.Validate(); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	return nil
}

func (m MsgUpdateFees) GetSignBytes() []byte { return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m)) }

func (m MsgUpdateFees) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{m.Configurer} }

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgUpdateConfig is used to update
// configuration using a multisig strategy
type MsgUpdateConfig struct {
	// Signer is the address of the entity who is doing the transaction
	Signer sdk.AccAddress
	// NewConfiguration contains the new configuration data
	NewConfiguration Config
}

var _ sdk.Msg = (*MsgUpdateConfig)(nil)

// Route implements sdk.Msg
func (m MsgUpdateConfig) Route() string { return RouterKey }

// Type implements sdk.Msg
func (m MsgUpdateConfig) Type() string { return "update_config" }

// ValidateBasic implements sdk.Msg
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
	// Fees represent the new fees to apply
	Fees *Fees
	// Configurer is the address that is singing the message
	Configurer sdk.AccAddress
}

// Route implements sdk.Msg
func (m MsgUpdateFees) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m MsgUpdateFees) Type() string {
	return "update_fees"
}

// ValidateBasic implements sdk.Msg
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

// GetSignBytes implements sdk.Msg
func (m MsgUpdateFees) GetSignBytes() []byte { return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m)) }

// GetSigners implements sdk.Msg
func (m MsgUpdateFees) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{m.Configurer} }

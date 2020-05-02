package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgUpdateConfig is used to update
// configuration using a multisig strategy
type MsgUpdateConfig struct {
	Signers          []sdk.AccAddress
	NewConfiguration Config
}

func (m MsgUpdateConfig) Route() string {
	return RouterKey
}

func (m MsgUpdateConfig) Type() string {
	return "update_config"
}

func (m MsgUpdateConfig) ValidateBasic() error {
	if m.Signers == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "no signers specified")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m MsgUpdateConfig) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m MsgUpdateConfig) GetSigners() []sdk.AccAddress {
	return m.Signers
}

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgUpdateConfig is used to update
// configuration using a multisig strategy
type MsgUpdateConfig struct {
	Configurer       sdk.AccAddress
	NewConfiguration Config
}

var _ sdk.Msg = (*MsgUpdateConfig)(nil)

func (m MsgUpdateConfig) Route() string { return RouterKey }

func (m MsgUpdateConfig) Type() string { return "update_config" }

func (m MsgUpdateConfig) ValidateBasic() error {
	if m.Configurer.Empty() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "no configurer specified")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m MsgUpdateConfig) GetSignBytes() []byte { return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m)) }

// GetSigners implements sdk.Msg
func (m MsgUpdateConfig) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{m.Configurer} }

// MsgUpdateDefaultFee inserts or update the default fees
type MsgUpsertDefaultFee struct {
	Configurer sdk.AccAddress
	Module     string
	MsgType    string
	Fee        sdk.Coin
}

var _ sdk.Msg = (*MsgUpsertDefaultFee)(nil)

func (m MsgUpsertDefaultFee) Route() string { return RouterKey }

func (m MsgUpsertDefaultFee) Type() string { return "upsert_default_fee" }

func (m MsgUpsertDefaultFee) ValidateBasic() error {
	if m.Configurer.Empty() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "missing configurer")
	}
	if m.Fee == (sdk.Coin{}) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "missing fee")
	}
	if m.Module == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "missing module")
	}
	if m.MsgType == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "missing msg type")
	}
	return nil
}

func (m MsgUpsertDefaultFee) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgUpsertDefaultFee) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{m.Configurer} }

// MsgUpsertLevelFee inserts or update a level fee
type MsgUpsertLevelFee struct {
	Configurer sdk.AccAddress
	Module     string
	MsgType    string
	Level      sdk.Int
	Fee        sdk.Coin
}

var _ sdk.Msg = (*MsgUpsertLevelFee)(nil)

func (m MsgUpsertLevelFee) Route() string {
	return RouterKey
}

func (m MsgUpsertLevelFee) Type() string {
	return "upsert_level_fee"
}

func (m MsgUpsertLevelFee) ValidateBasic() error {
	if m.Configurer.Empty() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "missing configurer")
	}
	if m.Fee == (sdk.Coin{}) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "missing fee")
	}
	if m.Module == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "missing module")
	}
	if m.MsgType == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "missing msg type")
	}
	return nil
}

func (m MsgUpsertLevelFee) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgUpsertLevelFee) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{m.Configurer} }

// MsgDeleteLevelFee deletes a level fee
type MsgDeleteLevelFee struct {
	Configurer sdk.AccAddress
	Module     string
	MsgType    string
	Level      sdk.Int
}

var _ sdk.Msg = (*MsgDeleteLevelFee)(nil)

func (m MsgDeleteLevelFee) Route() string { return RouterKey }

func (m MsgDeleteLevelFee) Type() string { return "delete_level_fee" }

func (m MsgDeleteLevelFee) ValidateBasic() error {
	if m.Configurer.Empty() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "missing configurer")
	}
	if m.Module == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "missing module")
	}
	if m.MsgType == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "missing msg type")
	}
	return nil
}

func (m MsgDeleteLevelFee) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgDeleteLevelFee) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{m.Configurer} }

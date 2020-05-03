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

func (m MsgUpdateConfig) Route() string { return RouterKey }

func (m MsgUpdateConfig) Type() string { return "update_config" }

func (m MsgUpdateConfig) ValidateBasic() error {
	if m.Signers == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "no signers specified")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m MsgUpdateConfig) GetSignBytes() []byte { return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m)) }

// GetSigners implements sdk.Msg
func (m MsgUpdateConfig) GetSigners() []sdk.AccAddress { return m.Signers }

// MsgUpdateDefaultFee inserts or update the default fees
type MsgUpsertDefaultFee struct {
	Signers []sdk.AccAddress
	Module  string
	MsgType string
	Fee     sdk.Coin
}

func (m MsgUpsertDefaultFee) Route() string { return RouterKey }

func (m MsgUpsertDefaultFee) Type() string { return "upsert_default_fee" }

func (m MsgUpsertDefaultFee) ValidateBasic() error {
	if len(m.Signers) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "missing signers")
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

func (m MsgUpsertDefaultFee) GetSigners() []sdk.AccAddress { return m.Signers }

// MsgUpsertLevelFee inserts or update a level fee
type MsgUpsertLevelFee struct {
	Signers []sdk.AccAddress
	Module  string
	MsgType string
	Level   sdk.Int
	Fee     sdk.Coin
}

func (m MsgUpsertLevelFee) Route() string {
	return RouterKey
}

func (m MsgUpsertLevelFee) Type() string {
	return "upsert_level_fee"
}

func (m MsgUpsertLevelFee) ValidateBasic() error {
	if len(m.Signers) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "missing signers")
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

func (m MsgUpsertLevelFee) GetSigners() []sdk.AccAddress { return m.Signers }

// MsgDeleteLevelFee deletes a level fee
type MsgDeleteLevelFee struct {
	Signers []sdk.AccAddress
	Module  string
	MsgType string
	Level   sdk.Int
}

func (m MsgDeleteLevelFee) Route() string { return RouterKey }

func (m MsgDeleteLevelFee) Type() string { return "delete_level_fee" }

func (m MsgDeleteLevelFee) ValidateBasic() error {
	if len(m.Signers) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "missing signers")
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

func (m MsgDeleteLevelFee) GetSigners() []sdk.AccAddress { return m.Signers }

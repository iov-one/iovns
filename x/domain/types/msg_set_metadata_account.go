package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgSetAccountMetadata is the function used
// to set accounts metadata
type MsgSetAccountMetadata struct {
	// Domain is the domain name of the account
	Domain string
	// Name is the name of the account
	Name string
	// NewMetadataURI is the metadata URI of the account
	// we want to update or insert
	NewMetadataURI string
	// Signer is the owner of the account
	Signer sdk.AccAddress
}

// Route implements sdk.Msg
func (m MsgSetAccountMetadata) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m MsgSetAccountMetadata) Type() string {
	return "set_account_metadata"
}

// ValidateBasic implements sdk.Msg
func (m MsgSetAccountMetadata) ValidateBasic() error {
	if m.Domain == "" {
		return sdkerrors.Wrapf(ErrInvalidDomainName, "empty")
	}
	if m.Name == "" {
		return sdkerrors.Wrapf(ErrInvalidAccountName, "empty")
	}
	if m.NewMetadataURI == "" {
		return sdkerrors.Wrapf(ErrInvalidRequest, "metadata uri is empty")
	}
	if m.Signer.Empty() {
		return sdkerrors.Wrap(ErrInvalidOwner, "empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m MsgSetAccountMetadata) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m MsgSetAccountMetadata) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Signer}
}

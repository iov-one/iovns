package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgAddAccountCertificates is the message used
// when a user wants to add new certificates
// to his account
type MsgAddAccountCertificates struct {
	// Domain is the domain of the account
	Domain string
	// Name is the name of the account
	Name string
	// Owner is the owner of the account
	Owner sdk.AccAddress
	// NewCertificate is the new certificate to add
	NewCertificate []byte
}

// Route implements sdk.Msg
func (m *MsgAddAccountCertificates) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgAddAccountCertificates) Type() string {
	return "add_certificates_account"
}

// ValidateBasic implements sdk.Msg
func (m *MsgAddAccountCertificates) ValidateBasic() error {
	if m.Domain == "" {
		return sdkerrors.Wrapf(ErrInvalidDomainName, "empty")
	}
	if m.Name == "" {
		return sdkerrors.Wrap(ErrInvalidAccountName, "empty")
	}
	if m.Owner == nil {
		return sdkerrors.Wrap(ErrInvalidOwner, "empty")
	}
	if m.NewCertificate == nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "certificate is empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgAddAccountCertificates) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgAddAccountCertificates) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

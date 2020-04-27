package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgDeleteAccountCertificates is the request
// model used to remove certificates from an
// account
type MsgDeleteAccountCertificate struct {
	// Domain is the name of the domain of the account
	Domain string
	// Name is the name of the account
	Name string
	// DeleteCertificate is the certificate to delete
	DeleteCertificate []byte
	// Owner is the owner of the account
	Owner sdk.AccAddress
}

// Route implements sdk.Msg
func (m *MsgDeleteAccountCertificate) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgDeleteAccountCertificate) Type() string {
	return "delete_certificate_account"
}

// ValidateBasic implements sdk.Msg
func (m *MsgDeleteAccountCertificate) ValidateBasic() error {
	if m.Domain == "" {
		return sdkerrors.Wrapf(ErrInvalidDomainName, "empty")
	}
	if m.Name == "" {
		return sdkerrors.Wrap(ErrInvalidAccountName, "empty")
	}
	if m.Owner == nil {
		return sdkerrors.Wrap(ErrInvalidOwner, "empty")
	}
	if m.DeleteCertificate == nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "certificate is empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgDeleteAccountCertificate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgDeleteAccountCertificate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

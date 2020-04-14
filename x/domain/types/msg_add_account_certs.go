package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgAddAccountCertificates struct {
	Domain         string
	Name           string
	Owner          sdk.AccAddress
	NewCertificate []byte
}

func (m MsgAddAccountCertificates) Route() string {
	return RouterKey
}

func (m MsgAddAccountCertificates) Type() string {
	return "add_certificates_account"
}

func (m MsgAddAccountCertificates) ValidateBasic() error {
	if m.Domain == "" {
		return sdkerrors.Wrapf(ErrInvalidDomainName, "domain name is empty")
	}
	if m.Name == "" {
		return sdkerrors.Wrap(ErrInvalidAccountName, "account name is empty")
	}
	if m.Owner == nil {
		return sdkerrors.Wrap(ErrInvalidOwner, "owner name is empty")
	}
	if m.NewCertificate == nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "certificate is empty")
	}
	return nil
}

func (m MsgAddAccountCertificates) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgAddAccountCertificates) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

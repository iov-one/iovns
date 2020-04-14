package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgDeleteAccountCertificate struct {
	Domain            string
	Name              string
	DeleteCertificate []byte
	Owner             sdk.AccAddress
}

func (m MsgDeleteAccountCertificate) Route() string {
	return RouterKey
}

func (m MsgDeleteAccountCertificate) Type() string {
	return "delete_certificate_account"
}

func (m MsgDeleteAccountCertificate) ValidateBasic() error {
	if m.Domain == "" {
		return sdkerrors.Wrapf(ErrInvalidDomainName, "domain name is empty")
	}
	if m.Name == "" {
		return sdkerrors.Wrap(ErrInvalidAccountName, "account name is empty")
	}
	if m.Owner == nil {
		return sdkerrors.Wrap(ErrInvalidOwner, "owner name is empty")
	}
	if m.DeleteCertificate == nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "certificate is empty")
	}
	return nil
}

func (m MsgDeleteAccountCertificate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgDeleteAccountCertificate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

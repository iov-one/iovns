package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgRenewDomain struct {
	// Domain is the domain name to renew
	Domain string
}

func (m *MsgRenewDomain) Route() string {
	return RouterKey
}

func (m *MsgRenewDomain) Type() string {
	return "renew_domain"
}

func (m *MsgRenewDomain) ValidateBasic() error {
	if m.Domain == "" {
		return sdkerrors.Wrapf(ErrInvalidDomainName, "domain name is empty")
	}
	return nil
}

func (m *MsgRenewDomain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgRenewDomain) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

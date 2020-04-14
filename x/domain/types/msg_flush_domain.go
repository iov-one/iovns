package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgFlushDomain is used to flush a domain
type MsgFlushDomain struct {
	Domain string
	Owner  sdk.AccAddress
}

func (m MsgFlushDomain) Route() string {
	return RouterKey
}

func (m MsgFlushDomain) Type() string {
	return "delete_domain"
}

func (m MsgFlushDomain) ValidateBasic() error {
	if m.Domain == "" {
		return sdkerrors.Wrap(ErrInvalidDomainName, "domain name is empty")
	}
	if m.Owner == nil {
		return sdkerrors.Wrap(ErrInvalidOwner, "owner is empty")
	}
	return nil
}

func (m MsgFlushDomain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgFlushDomain) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

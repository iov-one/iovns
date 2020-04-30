package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgFlushDomain is used to flush a domain
type MsgFlushDomain struct {
	// Domain is the domain name to flush
	Domain string
	// Owner is the owner of the domain
	Owner sdk.AccAddress
}

// Route implements sdk.Msg
func (m *MsgFlushDomain) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgFlushDomain) Type() string {
	return "delete_domain"
}

// ValidateBasic implements sdk.Msg
func (m *MsgFlushDomain) ValidateBasic() error {
	if m.Domain == "" {
		return sdkerrors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Owner == nil {
		return sdkerrors.Wrap(ErrInvalidOwner, "empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgFlushDomain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgFlushDomain) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

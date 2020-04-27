package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgDeleteDomain struct {
	Domain string
	Owner  sdk.AccAddress
}

func (m *MsgDeleteDomain) Route() string {
	return RouterKey
}

func (m *MsgDeleteDomain) Type() string {
	return "delete_domain"
}

func (m *MsgDeleteDomain) ValidateBasic() error {
	if m.Domain == "" {
		return sdkerrors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Owner == nil {
		return sdkerrors.Wrap(ErrInvalidOwner, "empty")
	}
	// success
	return nil
}

func (m *MsgDeleteDomain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgDeleteDomain) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

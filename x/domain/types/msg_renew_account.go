package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgRenewAccount struct {
	Domain string
	Name   string
}

func (m *MsgRenewAccount) Route() string {
	return RouterKey
}

func (m *MsgRenewAccount) Type() string {
	return "renew_account"
}

func (m *MsgRenewAccount) ValidateBasic() error {
	if m.Domain == "" {
		return sdkerrors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Name == "" {
		return sdkerrors.Wrap(ErrInvalidAccountName, "empty")
	}
	return nil
}

func (m *MsgRenewAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgRenewAccount) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

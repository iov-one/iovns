package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgTransferDomain struct {
	Domain   string
	Owner    sdk.AccAddress
	NewAdmin sdk.AccAddress
}

func (m MsgTransferDomain) Route() string {
	return RouterKey
}

func (m MsgTransferDomain) Type() string {
	return "transfer_domain"
}

func (m MsgTransferDomain) ValidateBasic() error {
	if m.Domain == "" {
		return sdkerrors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Owner == nil {
		return sdkerrors.Wrap(ErrInvalidOwner, "empty")
	}
	if m.NewAdmin == nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "new admin is empty")
	}
	return nil
}

func (m MsgTransferDomain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgTransferDomain) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

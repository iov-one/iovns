package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
)

type MsgReplaceAccountTargets struct {
	Domain     string
	Name       string
	NewTargets []iovns.BlockchainAddress
	Owner      sdk.AccAddress
}

func (m *MsgReplaceAccountTargets) Route() string {
	return RouterKey
}

func (m *MsgReplaceAccountTargets) Type() string {
	return "replace_account_targets"
}

func (m *MsgReplaceAccountTargets) ValidateBasic() error {
	if m.Domain == "" {
		return sdkerrors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Name == "" {
		return sdkerrors.Wrap(ErrInvalidAccountName, "empty")
	}
	if m.Owner == nil {
		return sdkerrors.Wrap(ErrInvalidOwner, "empty")
	}
	if len(m.NewTargets) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "empty blockchain targets")
	}
	return nil
}

func (m *MsgReplaceAccountTargets) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgReplaceAccountTargets) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

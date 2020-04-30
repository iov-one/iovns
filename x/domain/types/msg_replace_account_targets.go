package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
)

// MsgReplaceAccountTargets is the request model
// used to renew blockchain addresses associated
// with an account
type MsgReplaceAccountTargets struct {
	// Domain is the domain name of the account
	Domain string
	// Name is the name of the account
	Name string
	// NewTargets are the new blockchain addresses
	NewTargets []iovns.BlockchainAddress
	// Owner is the owner of the account
	Owner sdk.AccAddress
}

// Route implements sdk.Msg
func (m *MsgReplaceAccountTargets) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgReplaceAccountTargets) Type() string {
	return "replace_account_targets"
}

// ValidateBasic implements sdk.Msg
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

// GetSignBytes implements sdk.Msg
func (m *MsgReplaceAccountTargets) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgReplaceAccountTargets) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

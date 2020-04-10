package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgTransferAccount is the message used to transfer accounts
type MsgTransferAccount struct {
	Domain   string
	Name     string
	Owner    sdk.AccAddress
	NewOwner sdk.AccAddress
}

// Route returns the name of the module
func (m MsgTransferAccount) Route() string {
	return RouterKey
}

// Type returns the action
func (m MsgTransferAccount) Type() string {
	return "transfer_account"
}

// ValidateBasic does stateless checks on the request
// it checks if the domain name is valid
func (m MsgTransferAccount) ValidateBasic() error {
	if m.Domain == "" {
		return sdkerrors.Wrap(ErrInvalidDomainName, "domain name is missing")
	}
	if m.Name == "" {
		return sdkerrors.Wrap(ErrInvalidAccountName, "account name is empty")
	}
	if m.Owner == nil {
		return sdkerrors.Wrap(ErrInvalidOwner, "owner is empty")
	}
	if m.NewOwner == nil {
		return sdkerrors.Wrap(ErrInvalidOwner, "new owner is empty")
	}
	// success
	return nil
}

// GetSignBytes returns an ordered json of the request
func (m MsgTransferAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners returns the list of address that should list the request
// in this case the admin of the domain
func (m MsgTransferAccount) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

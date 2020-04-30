package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
)

// MsgRegisterAccount is the request
// model used to register new accounts
type MsgRegisterAccount struct {
	// Domain is the domain of the account
	Domain string
	// Name is the name of the account
	Name string
	// Owner is the owner of the account
	Owner sdk.AccAddress
	// Targets are the blockchain addresses of the account
	Targets []iovns.BlockchainAddress
	// Broker is the account that facilitated the transaction
	Broker sdk.AccAddress
}

// Route implements sdk.Msg
func (m *MsgRegisterAccount) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgRegisterAccount) Type() string {
	return "register_account"
}

// ValidateBasic implements sdk.Msg
func (m *MsgRegisterAccount) ValidateBasic() error {
	if m.Domain == "" {
		return sdkerrors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Owner.Empty() {
		return sdkerrors.Wrap(ErrInvalidOwner, "empty")
	}
	if m.Name == "" {
		return sdkerrors.Wrap(ErrInvalidAccountName, "empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgRegisterAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgRegisterAccount) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

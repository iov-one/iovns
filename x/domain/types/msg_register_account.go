package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
)

type MsgRegisterAccount struct {
	Domain  string
	Name    string
	Owner   sdk.AccAddress
	Targets []iovns.BlockchainAddress
	Broker  sdk.AccAddress
}

// Route returns the route key for the request
func (m MsgRegisterAccount) Route() string {
	return RouterKey
}

// Type returns the type of the msg
func (m MsgRegisterAccount) Type() string {
	return "register_account"
}

// ValidateBasic checks the request in a stateless way
func (m MsgRegisterAccount) ValidateBasic() error {
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

// GetSignBytes returns the expected signature
func (m MsgRegisterAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners returns the expected signers of the request
func (m MsgRegisterAccount) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

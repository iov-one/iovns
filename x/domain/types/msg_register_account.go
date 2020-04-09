package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovnsd"
)

type MsgRegisterAccount struct {
	Domain  string
	Name    string
	Owner   sdk.AccAddress
	Targets []iovnsd.BlockchainAddress
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
		return sdkerrors.Wrap(ErrInvalidDomain, "empty domain name")
	}
	if m.Owner.Empty() {
		return sdkerrors.Wrap(ErrInvalidOwner, "owner empty")
	}
	if m.Name == "" {
		return sdkerrors.Wrap(ErrInvalidName, "account name empty")
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

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgRegisterDomain is the request used to register new domains
type MsgRegisterDomain struct {
	// Name is the name of the domain we want to register
	Name string `json:"domain" arg:"--domain" helper:"name of the domain"`
	// Admin is the address of the newly registered domain
	Admin sdk.AccAddress `json:"admin"`
	// HasSuperuser defines if the domain registered has an owner or not
	HasSuperuser bool `json:"has_superuser"`
	// Broker TODO document
	Broker sdk.AccAddress `json:"broker" arg:"--broker" helper:"the broker"`
	// AccountRenew defines the expiration time in seconds of each newly registered account.
	AccountRenew int64 `json:"account_renew" arg:"--account-renew" helper:"account's renewal time in seconds"`
	// TODO MSGFEEs
}

// Route implements sdk.Msg
func (m *MsgRegisterDomain) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgRegisterDomain) Type() string {
	return "register_domain"
}

// ValidateBasic implements sdk.Msg
func (m *MsgRegisterDomain) ValidateBasic() error {
	if m.Admin == nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "admin is missing")
	}
	if m.AccountRenew == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "account renew value can not be zero")
	}
	if m.Name == "" {
		return sdkerrors.Wrap(ErrInvalidDomainName, "empty")
	}
	// success
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgRegisterDomain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgRegisterDomain) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Admin}
}

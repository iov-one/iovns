package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgRegisterDomain is the request used to register new domains
type MsgRegisterDomain struct {
	// Name is the name of the domain we want to register
	Name string `json:"domain"`
	// Admin is the address of the newly registered domain
	Admin sdk.AccAddress `json:"admin"`
	// HasSuperuser defines if the domain registered has an owner or not
	HasSuperuser bool `json:"has_superuser"`
	// Broker TODO document
	Broker sdk.AccAddress `json:"broker"`
	// AccountRenew defines the expiration time in seconds of each newly registered account.
	AccountRenew int64
	// TODO MSGFEEs
}

// Route returns the name of the module
func (m MsgRegisterDomain) Route() string {
	return RouterKey
}

// Type returns the action
func (m MsgRegisterDomain) Type() string {
	return "register_domain"
}

// ValidateBasic does stateless checks on the request
// it checks if the domain name is valid
func (m MsgRegisterDomain) ValidateBasic() error {
	if m.Admin == nil {
		return sdkerrors.Wrap(ErrInvalidRegisterDomainRequest, "admin is missing")
	}
	if m.AccountRenew == 0 {
		return sdkerrors.Wrap(ErrInvalidRegisterDomainRequest, "account renew value can not be zero")
	}
	if m.Name == "" {
		return sdkerrors.Wrap(ErrInvalidRegisterDomainRequest, "domain name is missing")
	}
	// success
	return nil
}

// GetSignBytes returns an ordered json of the request
func (m MsgRegisterDomain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners returns the list of address that should list the request
// in this case the admin of the domain
func (m MsgRegisterDomain) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Admin}
}

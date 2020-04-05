package types

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: Describe your actions, these will implment the interface of `sdk.Msg`

var _ sdk.Msg = MsgRegisterDomain{}

// MsgRegisterDomain is used to register new domains
type MsgRegisterDomain struct {
	// Domain is the name of the domain to register
	Domain string `json:"domain"`
	// Admin is the administrator of the domain TODO can someone buy a domain for someone else?
	Admin sdk.AccAddress `json:"admin"`
	// HasSuperuser TODO explain better what this does
	HasSuperuser bool `json:"has_superuser"`
	// Broker is the broker that facilitated the request
	Broker sdk.AccAddress `json:"acc_address"`
	// AccountRenew TODO explain
	AccountRenew int64 `json:"account_renew"`
}

// ValidateBasic is used to validate the request without checking the state
func (msg MsgRegisterDomain) ValidateBasic() (err error) {
	// validate domain
	if err = validateDomain(msg.Domain); err != nil {
		return sdkerrors.Wrap(ErrInvalidDomain, err.Error())
	}
	// check owner
	if msg.Admin.Empty() {
		return sdkerrors.Wrap(ErrInvalidOwner, "empty")
	}
	// nothing else to validate in a stateless way
	return
}

// validateDomain checks if a domain formation is valid
// TODO check extra constraints
func validateDomain(domain string) error {
	switch domain {
	// check empty domain
	case "":
		return errors.New("empty string")
	default:
		return nil
	}
}

// implement sdk.Msg

// nolint
// Route returns the name of the module
func (msg MsgRegisterDomain) Route() string { return RouterKey }

// Type returns the action
func (msg MsgRegisterDomain) Type() string { return "register_domain" }

// GetSignBytes encodes the message for signing
func (msg MsgRegisterDomain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners returns who is supposed to sign the transaction
func (msg MsgRegisterDomain) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Admin}
}

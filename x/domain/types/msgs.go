package types

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
)

// MsgAddAccountCertificates is the message used
// when a user wants to add new certificates
// to his account
type MsgAddAccountCertificates struct {
	// Domain is the domain of the account
	Domain string
	// Name is the name of the account
	Name string
	// Owner is the owner of the account
	Owner types.AccAddress
	// NewCertificate is the new certificate to add
	NewCertificate []byte
}

// Route implements sdk.Msg
func (m *MsgAddAccountCertificates) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgAddAccountCertificates) Type() string {
	return "add_certificates_account"
}

// ValidateBasic implements sdk.Msg
func (m *MsgAddAccountCertificates) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrapf(ErrInvalidDomainName, "empty")
	}
	if m.Name == "" {
		return errors.Wrap(ErrInvalidAccountName, "empty")
	}
	if m.Owner == nil {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	if m.NewCertificate == nil {
		return errors.Wrap(ErrInvalidRequest, "certificate is empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgAddAccountCertificates) GetSignBytes() []byte {
	return types.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgAddAccountCertificates) GetSigners() []types.AccAddress {
	return []types.AccAddress{m.Owner}
}

// MsgDeleteAccountCertificate is the request
// model used to remove certificates from an
// account
type MsgDeleteAccountCertificate struct {
	// Domain is the name of the domain of the account
	Domain string
	// Name is the name of the account
	Name string
	// DeleteCertificate is the certificate to delete
	DeleteCertificate []byte
	// Owner is the owner of the account
	Owner types.AccAddress
}

// Route implements sdk.Msg
func (m *MsgDeleteAccountCertificate) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgDeleteAccountCertificate) Type() string {
	return "delete_certificate_account"
}

// ValidateBasic implements sdk.Msg
func (m *MsgDeleteAccountCertificate) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrapf(ErrInvalidDomainName, "empty")
	}
	if m.Name == "" {
		return errors.Wrap(ErrInvalidAccountName, "empty")
	}
	if m.Owner == nil {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	if m.DeleteCertificate == nil {
		return errors.Wrap(ErrInvalidRequest, "certificate is empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgDeleteAccountCertificate) GetSignBytes() []byte {
	return types.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgDeleteAccountCertificate) GetSigners() []types.AccAddress {
	return []types.AccAddress{m.Owner}
}

// MsgDeleteAccount is the request model
// used to delete an account
type MsgDeleteAccount struct {
	// Domain is the name of the domain of the account
	Domain string
	// Name is the name of the account
	Name string
	// Owner is the owner of the account
	Owner types.AccAddress
}

// Route implements sdk.Msg
func (m *MsgDeleteAccount) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgDeleteAccount) Type() string {
	return "delete_account"
}

// ValidateBasic implements sdk.Msg
func (m *MsgDeleteAccount) ValidateBasic() error {
	if m.Owner == nil {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	if m.Name == "" {
		return errors.Wrap(ErrInvalidAccountName, "empty")
	}
	if m.Domain == "" {
		return errors.Wrap(ErrInvalidDomainName, "empty")
	}
	// success
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgDeleteAccount) GetSignBytes() []byte {
	return types.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgDeleteAccount) GetSigners() []types.AccAddress {
	return []types.AccAddress{m.Owner}
}

// MsgDeleteDomain is the request
// model to delete a domain
type MsgDeleteDomain struct {
	Domain string
	Owner  types.AccAddress
}

// Route implements sdk.Msg
func (m *MsgDeleteDomain) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgDeleteDomain) Type() string {
	return "delete_domain"
}

// ValidateBasic implements sdk.Msg
func (m *MsgDeleteDomain) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Owner == nil {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	// success
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgDeleteDomain) GetSignBytes() []byte {
	return types.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgDeleteDomain) GetSigners() []types.AccAddress {
	return []types.AccAddress{m.Owner}
}

// MsgFlushDomain is used to flush a domain
type MsgFlushDomain struct {
	// Domain is the domain name to flush
	Domain string
	// Owner is the owner of the domain
	Owner types.AccAddress
}

// Route implements sdk.Msg
func (m *MsgFlushDomain) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgFlushDomain) Type() string {
	return "delete_domain"
}

// ValidateBasic implements sdk.Msg
func (m *MsgFlushDomain) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Owner == nil {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgFlushDomain) GetSignBytes() []byte {
	return types.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgFlushDomain) GetSigners() []types.AccAddress {
	return []types.AccAddress{m.Owner}
}

// MsgRegisterAccount is the request
// model used to register new accounts
type MsgRegisterAccount struct {
	// Domain is the domain of the account
	Domain string
	// Name is the name of the account
	Name string
	// Owner is the owner of the account
	Owner types.AccAddress
	// Targets are the blockchain addresses of the account
	Targets []iovns.BlockchainAddress
	// Broker is the account that facilitated the transaction
	Broker types.AccAddress
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
		return errors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Owner.Empty() {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	if m.Name == "" {
		return errors.Wrap(ErrInvalidAccountName, "empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgRegisterAccount) GetSignBytes() []byte {
	return types.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgRegisterAccount) GetSigners() []types.AccAddress {
	return []types.AccAddress{m.Owner}
}

// MsgRegisterDomain is the request used to register new domains
type MsgRegisterDomain struct {
	// Name is the name of the domain we want to register
	Name string `json:"domain" arg:"--domain" helper:"name of the domain"`
	// Admin is the address of the newly registered domain
	Admin types.AccAddress `json:"admin"`
	// HasSuperuser defines if the domain registered has an owner or not
	HasSuperuser bool `json:"has_superuser"`
	// Broker TODO document
	Broker types.AccAddress `json:"broker" arg:"--broker" helper:"the broker"`
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
		return errors.Wrap(ErrInvalidRequest, "admin is missing")
	}
	if m.AccountRenew == 0 {
		return errors.Wrap(ErrInvalidRequest, "account renew value can not be zero")
	}
	if m.Name == "" {
		return errors.Wrap(ErrInvalidDomainName, "empty")
	}
	// success
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgRegisterDomain) GetSignBytes() []byte {
	return types.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgRegisterDomain) GetSigners() []types.AccAddress {
	return []types.AccAddress{m.Admin}
}

// MsgRenewAccount is the request
// model used to renew accounts
type MsgRenewAccount struct {
	// Domain is the domain of the account
	Domain string
	// Name is the name of the account
	Name string
}

// Route implements sdk.Msg
func (m *MsgRenewAccount) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgRenewAccount) Type() string {
	return "renew_account"
}

// ValidateBasic implements sdk.Msg
func (m *MsgRenewAccount) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Name == "" {
		return errors.Wrap(ErrInvalidAccountName, "empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgRenewAccount) GetSignBytes() []byte {
	return types.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgRenewAccount) GetSigners() []types.AccAddress {
	return []types.AccAddress{}
}

// MsgRenewDomain is the request
// model used to renew a domain
type MsgRenewDomain struct {
	// Domain is the domain name to renew
	Domain string
}

// Route implements sdk.Msg
func (m *MsgRenewDomain) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgRenewDomain) Type() string {
	return "renew_domain"
}

// ValidateBasic implements sdk.Msg
func (m *MsgRenewDomain) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrapf(ErrInvalidDomainName, "empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgRenewDomain) GetSignBytes() []byte {
	return types.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgRenewDomain) GetSigners() []types.AccAddress {
	return []types.AccAddress{}
}

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
	Owner types.AccAddress
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
		return errors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Name == "" {
		return errors.Wrap(ErrInvalidAccountName, "empty")
	}
	if m.Owner == nil {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	if len(m.NewTargets) == 0 {
		return errors.Wrap(ErrInvalidRequest, "empty blockchain targets")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgReplaceAccountTargets) GetSignBytes() []byte {
	return types.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgReplaceAccountTargets) GetSigners() []types.AccAddress {
	return []types.AccAddress{m.Owner}
}

// MsgSetAccountMetadata is the function used
// to set accounts metadata
type MsgSetAccountMetadata struct {
	// Domain is the domain name of the account
	Domain string
	// Name is the name of the account
	Name string
	// NewMetadataURI is the metadata URI of the account
	// we want to update or insert
	NewMetadataURI string
	// Owner is the owner of the account
	Owner types.AccAddress
}

// Route implements sdk.Msg
func (m *MsgSetAccountMetadata) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgSetAccountMetadata) Type() string {
	return "set_account_metadata"
}

// ValidateBasic implements sdk.Msg
func (m *MsgSetAccountMetadata) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrapf(ErrInvalidDomainName, "empty")
	}
	if m.Name == "" {
		return errors.Wrapf(ErrInvalidAccountName, "empty")
	}
	if m.NewMetadataURI == "" {
		return errors.Wrapf(ErrInvalidRequest, "metadata uri is empty")
	}
	if m.Owner.Empty() {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgSetAccountMetadata) GetSignBytes() []byte {
	return types.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgSetAccountMetadata) GetSigners() []types.AccAddress {
	return []types.AccAddress{m.Owner}
}

// MsgTransferAccount is the request
// model used to transfer accounts
type MsgTransferAccount struct {
	// Domain is the domain name of the account
	Domain string
	// Account is the account name
	Name string
	// Owner is the actual owner of the account
	Owner types.AccAddress
	// NewOwner is the new owner of the account
	NewOwner types.AccAddress
}

// Route implements sdk.Msg
func (m *MsgTransferAccount) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgTransferAccount) Type() string {
	return "transfer_account"
}

// ValidateBasic implements sdk.Msg
func (m *MsgTransferAccount) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Name == "" {
		return errors.Wrap(ErrInvalidAccountName, "empty")
	}
	if m.Owner == nil {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	if m.NewOwner == nil {
		return errors.Wrap(ErrInvalidOwner, "new owner is empty")
	}
	// success
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgTransferAccount) GetSignBytes() []byte {
	return types.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgTransferAccount) GetSigners() []types.AccAddress {
	return []types.AccAddress{m.Owner}
}

// MsgTransferDomain is the request model
// used to transfer a domain
type MsgTransferDomain struct {
	Domain   string
	Owner    types.AccAddress
	NewAdmin types.AccAddress
}

// Route implements sdk.Msg
func (m *MsgTransferDomain) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgTransferDomain) Type() string {
	return "transfer_domain"
}

// ValidateBasic implements sdk.Msg
func (m *MsgTransferDomain) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Owner == nil {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	if m.NewAdmin == nil {
		return errors.Wrap(ErrInvalidRequest, "new admin is empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgTransferDomain) GetSignBytes() []byte {
	return types.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgTransferDomain) GetSigners() []types.AccAddress {
	return []types.AccAddress{m.Owner}
}

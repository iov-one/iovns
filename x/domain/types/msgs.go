package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	fee "github.com/iov-one/iovns/x/fee/types"
)

// MsgAddAccountCertificates is the message used
// when a user wants to add new certificates
// to his account
type MsgAddAccountCertificates struct {
	// Domain is the domain of the account
	Domain string `json:"domain"`
	// Name is the name of the account
	Name string `json:"name"`
	// Owner is the owner of the account
	Owner sdk.AccAddress `json:"owner"`
	// NewCertificate is the new certificate to add
	NewCertificate []byte         `json:"new_certificate"`
	FeePayerAddr   sdk.AccAddress `json:"fee_payer"`
}

var _ fee.ProductMsg = (*MsgAddAccountCertificates)(nil)

// Route implements sdk.Msg
func (m *MsgAddAccountCertificates) FeePayer() sdk.AccAddress {
	if !m.FeePayerAddr.Empty() {
		return m.FeePayerAddr
	}
	return m.Owner
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
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgAddAccountCertificates) GetSigners() []sdk.AccAddress {
	if m.FeePayerAddr.Empty() {
		return []sdk.AccAddress{m.Owner}
	} else {
		return []sdk.AccAddress{m.FeePayerAddr, m.Owner}
	}
}

// MsgDeleteAccountCertificate is the request
// model used to remove certificates from an
// account
type MsgDeleteAccountCertificate struct {
	// Domain is the name of the domain of the account
	Domain string `json:"domain"`
	// Name is the name of the account
	Name string `json:"name"`
	// DeleteCertificate is the certificate to delete
	DeleteCertificate []byte `json:"delete_certificate"`
	// Owner is the owner of the account
	Owner        sdk.AccAddress `json:"owner"`
	FeePayerAddr sdk.AccAddress `json:"fee_payer"`
}

var _ fee.ProductMsg = (*MsgDeleteAccountCertificate)(nil)

// Route implements sdk.Msg
func (m *MsgDeleteAccountCertificate) FeePayer() sdk.AccAddress {
	if !m.FeePayerAddr.Empty() {
		return m.FeePayerAddr
	}
	return m.Owner
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
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgDeleteAccountCertificate) GetSigners() []sdk.AccAddress {
	if m.FeePayerAddr.Empty() {
		return []sdk.AccAddress{m.Owner}
	} else {
		return []sdk.AccAddress{m.FeePayerAddr, m.Owner}
	}
}

// MsgDeleteAccount is the request model
// used to delete an account
type MsgDeleteAccount struct {
	// Domain is the name of the domain of the account
	Domain string `json:"domain"`
	// Name is the name of the account
	Name string `json:"name"`
	// Owner is the owner of the account
	Owner        sdk.AccAddress `json:"owner"`
	FeePayerAddr sdk.AccAddress `json:"fee_payer"`
}

var _ fee.ProductMsg = (*MsgDeleteAccount)(nil)

// Route implements sdk.Msg
func (m *MsgDeleteAccount) FeePayer() sdk.AccAddress {
	if !m.FeePayerAddr.Empty() {
		return m.FeePayerAddr
	}
	return m.Owner
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
	if m.Domain == "" {
		return errors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Name == "" {
		return errors.Wrap(ErrOpEmptyAcc, "empty")
	}
	// success
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgDeleteAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgDeleteAccount) GetSigners() []sdk.AccAddress {
	if m.FeePayerAddr.Empty() {
		return []sdk.AccAddress{m.Owner}
	} else {
		return []sdk.AccAddress{m.FeePayerAddr, m.Owner}
	}
}

// MsgDeleteDomain is the request
// model to delete a domain
type MsgDeleteDomain struct {
	Domain       string         `json:"domain"`
	Owner        sdk.AccAddress `json:"owner"`
	FeePayerAddr sdk.AccAddress `json:"fee_payer"`
}

var _ fee.ProductMsg = (*MsgDeleteDomain)(nil)

func (m *MsgDeleteDomain) FeeCalculationFunc() (*fee.Params, error) {
	return nil, fee.ErrDefaultFee
}

// Route implements sdk.Msg
func (m *MsgDeleteDomain) FeePayer() sdk.AccAddress {
	if !m.FeePayerAddr.Empty() {
		return m.FeePayerAddr
	}
	return m.Owner
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
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgDeleteDomain) GetSigners() []sdk.AccAddress {
	if m.FeePayerAddr.Empty() {
		return []sdk.AccAddress{m.Owner}
	} else {
		return []sdk.AccAddress{m.FeePayerAddr, m.Owner}
	}
}

// MsgRegisterAccount is the request
// model used to register new accounts
type MsgRegisterAccount struct {
	// Domain is the domain of the account
	Domain string `json:"domain"`
	// Name is the name of the account
	Name string `json:"name"`
	// Owner is the owner of the account
	Owner sdk.AccAddress `json:"owner"`
	// Registerer is the user who registers this account
	Registerer sdk.AccAddress `json:"registerer"`
	// Targets are the blockchain addresses of the account
	Targets []BlockchainAddress `json:"targets"`
	// Broker is the account that facilitated the transaction
	Broker       sdk.AccAddress `json:"broker"`
	FeePayerAddr sdk.AccAddress `json:"fee_payer"`
}

var _ fee.ProductMsg = (*MsgRegisterAccount)(nil)

// Route implements sdk.Msg
func (m *MsgRegisterAccount) FeePayer() sdk.AccAddress {
	if !m.FeePayerAddr.Empty() {
		return m.FeePayerAddr
	}
	return m.Registerer
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
	if m.Registerer.Empty() {
		return errors.Wrap(ErrInvalidRegisterer, "empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgRegisterAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgRegisterAccount) GetSigners() []sdk.AccAddress {
	if m.FeePayerAddr.Empty() {
		return []sdk.AccAddress{m.Registerer}
	} else {
		return []sdk.AccAddress{m.FeePayerAddr, m.Registerer}
	}
}

// MsgRegisterDomain is the request used to register new domains
type MsgRegisterDomain struct {
	// Name is the name of the domain we want to register
	Name string `json:"domain" arg:"--domain" helper:"name of the domain"`
	// Admin is the address of the newly registered domain
	Admin sdk.AccAddress `json:"admin"`
	// DomainType defines the type of the domain
	DomainType DomainType `json:"type"`
	// Broker TODO document
	Broker       sdk.AccAddress `json:"broker" arg:"--broker" helper:"the broker"`
	FeePayerAddr sdk.AccAddress `json:"fee_payer"`
}

var _ fee.ProductMsg = (*MsgRegisterDomain)(nil)

// Route implements sdk.Msg
func (m *MsgRegisterDomain) FeePayer() sdk.AccAddress {
	if !m.FeePayerAddr.Empty() {
		return m.FeePayerAddr
	}
	return m.Admin
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
	if err := ValidateDomainType(m.DomainType); err != nil {
		return err
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
	if m.FeePayerAddr.Empty() {
		return []sdk.AccAddress{m.Admin}
	} else {
		return []sdk.AccAddress{m.FeePayerAddr, m.Admin}
	}
}

// MsgRenewAccount is the request
// model used to renew accounts
type MsgRenewAccount struct {
	// Domain is the domain of the account
	Domain string `json:"domain"`
	// Name is the name of the account
	Name string `json:"name"`
	// Signer is the signer of the request
	Signer       sdk.AccAddress `json:"signer"`
	FeePayerAddr sdk.AccAddress `json:"fee_payer"`
}

var _ fee.ProductMsg = (*MsgRenewAccount)(nil)

// Route implements sdk.Msg
func (m *MsgRenewAccount) FeePayer() sdk.AccAddress {
	if !m.FeePayerAddr.Empty() {
		return m.FeePayerAddr
	}
	return m.Signer
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
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgRenewAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgRenewAccount) GetSigners() []sdk.AccAddress {
	if m.FeePayerAddr.Empty() {
		return []sdk.AccAddress{m.Signer}
	} else {
		return []sdk.AccAddress{m.FeePayerAddr, m.Signer}
	}
}

// MsgRenewDomain is the request
// model used to renew a domain
type MsgRenewDomain struct {
	// Domain is the domain name to renew
	Domain string `json:"domain"`
	// Signer is the request signer
	Signer       sdk.AccAddress `json:"signer"`
	FeePayerAddr sdk.AccAddress `json:"fee_payer"`
}

var _ fee.ProductMsg = (*MsgRenewDomain)(nil)

// Route implements sdk.Msg
func (m *MsgRenewDomain) FeePayer() sdk.AccAddress {
	if !m.FeePayerAddr.Empty() {
		return m.FeePayerAddr
	}
	return m.Signer
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
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgRenewDomain) GetSigners() []sdk.AccAddress {
	if m.FeePayerAddr.Empty() {
		return []sdk.AccAddress{m.Signer}
	} else {
		return []sdk.AccAddress{m.FeePayerAddr, m.Signer}
	}
}

// MsgReplaceAccountTargets is the request model
// used to renew blockchain addresses associated
// with an account
type MsgReplaceAccountTargets struct {
	// Domain is the domain name of the account
	Domain string `json:"domain"`
	// Name is the name of the account
	Name string `json:"name"`
	// NewTargets are the new blockchain addresses
	NewTargets []BlockchainAddress `json:"new_targets"`
	// Owner is the owner of the account
	Owner        sdk.AccAddress `json:"owner"`
	FeePayerAddr sdk.AccAddress `json:"fee_payer"`
}

var _ fee.ProductMsg = (*MsgReplaceAccountTargets)(nil)

// Route implements sdk.Msg
func (m *MsgReplaceAccountTargets) FeePayer() sdk.AccAddress {
	if !m.FeePayerAddr.Empty() {
		return m.FeePayerAddr
	}
	return m.Owner
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
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgReplaceAccountTargets) GetSigners() []sdk.AccAddress {
	if m.FeePayerAddr.Empty() {
		return []sdk.AccAddress{m.Owner}
	} else {
		return []sdk.AccAddress{m.FeePayerAddr, m.Owner}
	}
}

// MsgReplaceAccountMetadata is the function used
// to set accounts metadata
type MsgReplaceAccountMetadata struct {
	// Domain is the domain name of the account
	Domain string `json:"domain"`
	// Name is the name of the account
	Name string `json:"name"`
	// NewMetadataURI is the metadata URI of the account
	// we want to update or insert
	NewMetadataURI string `json:"new_metadata_uri"`
	// Owner is the owner of the account
	Owner        sdk.AccAddress `json:"owner"`
	FeePayerAddr sdk.AccAddress `json:"fee_payer"`
}

var _ fee.ProductMsg = (*MsgReplaceAccountMetadata)(nil)

// Route implements sdk.Msg
func (m *MsgReplaceAccountMetadata) FeePayer() sdk.AccAddress {
	if !m.FeePayerAddr.Empty() {
		return m.FeePayerAddr
	}
	return m.Owner
}

// Route implements sdk.Msg
func (m *MsgReplaceAccountMetadata) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgReplaceAccountMetadata) Type() string {
	return "set_account_metadata"
}

// ValidateBasic implements sdk.Msg
func (m *MsgReplaceAccountMetadata) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrapf(ErrInvalidDomainName, "empty")
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
func (m *MsgReplaceAccountMetadata) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgReplaceAccountMetadata) GetSigners() []sdk.AccAddress {
	if m.FeePayerAddr.Empty() {
		return []sdk.AccAddress{m.Owner}
	} else {
		return []sdk.AccAddress{m.FeePayerAddr, m.Owner}
	}
}

// MsgTransferAccount is the request
// model used to transfer accounts
type MsgTransferAccount struct {
	// Domain is the domain name of the account
	Domain string `json:"domain"`
	// Account is the account name
	Name string `json:"name"`
	// Owner is the actual owner of the account
	Owner sdk.AccAddress `json:"owner"`
	// NewOwner is the new owner of the account
	NewOwner sdk.AccAddress `json:"new_owner"`
	// Reset indicates if the accounts content will be resetted
	Reset        bool           `json:"reset"`
	FeePayerAddr sdk.AccAddress `json:"fee_payer"`
}

var _ fee.ProductMsg = (*MsgTransferAccount)(nil)

// Route implements sdk.Msg
func (m *MsgTransferAccount) FeePayer() sdk.AccAddress {
	if !m.FeePayerAddr.Empty() {
		return m.FeePayerAddr
	}
	return m.Owner
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
		return errors.Wrap(ErrOpEmptyAcc, "empty")
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
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgTransferAccount) GetSigners() []sdk.AccAddress {
	if m.FeePayerAddr.Empty() {
		return []sdk.AccAddress{m.Owner}
	} else {
		return []sdk.AccAddress{m.FeePayerAddr, m.Owner}
	}
}

// TransferFlag defines the type of domain transfer
type TransferFlag int

const (
	// TransferFlush clears all domain account data, except empty account)
	TransferFlush = iota
	// TransferOwned transfers only accounts owned by the current owner
	TransferOwned
	// ResetNone leaves things as they are except for empty account
	ResetNone
	// TransferAll is not available is here only for tests backwards compatibility and will be removed. TODO deprecate
	TransferAll
)

// MsgTransferDomain is the request model
// used to transfer a domain
type MsgTransferDomain struct {
	Domain       string         `json:"domain"`
	Owner        sdk.AccAddress `json:"owner"`
	NewAdmin     sdk.AccAddress `json:"new_admin"`
	TransferFlag TransferFlag   `json:"transfer_flag"`
	FeePayerAddr sdk.AccAddress `json:"fee_payer"`
}

var _ fee.ProductMsg = (*MsgTransferDomain)(nil)

// Route implements sdk.Msg
func (m *MsgTransferDomain) FeePayer() sdk.AccAddress {
	if !m.FeePayerAddr.Empty() {
		return m.FeePayerAddr
	}
	return m.Owner
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
	switch m.TransferFlag {
	case TransferOwned:
	case ResetNone:
	case TransferFlush:
	default:
		return errors.Wrapf(ErrInvalidRequest, "unknown reset flag: %d", m.TransferFlag)
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgTransferDomain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgTransferDomain) GetSigners() []sdk.AccAddress {
	if m.FeePayerAddr.Empty() {
		return []sdk.AccAddress{m.Owner}
	} else {
		return []sdk.AccAddress{m.FeePayerAddr, m.Owner}
	}
}

func snakeCaseAppend(args ...string) string {
	var str string
	for _, a := range args {
		str = a + "_"
	}
	return str
}

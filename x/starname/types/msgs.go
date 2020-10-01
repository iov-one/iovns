package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgWithFeePayer abstracts the Msg type to support a fee payer
// which takes care of handling product fees
type MsgWithFeePayer interface {
	sdk.Msg
	FeePayer() sdk.AccAddress
}

var _ MsgWithFeePayer = (*MsgAddAccountCertificates)(nil)

// FeePayer implements FeePayer interface
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

var _ MsgWithFeePayer = (*MsgDeleteAccountCertificate)(nil)

// FeePayer implements FeePayer interface
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

var _ MsgWithFeePayer = (*MsgDeleteAccount)(nil)

// FeePayer implements FeePayer interface
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

var _ MsgWithFeePayer = (*MsgDeleteDomain)(nil)

// FeePayer implements FeePayer interface
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

var _ MsgWithFeePayer = (*MsgRegisterAccount)(nil)

// FeePayer implements FeePayer interface
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

var _ MsgWithFeePayer = (*MsgRegisterDomain)(nil)

// FeePayer implements FeePayer interface
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

var _ MsgWithFeePayer = (*MsgRenewAccount)(nil)

// FeePayer implements FeePayer interface
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

var _ MsgWithFeePayer = (*MsgRenewDomain)(nil)

// FeePayer implements FeePayer interface
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

var _ MsgWithFeePayer = (*MsgReplaceAccountResources)(nil)

// FeePayer implements FeePayer interface
func (m *MsgReplaceAccountResources) FeePayer() sdk.AccAddress {
	if !m.FeePayerAddr.Empty() {
		return m.FeePayerAddr
	}
	return m.Owner
}

// Route implements sdk.Msg
func (m *MsgReplaceAccountResources) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgReplaceAccountResources) Type() string {
	return "replace_account_resources"
}

// ValidateBasic implements sdk.Msg
func (m *MsgReplaceAccountResources) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Owner == nil {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	if len(m.NewResources) == 0 {
		return errors.Wrap(ErrInvalidRequest, "empty resources")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgReplaceAccountResources) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgReplaceAccountResources) GetSigners() []sdk.AccAddress {
	if m.FeePayerAddr.Empty() {
		return []sdk.AccAddress{m.Owner}
	} else {
		return []sdk.AccAddress{m.FeePayerAddr, m.Owner}
	}
}

var _ MsgWithFeePayer = (*MsgReplaceAccountMetadata)(nil)

// FeePayer implements FeePayer interface
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

var _ MsgWithFeePayer = (*MsgTransferAccount)(nil)

// FeePayer implements FeePayer interface
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
	// TransferResetNone leaves things as they are except for empty account
	TransferResetNone
	// TransferAll is not available is here only for tests backwards compatibility and will be removed. TODO deprecate
	TransferAll
)

var _ MsgWithFeePayer = (*MsgTransferDomain)(nil)

// FeePayer implements FeePayer interface
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
	case TransferResetNone:
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

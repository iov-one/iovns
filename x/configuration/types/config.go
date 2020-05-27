package types

import (
	"fmt"
	"regexp"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Config is the configuration of the network
type Config struct {
	// Configurer is the configuration owner, the addresses allowed to handle fees
	// and register domains with no superuser
	Configurer sdk.AccAddress `json:"configurer"`
	// ValidDomainName defines a regexp that determines if a domain name is valid or not
	ValidDomainName string `json:"valid_domain_name"`
	// ValidAccountName defines a regexp that determines if an account name is valid or not
	ValidAccountName string `json:"valid_account_name"`
	// ValidBlockchainID defines a regexp that determines if a blockchain id is valid or not
	ValidBlockchainID string `json:"valid_blockchain_id"`
	// ValidBlockchainAddress determines a regexp for a valid blockchain address
	ValidBlockchainAddress string `json:"valid_blockchain_address"`

	// DomainRenewalPeriod defines the duration of the domain renewal period in seconds
	DomainRenewalPeriod time.Duration `json:"domain_renew_period"`
	// DomainRenewalCountMax defines maximum number of domain renewals a user can do
	DomainRenewalCountMax uint32 `json:"domain_renew_count_max"`
	// DomainGracePeriod defines the grace period for a domain deletion in seconds
	DomainGracePeriod time.Duration `json:"domain_grace_period"`

	// AccountRenewalPeriod defines the duration of the account renewal period in seconds
	AccountRenewalPeriod time.Duration `json:"account_renew_period"`
	// AccountRenewalCountMax defines maximum number of account renewals a user can do
	AccountRenewalCountMax uint32 `json:"account_renew_count_max"`
	// DomainGracePeriod defines the grace period for a domain deletion in seconds
	AccountGracePeriod time.Duration `json:"account_grace_period"`
	// BlockchainTargetMax defines maximum number of blockchain targets could be saved under an account
	BlockchainTargetMax uint32 `json:"blockchain_target_max"`
	// CertificateSizeMax defines maximum size of a certificate that could be saved under an account
	CertificateSizeMax uint64 `json:"certificate_size_max"`
	// CertificateCountMax defines maximum number of certificates that could be saved under an account
	CertificateCountMax uint32 `json:"certificate_count_max"`
	// MetadataSizeMax defines maximum size of metadata that could be saved under an account
	MetadataSizeMax uint64 `json:"metadata_size_max"`
}

func (c Config) Validate() error {
	if c.Configurer == nil {
		return fmt.Errorf("empty configurer")
	}
	if c.DomainRenewalPeriod < 0 {
		return fmt.Errorf("empty domain renew")
	}
	if c.DomainGracePeriod < 0 {
		return fmt.Errorf("empty domain grace period")
	}
	if c.AccountRenewalPeriod < 0 {
		return fmt.Errorf("empty account renew")
	}
	if c.AccountGracePeriod < 0 {
		return fmt.Errorf("empty account grace period")
	}
	if _, err := regexp.Compile(c.ValidAccountName); err != nil {
		return err
	}
	if _, err := regexp.Compile(c.ValidBlockchainAddress); err != nil {
		return err
	}
	if _, err := regexp.Compile(c.ValidBlockchainID); err != nil {
		return err
	}
	if _, err := regexp.Compile(c.ValidDomainName); err != nil {
		return err
	}

	return nil
}

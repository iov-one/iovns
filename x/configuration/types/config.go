package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

// Config is the configuration of the network
type Config struct {
	// Owners are the configuration owner, the addresses allowed to handle fees
	// and register domains with no superuser
	Owners []sdk.AccAddress `json:"owners"`
	// ValidDomain defines a regexp that determines if a domain name is valid or not
	ValidDomain string `json:"valid_domain"`
	// ValidName defines a regexp that determines if an account name is valid or not
	ValidName string `json:"valid_name"`
	// ValidBlockchainID defines a regexp that determines if a blockchain id is valid or not
	ValidBlockchainID string `json:"valid_blockchain_id"`
	// ValidBlockchainAddress determines a regexp for a valid blockchain address
	ValidBlockchainAddress string `json:"valid_blockchain_address"`
	// DomainRenew defines the duration of the domain renewal period
	DomainRenew time.Duration `json:"domain_renew"`
	// DomainGracePeriod defines the grace period for a domain deletion in seconds
	DomainGracePeriod int64 `json:"domain_grace_period"`
}

package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns"
)

// Domain defines a domain
type Domain struct {
	// Name is the name of the domain
	Name string `json:"name"`
	// Admin is the owner of the domain
	Admin sdk.AccAddress `json:"admin"`
	// ValidUntil is a unix timestamp that defines for how long the domain is valid
	ValidUntil int64 `json:"valid_until"`
	// HasSuperuser checks if the domain is owned by a super user or not
	HasSuperuser bool `json:"has_super_user"`
	// AccountRenew defines the duration of each created or renewed account
	// under the domain
	AccountRenew time.Duration `json:"account_renew"`
	// Broker TODO needs comment
	Broker sdk.AccAddress `json:"broker"`
}

// Account defines an account that belongs to a domain
type Account struct {
	// Domain references the domain this account belongs to
	Domain string `json:"domain"`
	// Name is the name of the account
	Name string `json:"name"`
	// Owner is the address that owns the account
	Owner sdk.AccAddress `json:"owner"`
	// ValidUntil defines a unix timestamp of the expiration of the account
	ValidUntil int64 `json:"valid_until"`
	// Targets is the list of blockchain addresses this account belongs to
	Targets []iovns.BlockchainAddress `json:"targets"`
	// Certificates contains the list of certificates to identify the account owner
	Certificates [][]byte `json:"certificates"`
	// Broker can be empty
	// it identifies an entity that facilitated the transaction of the account
	Broker sdk.AccAddress `json:"broker"`
	// MetadataURI contains a link to extra information regarding the account
	MetadataURI string `json:"metadata_uri"`
}

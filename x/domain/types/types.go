package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovnsd"
	"time"
)

// Domain defines a domain
type Domain struct {
	// Name is the name of the domain
	Name string
	// Admin is the owner of the domain
	Admin sdk.AccAddress
	// ValidUntil is a unix timestamp that defines for how long the domain is valid
	ValidUntil int64
	// HasSuperuser checks if the domain is owned by a super user or not
	HasSuperuser bool
	// AccountRenew defines the duration of each created or renewed account
	// under the domain
	AccountRenew time.Duration
	// Broker TODO needs comment
	Broker sdk.AccAddress
}

func (d Domain) String() string {
	panic("implement plz")
}

// Account defines an account that belongs to a domain
type Account struct {
	// Domain references the domain this account belongs to
	Domain string
	// Name is the name of the account
	Name string
	// Owner is the address that owns the account
	Owner sdk.AccAddress
	// ValidUntil defines a unix timestamp of the expiration of the account
	ValidUntil int64
	// Targets is the list of blockchain addresses this account belongs to
	Targets []iovnsd.BlockchainAddress
	// Certificates contains the list of certificates to identify the account owner
	Certificates [][]byte
	// Broker can be empty
	// it identifies an entity that facilitated the transaction of the account
	Broker sdk.AccAddress
}

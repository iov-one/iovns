package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/pkg/index"
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

// Account defines an account that belongs to a domain

// owner:account
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
	Targets []BlockchainAddress
	// Certificates contains the list of certificates to identify the account owner
	Certificates []Certificate
	// Broker can be empty
	// it identifies an entity that facilitated the transaction of the account
	Broker sdk.AccAddress
	// MetadataURI contains a link to extra information regarding the account
	MetadataURI string
}

func (a Account) Pack() ([]byte, error) {
	key, err := index.PackBytes([][]byte{[]byte(a.Domain), []byte(a.Name)})
	if err != nil {
		return nil, err
	}
	return key, nil
}

// Unpack converts a byte key into an account
// composed only of domain and name
func (a *Account) Unpack(key []byte) error {
	keys, err := index.UnpackBytes(key)
	if err != nil {
		return err
	}
	if len(keys) != 2 {
		return fmt.Errorf("unexpected number of keys for %T: %d", a, len(keys))
	}
	a.Domain = string(keys[0])
	a.Name = string(keys[1])
	return nil
}

// BlockchainAddress defines an address coming from different DLTs
type BlockchainAddress struct {
	// ID defines a blockchain ID
	ID string
	// Address is the blockchain address
	Address string
}

func (b BlockchainAddress) Index() ([]byte, error) {
	return index.PackBytes([][]byte{[]byte(b.ID), []byte(b.Address)})
}

// Certificate defines a certificate
type Certificate []byte

func (c Certificate) Index() ([]byte, error) {
	return c, nil
}

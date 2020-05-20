package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/pkg/index"
)

// emptyAccountNameIndexIdentifier defines how empty
// account names of a domain are identified in indexes
const emptyAccountNameIndexIdentifier = "*"

// Domain defines a domain
type Domain struct {
	// Name is the name of the domain
	Name string `json:"name"`
	// Admin is the owner of the domain
	Admin sdk.AccAddress `json:"admin"`
	// ValidUntil is a unix timestamp defines the time when the domain will become invalid
	ValidUntil int64 `json:"valid_until"`
	// HasSuperuser checks if the domain is owned by a super user or not
	HasSuperuser bool `json:"has_super_user"`
	// AccountRenew defines the duration of each created or renewed account
	// under the domain
	AccountRenew time.Duration `json:"account_renew"`
	// Broker TODO needs comment
	Broker sdk.AccAddress `json:"broker"`
}

// Index implements Indexer and packs the
// domain into an index key using its name
func (d Domain) Index() ([]byte, error) {
	key, err := index.PackBytes([][]byte{[]byte(d.Name)})
	if err != nil {
		return nil, err
	}
	return key, nil
}

// Pack implements Indexed and allows
// the domain to be saved as a value
// in an index deterministically
func (d Domain) Pack() ([]byte, error) {
	return d.Index()
}

// Unpack implements Unpacker and allows
// the domain to be retrieved from an index key
func (d *Domain) Unpack(key []byte) error {
	unpackedKeys, err := index.UnpackBytes(key)
	if err != nil {
		return err
	}
	if len(unpackedKeys) != 1 {
		return fmt.Errorf("unpack domain expected one key, got: %d", len(unpackedKeys))
	}
	d.Name = string(unpackedKeys[0])
	return nil
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
	Targets []BlockchainAddress `json:"targets"`
	// Certificates contains the list of certificates to identify the account owner
	Certificates []Certificate `json:"certificates"`
	// Broker can be empty
	// it identifies an entity that facilitated the transaction of the account
	Broker sdk.AccAddress `json:"broker"`
	// MetadataURI contains a link to extra information regarding the account
	MetadataURI string `json:"metadata_uri"`
}

// Pack implements Indexed and allows
// the account to be saved as a value
// in an index deterministically
func (a Account) Pack() ([]byte, error) {
	// in order to avoid empty keys
	// indexing in case account name
	// is empty, we index it as '*'
	var name = a.Name
	if a.Name == "" {
		name = emptyAccountNameIndexIdentifier
	}
	key, err := index.PackBytes([][]byte{[]byte(a.Domain), []byte(name)})
	if err != nil {
		return nil, err
	}
	return key, nil
}

// Unpack implements Unpacker and allows
// the account to be retrieved from an index key
func (a *Account) Unpack(key []byte) error {
	keys, err := index.UnpackBytes(key)
	if err != nil {
		return err
	}
	if len(keys) != 2 {
		return fmt.Errorf("unexpected number of keys for %T: %d", a, len(keys))
	}
	a.Domain = string(keys[0])
	name := string(keys[1])
	if name == emptyAccountNameIndexIdentifier {
		name = ""
	}
	a.Name = name
	return nil
}

// BlockchainAddress defines an address coming from different DLTs
type BlockchainAddress struct {
	// ID defines a blockchain ID
	ID string `json:"id"`
	// Address is the blockchain address
	Address string `json:"address"`
}

// Index implements Indexer and packs the
// blockchain address into an index key using
// its blockchain ID and address
func (b BlockchainAddress) Index() ([]byte, error) {
	return index.PackBytes([][]byte{[]byte(b.ID), []byte(b.Address)})
}

// Certificate defines a certificate
type Certificate []byte

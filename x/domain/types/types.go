package types

import (
	"bytes"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/pkg/utils"
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
	Certificates [][]byte
	// Broker can be empty
	// it identifies an entity that facilitated the transaction of the account
	Broker sdk.AccAddress
	// MetadataURI contains a link to extra information regarding the account
	MetadataURI string
}

func (a Account) Index() []byte {
	encodedDomain := utils.Base64Encode(a.Domain)
	encodedName := utils.Base64Encode(a.Name)
	return bytes.Join([][]byte{encodedDomain, encodedName}, iovns.Separator)
}

// Unpack converts a byte key into an account
// composed only of domain and name
func (a *Account) Unpack(key []byte) error {
	if a == nil {
		*a = Account{}
	}
	splits := bytes.Split(key, iovns.Separator)
	if len(splits) != 2 {
		return fmt.Errorf("unpack: unxpected number of splits from key %x: %d", key, len(splits))
	}
	domain, err := utils.Base64Decode(splits[0])
	if err != nil {
		return fmt.Errorf("unpack: %s", err)
	}
	a.Domain = domain
	name, err := utils.Base64Decode(splits[1])
	if err != nil {
		return fmt.Errorf("unpack: %s", err)
	}
	a.Name = name
	return nil
}

// BlockchainAddress defines an address coming from different DLTs
type BlockchainAddress struct {
	// ID defines a blockchain ID
	ID string
	// Address is the blockchain address
	Address string
}

func (b BlockchainAddress) Index() []byte {
	encodedID := utils.Base64Encode(b.ID)
	encodedAddress := utils.Base64Encode(b.Address)
	return bytes.Join([][]byte{encodedID, encodedAddress}, iovns.Separator)
}

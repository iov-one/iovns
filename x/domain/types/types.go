package types

import (
	"github.com/iov-one/iovns/pkg/crud"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/pkg/index"
)

// emptyAccountNameIndexIdentifier defines how empty
// account names of a domain are identified in indexes
const emptyAccountNameIndexIdentifier = "*"

const DomainAdminIndex = 0x1
const AccountAdminIndex = 0x1
const AccountDomainIndex = 0x2
const AccountResourcesIndex = 0x3

// Domain defines a domain
type Domain struct {
	// Name is the name of the domain
	Name string `json:"name" crud:"primaryKey"`
	// Owner is the owner of the domain
	Admin sdk.AccAddress `json:"admin" crud:"01"`
	// ValidUntil is a unix timestamp defines the time when the domain will become invalid
	ValidUntil int64 `json:"valid_until"`
	// Type defines the type of the domain
	Type DomainType `json:"type"`
	// Broker TODO needs comment
	Broker sdk.AccAddress `json:"broker"`
}

func (d *Domain) PrimaryKey() crud.PrimaryKey {
	if d.Name == "" {
		return nil
	}
	return []byte(d.Name)
}

func (d *Domain) SecondaryKeys() []crud.SecondaryKey {
	return []crud.SecondaryKey{
		{
			Key:         d.Admin,
			StorePrefix: []byte{DomainAdminIndex},
		},
	}
}

type DomainType string

const (
	OpenDomain   DomainType = "open"
	ClosedDomain            = "closed"
)

func ValidateDomainType(typ DomainType) error {
	switch typ {
	case OpenDomain, ClosedDomain:
		return nil
	default:
		return errors.Wrapf(ErrInvalidDomainType, "invalid domain type: %s", typ)
	}
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
	// Resources is the list of resources an account resolves to
	Resources []Resource `json:"resources"`
	// Certificates contains the list of certificates to identify the account owner
	Certificates []Certificate `json:"certificates"`
	// Broker can be empty
	// it identifies an entity that facilitated the transaction of the account
	Broker sdk.AccAddress `json:"broker"`
	// MetadataURI contains a link to extra information regarding the account
	MetadataURI string `json:"metadata_uri"`
}

func (a *Account) PrimaryKey() crud.PrimaryKey {
	j := strings.Join([]string{a.Domain, a.Name}, "*")
	return []byte(j)
}

func (a *Account) SecondaryKeys() []crud.SecondaryKey {
	// index by owner
	ownerIndex := crud.SecondaryKey{
		Key:         a.Owner,
		StorePrefix: []byte{AccountAdminIndex},
	}
	// index by domain
	domainIndex := crud.SecondaryKey{
		Key:         []byte(a.Domain),
		StorePrefix: []byte{AccountDomainIndex},
	}
	// index by resources
	resourcesIndexes := make([]crud.SecondaryKey, len(a.Resources))
	for i, res := range a.Resources {
		// exclude empty resources
		if res.Resource == "" || res.URI == "" {
			continue
		}
		resKey := strings.Join([]string{res.URI, res.Resource}, "")
		resourcesIndexes[i] = crud.SecondaryKey{
			Key:         []byte(resKey),
			StorePrefix: []byte{AccountResourcesIndex},
		}
	}
	// return keys
	return append([]crud.SecondaryKey{ownerIndex, domainIndex}, resourcesIndexes...)
}

// Resource defines a resource an account can resolve to
type Resource struct {
	// URI defines the ID of the resource
	URI string `json:"uri"`
	// Resource is the resource
	Resource string `json:"resource"`
}

// Index implements Indexer and packs the
// resource into an index key using
// its URI and Resource
func (b Resource) Index() ([]byte, error) {
	return index.PackBytes([][]byte{[]byte(b.URI), []byte(b.Resource)})
}

// Certificate defines a certificate
type Certificate []byte

package types

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	crud "github.com/iov-one/cosmos-sdk-crud/pkg/crud/types"
	"github.com/iov-one/iovns/pkg/utils"

	"github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const DomainAdminIndex = 0x1
const AccountAdminIndex = 0x1
const AccountDomainIndex = 0x2
const AccountResourcesIndex = 0x3

// StarnameSeparator defines the starname separator identifier
const StarnameSeparator = "*"

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
	return crud.NewPrimaryKeyFromString(d.Name)
}

func (d *Domain) SecondaryKeys() []crud.SecondaryKey {
	if d.Admin.Empty() {
		return nil
	}
	return []crud.SecondaryKey{crud.NewSecondaryKey(DomainAdminIndex, d.Admin)}
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
// NOTE: It should not be confused with cosmos-sdk auth account
// github.com/cosmos/cosmos-sdk/x/auth.Account
type Account struct {
	// Domain references the domain this account belongs to
	Domain string `json:"domain"`
	// Name is the name of the account
	Name *string `json:"name"`
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

type accountCodecAlias struct {
	Underlying *Account
	NameNil    bool
}

func (a *Account) MarshalCRUD() interface{} {
	return accountCodecAlias{
		Underlying: a,
		NameNil:    a.Name == nil,
	}
}

func (a *Account) UnmarshalCRUD(cdc *codec.Codec, b []byte) {
	trg := new(accountCodecAlias)
	cdc.MustUnmarshalBinaryBare(b, trg)
	*a = *trg.Underlying
	if a.Name == nil && !trg.NameNil {
		a.Name = utils.StrPtr("")
	}
}

func (a *Account) PrimaryKey() crud.PrimaryKey {
	if len(a.Domain) == 0 || a.Name == nil {
		return nil
	}
	j := strings.Join([]string{a.Domain, *a.Name}, "*")
	return crud.NewPrimaryKeyFromString(j)
}

func (a *Account) SecondaryKeys() []crud.SecondaryKey {
	var sk []crud.SecondaryKey
	// index by owner
	if !a.Owner.Empty() {
		ownerIndex := crud.NewSecondaryKey(AccountAdminIndex, a.Owner)
		sk = append(sk, ownerIndex)
	}
	// index by domain
	if len(a.Domain) != 0 {
		domainIndex := crud.NewSecondaryKey(AccountDomainIndex, []byte(a.Domain))
		sk = append(sk, domainIndex)
	}
	// index by resources
	for _, res := range a.Resources {
		// exclude empty resources
		if res.Resource == "" || res.URI == "" {
			continue
		}
		resKey := strings.Join([]string{res.URI, res.Resource}, "")
		// append resource
		sk = append(sk, crud.NewSecondaryKey(AccountResourcesIndex, []byte(resKey)))
	}
	// return keys
	return sk
}

// Resource defines a resource an account can resolve to
type Resource struct {
	// URI defines the ID of the resource
	URI string `json:"uri"`
	// Resource is the resource
	Resource string `json:"resource"`
}

// Certificate defines a certificate
type Certificate []byte

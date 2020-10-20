package types

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	crud "github.com/iov-one/cosmos-sdk-crud/pkg/crud/types"
	"github.com/iov-one/iovns/pkg/utils"

	"github.com/cosmos/cosmos-sdk/types/errors"
)

const DomainAdminIndex = 0x1
const AccountAdminIndex = 0x1
const AccountDomainIndex = 0x2
const AccountResourcesIndex = 0x3

// StarnameSeparator defines the starname separator identifier
const StarnameSeparator = "*"

func (m *Domain) PrimaryKey() crud.PrimaryKey {
	if m.Name == "" {
		return nil
	}
	return crud.NewPrimaryKeyFromString(m.Name)
}

func (m *Domain) SecondaryKeys() []crud.SecondaryKey {
	if m.Admin.Empty() {
		return nil
	}
	return []crud.SecondaryKey{crud.NewSecondaryKey(DomainAdminIndex, m.Admin)}
}

// DomainType defines the type of the domain
type DomainType string

const (
	// OpenDomain is the domain type in which an account owner is the only entity that can perform actions on the account
	OpenDomain DomainType = "open"
	// ClosedDomain is the domain type in which the domain owner has control over accounts too
	ClosedDomain = "closed"
)

func ValidateDomainType(typ DomainType) error {
	switch typ {
	case OpenDomain, ClosedDomain:
		return nil
	default:
		return errors.Wrapf(ErrInvalidDomainType, "invalid domain type: %s", typ)
	}
}

type accountCodecAlias struct {
	Underlying *Account
	NameNil    bool
}

func (m *Account) MarshalCRUD() interface{} {
	return accountCodecAlias{
		Underlying: m,
		NameNil:    m.Name == nil,
	}
}

func (m *Account) UnmarshalCRUD(cdc *codec.Codec, b []byte) {
	trg := new(accountCodecAlias)
	cdc.MustUnmarshalBinaryBare(b, trg)
	*m = *trg.Underlying
	if m.Name == nil && !trg.NameNil {
		m.Name = utils.StrPtr("")
	}
}

func (m *Account) PrimaryKey() crud.PrimaryKey {
	if len(m.Domain) == 0 || m.Name == nil {
		return nil
	}
	j := strings.Join([]string{m.Domain, *m.Name}, "*")
	return crud.NewPrimaryKeyFromString(j)
}

func (m *Account) SecondaryKeys() []crud.SecondaryKey {
	var sk []crud.SecondaryKey
	// index by owner
	if !m.Owner.Empty() {
		ownerIndex := crud.NewSecondaryKey(AccountAdminIndex, m.Owner)
		sk = append(sk, ownerIndex)
	}
	// index by domain
	if len(m.Domain) != 0 {
		domainIndex := crud.NewSecondaryKey(AccountDomainIndex, []byte(m.Domain))
		sk = append(sk, domainIndex)
	}
	// index by resources
	for _, res := range m.Resources {
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

// Certificate defines a certificate
type Certificate []byte

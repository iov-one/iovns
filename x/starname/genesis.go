package starname

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crud "github.com/iov-one/cosmos-sdk-crud/pkg/crud/types"
	"github.com/iov-one/iovns/x/starname/types"
)

// GenesisState represents the state of the domain module
type GenesisState struct {
	// DomainRecords contains the records of registered domains
	Domains []types.Domain `json:"domains"`
	// AccountRecords contains the records of registered accounts
	Accounts []types.Account `json:"accounts"`
}

// NewGenesisState builds a genesis state including the domains provided
func NewGenesisState(domains []types.Domain, accounts []types.Account) GenesisState {
	return GenesisState{Domains: domains, Accounts: accounts}
}

// ValidateGenesis validates a genesis state
// checking for domain validity and no domain name repetitions
func ValidateGenesis(data GenesisState) error {
	namesSet := make(map[string]struct{}, len(data.Domains))
	for _, domain := range data.Domains {
		if _, ok := namesSet[domain.Name]; ok {
			return fmt.Errorf("domain name %s declared twice", domain.Name)
		}
		namesSet[domain.Name] = struct{}{}
		if err := validateDomain(domain); err != nil {
			return err
		}
	}
	return nil
}

// DefaultGenesisState creates an empty genesis state for the domain module
func DefaultGenesisState() GenesisState {
	return GenesisState{Domains: []types.Domain{}, Accounts: []types.Account{}}
}

// InitGenesis builds a state from GenesisState
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	// insert domains
	ds := keeper.DomainStore(ctx)
	for _, domain := range data.Domains {
		// create domains
		ds.Create(&domain)
	}
	// insert accounts
	as := keeper.AccountStore(ctx)
	for _, account := range data.Accounts {
		as.Create(&account)
	}
}

// ExportGenesis saves the state of the domain module
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	ds := k.DomainStore(ctx)
	as := k.AccountStore(ctx)
	var domains []types.Domain
	var accounts []types.Account
	ds.IterateKeys(func(pk crud.PrimaryKey) bool {
		domain := new(types.Domain)
		ds.Read(pk, domain)
		domains = append(domains, *domain)
		return true
	})
	as.IterateKeys(func(pk crud.PrimaryKey) bool {
		account := new(types.Account)
		as.Read(pk, account)
		accounts = append(accounts, *account)
		return true
	})
	return GenesisState{
		Domains:  domains,
		Accounts: accounts,
	}
}

// validateDomain checks if a domain is valid or not
func validateDomain(d types.Domain) error {
	// TODO fill
	return nil
}

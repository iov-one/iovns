package domain

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/types"
)

// GenesisState represents the state of the domain module
type GenesisState struct {
	// DomainRecords contains the records of registered domains
	DomainsRecords []types.Domain `json:"domain_records"`
	// AccountRecords contains the records of registered accounts
	AccountsRecords []types.Account `json:"account_records"`
}

// NewGenesisState builds a genesis state including the domains provided
func NewGenesisState(domains []types.Domain, accounts []types.Account) GenesisState {
	return GenesisState{DomainsRecords: domains, AccountsRecords: accounts}
}

// ValidateGenesis validates a genesis state
// checking for domain validity and no domain name repetitions
func ValidateGenesis(data GenesisState) error {
	namesSet := make(map[string]struct{}, len(data.DomainsRecords))
	for _, domain := range data.DomainsRecords {
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
	return GenesisState{DomainsRecords: []types.Domain{}, AccountsRecords: []types.Account{}}
}

// InitGenesis builds a state from GenesisState
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	// insert domains
	for _, domain := range data.DomainsRecords {
		keeper.CreateDomain(ctx, domain)
	}
	// insert accounts
	for _, account := range data.AccountsRecords {
		keeper.CreateAccount(ctx, account)
	}
}

// ExportGenesis saves the state of the domain module
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	// save domain data
	domains := k.IterateAllDomains(ctx)
	// save account data
	accounts := k.IterateAllAccounts(ctx)
	return GenesisState{DomainsRecords: domains, AccountsRecords: accounts}
}

// validateDomain checks if a domain is valid or not
func validateDomain(d types.Domain) error {
	// TODO fill
	return nil
}

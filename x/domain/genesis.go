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
}

// NewGenesisState builds a genesis state including the domains provided
func NewGenesisState(domains []types.Domain) GenesisState {
	return GenesisState{DomainsRecords: domains}
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
	// TODO remove.
	return GenesisState{DomainsRecords: []types.Domain{
		{
			Name:         "test",
			Admin:        nil,
			ValidUntil:   0,
			HasSuperuser: false,
			AccountRenew: 0,
			Broker:       nil,
		},
	}}
}

// InitGenesis builds a state from GenesisState
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, domain := range data.DomainsRecords {
		keeper.CreateDomain(ctx, domain)
	}
}

// ExportGenesis saves the state of the domain module
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []types.Domain
	iterator := k.IterateAllDomains(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		domain, _ := k.GetDomain(ctx, string(iterator.Key()))
		records = append(records, domain)
	}
	return GenesisState{DomainsRecords: records}
}

// validateDomain checks if a domain is valid or not
func validateDomain(d types.Domain) error {
	// TODO fill
	return nil
}

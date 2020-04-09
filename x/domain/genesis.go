package domain

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovnsd/x/domain/types"
)

type GenesisState struct {
	DomainsRecords []types.Domain `json:"domain_records"`
}

func NewGenesisState(domains []types.Domain) GenesisState {
	return GenesisState{DomainsRecords: domains}
}

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

// validateDomain checks if a domain is valid or not
func validateDomain(d types.Domain) error {
	// TODO fill
	return nil
}

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

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, domain := range data.DomainsRecords {
		keeper.SetDomain(ctx, domain)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []types.Domain
	iterator := k.IterateAll(ctx)
	for ; iterator.Valid(); iterator.Next() {
		domain, _ := k.GetDomain(ctx, string(iterator.Key()))
		records = append(records, domain)
	}
	return GenesisState{DomainsRecords: records}
}

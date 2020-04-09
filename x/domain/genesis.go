package domain

import sdk "github.com/cosmos/cosmos-sdk/types"

type GenesisState struct {
	DomainsRecords []Domain `json:"domain_records"`
}

func NewGenesisState(domains []Domain) GenesisState {
	return GenesisState{DomainsRecords: domains}
}

func ValidateGenesis(data GenesisState) error {
	// TODO validate genesis by: checking no duplicate names, and domain validity
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{DomainsRecords: nil}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, domain := range data.DomainsRecords {
		keeper.SetDomain(ctx, domain)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []Domain
	iterator := k.IterateAll(ctx)
	for ; iterator.Valid(); iterator.Next() {
		domain, _ := k.GetDomain(ctx, string(iterator.Key()))
		records = append(records, domain)
	}
	return GenesisState{DomainsRecords: records}
}

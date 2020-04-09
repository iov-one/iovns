package domain

import (
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
	// TODO validate genesis by: checking no duplicate names, and domain validity
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

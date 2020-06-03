package keeper

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/types"
)

func Test_queryGetAccountsInDomain(t *testing.T) {
	testCases := map[string]subTest{
		"success default": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				ctx = ctx.WithBlockTime(time.Unix(0, 0))
				k.CreateDomain(ctx, types.Domain{Name: "test", Admin: aliceAddr})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "1", Owner: aliceAddr})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "2", Owner: aliceAddr})
			},
			Request: &QueryAccountsInDomain{
				Domain:         "test",
				ResultsPerPage: 0,
				Offset:         0,
			},
			Handler: queryAccountsInDomainHandler,
			WantErr: nil,
			PtrExpectedResponse: QueryAccountsInDomainResponse{
				Accounts: []types.Account{{Domain: "test", Name: "", Owner: aliceAddr, ValidUntil: types.MaxValidUntil}, {Domain: "test", Name: "1", Owner: aliceAddr}, {Domain: "test", Name: "2", Owner: aliceAddr}},
			},
		},
		"success with paging": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				k.CreateDomain(ctx, types.Domain{Name: "test", Admin: aliceAddr})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "1", Owner: aliceAddr})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "2", Owner: bobAddr})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "3", Owner: aliceAddr})
			},
			Request: &QueryAccountsInDomain{
				Domain:         "test",
				ResultsPerPage: 1,
				Offset:         2,
			},
			Handler: queryAccountsInDomainHandler,
			WantErr: nil,
			PtrExpectedResponse: QueryAccountsInDomainResponse{
				Accounts: []types.Account{{Domain: "test", Name: "1", Owner: aliceAddr}},
			},
		},
	}

	runQueryTests(t, testCases)
}

func Test_queryGetAccountsFromOwner(t *testing.T) {

	testCases := map[string]subTest{
		"success": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				k.CreateDomain(ctx, types.Domain{Name: "test", Admin: bobAddr})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "1", Owner: aliceAddr})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "2", Owner: aliceAddr})
			},
			Request: &QueryAccountsFromOwner{
				Owner:          aliceAddr,
				ResultsPerPage: 0,
				Offset:         0,
			},
			Handler: queryAccountsFromOwnerHandler,
			WantErr: nil,
			PtrExpectedResponse: &QueryAccountsFromOwnerResponse{
				Accounts: []types.Account{
					{Domain: "test", Name: "1", Owner: aliceAddr},
					{Domain: "test", Name: "2", Owner: aliceAddr},
				},
			},
		},
	}
	runQueryTests(t, testCases)
}

func Test_queryGetDomainsFromOwner(t *testing.T) {
	testCases := map[string]subTest{
		"success": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				k.CreateDomain(ctx, types.Domain{Name: "test", Admin: aliceAddr})
				k.CreateDomain(ctx, types.Domain{Name: "test2", Admin: aliceAddr})
			},
			Request: &QueryDomainsFromOwner{
				Owner:          aliceAddr,
				ResultsPerPage: 0,
				Offset:         0,
			},
			Handler: queryDomainsFromOwnerHandler,
			WantErr: nil,
			PtrExpectedResponse: &QueryDomainsFromOwnerResponse{
				Domains: []types.Domain{
					{Name: "test", Admin: aliceAddr},
					{Name: "test2", Admin: aliceAddr},
				},
			},
		},
	}
	runQueryTests(t, testCases)
}

func Test_queryResolveAccountHandler(t *testing.T) {
	testCases := map[string]subTest{
		"success": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				k.CreateAccount(ctx, types.Account{
					Domain: "test",
					Name:   "1",
					Owner:  bobAddr,
				})
			},
			Request: &QueryResolveAccount{
				Domain: "test",
				Name:   "1",
			},
			Handler: queryResolveAccountHandler,
			WantErr: nil,
			PtrExpectedResponse: &QueryResolveAccountResponse{Account: types.Account{
				Domain: "test",
				Name:   "1",
				Owner:  bobAddr,
			}},
		},
	}
	runQueryTests(t, testCases)
}

func Test_queryResolveDomainHandler(t *testing.T) {
	testCases := map[string]subTest{
		"success": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				k.CreateDomain(ctx, types.Domain{
					Name:  "test",
					Admin: bobAddr,
				})
			},
			Request:             &QueryResolveDomain{Name: "test"},
			Handler:             queryResolveDomainHandler,
			WantErr:             nil,
			PtrExpectedResponse: &QueryResolveDomainResponse{Domain: types.Domain{Name: "test", Admin: bobAddr}},
		},
	}

	runQueryTests(t, testCases)
}

func Test_queryTargetAccountsHandler(t *testing.T) {
	testCases := map[string]subTest{
		"success": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				target := types.BlockchainAddress{
					ID:      "id-1",
					Address: "addr-1",
				}
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "1",
					Owner:      bobAddr,
					ValidUntil: 0,
					Targets:    []types.BlockchainAddress{target},
				})
				k.CreateAccount(ctx, types.Account{
					Domain:     "test",
					Name:       "2",
					Owner:      bobAddr,
					ValidUntil: 0,
					Targets:    []types.BlockchainAddress{target},
				})
			},
			Request: &QueryTargetAccounts{
				Target: types.BlockchainAddress{
					ID:      "id-1",
					Address: "addr-1",
				},
			},
			Handler: queryTargetAccountsHandler,
			WantErr: nil,
			PtrExpectedResponse: &QueryTargetAccountsResponse{
				Accounts: []types.Account{
					{
						Domain:     "test",
						Name:       "1",
						Owner:      bobAddr,
						ValidUntil: 0,
						Targets: []types.BlockchainAddress{{
							ID:      "id-1",
							Address: "addr-1",
						}},
					},
					{
						Domain:     "test",
						Name:       "2",
						Owner:      bobAddr,
						ValidUntil: 0,
						Targets: []types.BlockchainAddress{{
							ID:      "id-1",
							Address: "addr-1",
						}},
					},
				}},
		},
	}

	runQueryTests(t, testCases)
}

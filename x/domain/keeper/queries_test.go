package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/iov-one/iovns/x/domain/types"
)

func Test_queryGetAccountsInDomain(t *testing.T) {
	testCases := map[string]subTest{
		"success default": {
			BeforeTest: func(t *testing.T, ctx types.Context, k Keeper) {
				k.CreateDomain(ctx, types2.Domain{Name: "test"})
				k.CreateAccount(ctx, types2.Account{Domain: "test", Name: "1"})
				k.CreateAccount(ctx, types2.Account{Domain: "test", Name: "2"})
			},
			Request: &QueryAccountsInDomain{
				Domain:         "test",
				ResultsPerPage: 0,
				Offset:         0,
			},
			Handler: queryAccountsInDomainHandler,
			WantErr: nil,
			PtrExpectedResponse: QueryAccountsInDomainResponse{
				Accounts: []types2.Account{{Domain: "test", Name: "1"}, {Domain: "test", Name: "2"}},
			},
		},
		"success with paging": {
			BeforeTest: func(t *testing.T, ctx types.Context, k Keeper) {
				k.CreateDomain(ctx, types2.Domain{Name: "test"})
				k.CreateAccount(ctx, types2.Account{Domain: "test", Name: "1"})
				k.CreateAccount(ctx, types2.Account{Domain: "test", Name: "2"})
				k.CreateAccount(ctx, types2.Account{Domain: "test", Name: "3"})
			},
			Request: &QueryAccountsInDomain{
				Domain:         "test",
				ResultsPerPage: 1,
				Offset:         2,
			},
			Handler: queryAccountsInDomainHandler,
			WantErr: nil,
			PtrExpectedResponse: QueryAccountsInDomainResponse{
				Accounts: []types2.Account{{Domain: "test", Name: "2"}},
			},
		},
	}

	runQueryTests(t, testCases)
}

func Test_queryGetAccountsFromOwner(t *testing.T) {

	testCases := map[string]subTest{
		"success": {
			BeforeTest: func(t *testing.T, ctx types.Context, k Keeper) {
				k.CreateDomain(ctx, types2.Domain{Name: "test"})
				k.CreateAccount(ctx, types2.Account{Domain: "test", Name: "1", Owner: aliceAddr})
				k.CreateAccount(ctx, types2.Account{Domain: "test", Name: "2", Owner: aliceAddr})
			},
			Request: &QueryAccountsFromOwner{
				Owner:          aliceAddr,
				ResultsPerPage: 0,
				Offset:         0,
			},
			Handler: queryAccountsFromOwnerHandler,
			WantErr: nil,
			PtrExpectedResponse: &QueryAccountsFromOwnerResponse{
				Accounts: []types2.Account{
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
			BeforeTest: func(t *testing.T, ctx types.Context, k Keeper) {
				k.CreateDomain(ctx, types2.Domain{Name: "test", Admin: aliceAddr})
				k.CreateDomain(ctx, types2.Domain{Name: "test2", Admin: aliceAddr})
			},
			Request: &QueryDomainsFromOwner{
				Owner:          aliceAddr,
				ResultsPerPage: 0,
				Offset:         0,
			},
			Handler: queryDomainsFromOwnerHandler,
			WantErr: nil,
			PtrExpectedResponse: &QueryDomainsFromOwnerResponse{
				Domains: []types2.Domain{
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
			BeforeTest: func(t *testing.T, ctx types.Context, k Keeper) {
				k.CreateAccount(ctx, types2.Account{
					Domain: "test",
					Name:   "1",
				})
			},
			Request: &QueryResolveAccount{
				Domain: "test",
				Name:   "1",
			},
			Handler: queryResolveAccountHandler,
			WantErr: nil,
			PtrExpectedResponse: &QueryResolveAccountResponse{Account: types2.Account{
				Domain: "test",
				Name:   "1",
			}},
		},
	}
	runQueryTests(t, testCases)
}

func Test_queryResolveDomainHandler(t *testing.T) {
	testCases := map[string]subTest{
		"success": {
			BeforeTest: func(t *testing.T, ctx types.Context, k Keeper) {
				k.CreateDomain(ctx, types2.Domain{
					Name: "test",
				})
			},
			Request:             &QueryResolveDomain{Name: "test"},
			Handler:             queryResolveDomainHandler,
			WantErr:             nil,
			PtrExpectedResponse: &QueryResolveDomainResponse{Domain: types2.Domain{Name: "test"}},
		},
	}

	runQueryTests(t, testCases)
}

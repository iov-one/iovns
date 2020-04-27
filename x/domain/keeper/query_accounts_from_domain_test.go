package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/types"
	"testing"
)

func Test_queryGetAccountsInDomain(t *testing.T) {
	testCases := map[string]subTest{
		"success default": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				k.CreateDomain(ctx, types.Domain{Name: "test"})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "1"})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "2"})
			},
			Request: &QueryAccountsInDomain{
				Domain:         "test",
				ResultsPerPage: 0,
				Offset:         0,
			},
			Handler: queryAccountsInDomainHandler,
			WantErr: nil,
			PtrExpectedResponse: QueryAccountsInDomainResponse{
				Accounts: []types.Account{{Domain: "test", Name: "1"}, {Domain: "test", Name: "2"}},
			},
		},
		"success with paging": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				k.CreateDomain(ctx, types.Domain{Name: "test"})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "1"})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "2"})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "3"})
			},
			Request: &QueryAccountsInDomain{
				Domain:         "test",
				ResultsPerPage: 1,
				Offset:         2,
			},
			Handler: queryAccountsInDomainHandler,
			WantErr: nil,
			PtrExpectedResponse: QueryAccountsInDomainResponse{
				Accounts: []types.Account{{Domain: "test", Name: "2"}},
			},
		},
	}

	runQueryTests(t, testCases)
}

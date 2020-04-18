package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/types"
	"testing"
)

func Test_queryGetAccountsFromOwner(t *testing.T) {

	testCases := map[string]subTest{
		"success": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				k.SetDomain(ctx, types.Domain{Name: "test"})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "1", Owner: aliceAddr})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "2", Owner: aliceAddr})
			},
			Request: &QueryAccountsFromOwner{
				Owner:          aliceAddr,
				ResultsPerPage: 0,
				Offset:         0,
			},
			Handler: queryGetAccountsFromOwner,
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

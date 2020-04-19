package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/types"
	"testing"
)

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

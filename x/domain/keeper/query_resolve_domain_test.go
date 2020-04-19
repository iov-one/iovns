package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/types"
	"testing"
)

func Test_queryResolveDomainHandler(t *testing.T) {
	testCases := map[string]subTest{
		"success": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				k.CreateDomain(ctx, types.Domain{
					Name: "test",
				})
			},
			Request:             &QueryResolveDomain{Name: "test"},
			Handler:             queryResolveDomainHandler,
			WantErr:             nil,
			PtrExpectedResponse: &QueryResolveDomainResponse{Domain: types.Domain{Name: "test"}},
		},
	}

	runQueryTests(t, testCases)
}

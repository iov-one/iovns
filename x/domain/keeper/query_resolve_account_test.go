package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/types"
	"testing"
)

func Test_queryResolveAccountHandler(t *testing.T) {
	testCases := map[string]subTest{
		"success": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				k.CreateAccount(ctx, types.Account{
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
			PtrExpectedResponse: &QueryResolveAccountResponse{Account: types.Account{
				Domain: "test",
				Name:   "1",
			}},
		},
	}
	runQueryTests(t, testCases)
}

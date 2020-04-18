package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/domain/types"
	"testing"
)

func Test_queryGetAccountsInDomain(t *testing.T) {
	testCases := map[string]subTest{
		"domain name empty": {
			Request: &QueryAccountsInDomain{
				Domain:         "",
				ResultsPerPage: 0,
				Offset:         0,
			},
			Handler:             queryGetAccountsInDomain,
			WantErr:             types.ErrInvalidDomainName,
			PtrExpectedResponse: nil,
		},
		"no accounts": {
			Request: &QueryAccountsInDomain{
				Domain:         "does not exist",
				ResultsPerPage: 0,
				Offset:         0,
			},
			Handler:             queryGetAccountsInDomain,
			WantErr:             nil,
			PtrExpectedResponse: QueryAccountsInDomainResponse{},
		},
		"success default": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				k.SetDomain(ctx, types.Domain{Name: "test"})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "1"})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "2"})
			},
			Request: &QueryAccountsInDomain{
				Domain:         "test",
				ResultsPerPage: 0,
				Offset:         0,
			},
			Handler: queryGetAccountsInDomain,
			WantErr: nil,
			PtrExpectedResponse: QueryAccountsInDomainResponse{
				Accounts: []types.Account{{Domain: "test", Name: "1"}, {Domain: "test", Name: "2"}},
			},
		},
		"success with paging": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				k.SetDomain(ctx, types.Domain{Name: "test"})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "1"})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "2"})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "3"})
			},
			Request: &QueryAccountsInDomain{
				Domain:         "test",
				ResultsPerPage: 1,
				Offset:         2,
			},
			Handler: queryGetAccountsInDomain,
			WantErr: nil,
			PtrExpectedResponse: QueryAccountsInDomainResponse{
				Accounts: []types.Account{{Domain: "test", Name: "2"}},
			},
		},
		"invalid paging": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				k.SetDomain(ctx, types.Domain{Name: "test"})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "1"})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "2"})
				k.CreateAccount(ctx, types.Account{Domain: "test", Name: "3"})
			},
			Request: &QueryAccountsInDomain{
				Domain:         "test",
				ResultsPerPage: 1,
				Offset:         4,
			},
			Handler:             queryGetAccountsInDomain,
			WantErr:             sdkerrors.ErrInvalidRequest,
			PtrExpectedResponse: nil,
		},
	}

	runQueryTests(t, testCases)
}

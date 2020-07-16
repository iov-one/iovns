package keeper

import (
	"github.com/iov-one/iovns/tutils"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/starname/types"
)

func Test_queryGetAccountsInDomain(t *testing.T) {
	testCases := map[string]subTest{
		"success default": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				ctx = ctx.WithBlockTime(time.Unix(0, 0))
				ds := k.DomainStore(ctx)
				ds.Create(&types.Domain{Name: "test", Admin: aliceAddr})
				as := k.AccountStore(ctx)
				as.Create(&types.Account{Domain: "test", Name: tutils.StrPtr("1"), Owner: aliceAddr})
				as.Create(&types.Account{Domain: "test", Name: tutils.StrPtr("2"), Owner: aliceAddr})
			},
			Request: &QueryAccountsInDomain{
				Domain:         "test",
				ResultsPerPage: 0,
				Offset:         0,
			},
			Handler: queryAccountsInDomainHandler,
			WantErr: nil,
			PtrExpectedResponse: QueryAccountsInDomainResponse{
				Accounts: []*types.Account{{Domain: "test", Name: tutils.StrPtr("1"), Owner: aliceAddr}, {Domain: "test", Name: tutils.StrPtr("2"), Owner: aliceAddr}},
			},
		},
		"success with paging": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				as := k.AccountStore(ctx)
				ds := k.DomainStore(ctx)
				ds.Create(&types.Domain{Name: "test", Admin: aliceAddr})
				as.Create(&types.Account{Domain: "test", Name: tutils.StrPtr("1"), Owner: aliceAddr})
				as.Create(&types.Account{Domain: "test", Name: tutils.StrPtr("2"), Owner: bobAddr})
				as.Create(&types.Account{Domain: "test", Name: tutils.StrPtr("3"), Owner: aliceAddr})
			},
			Request: &QueryAccountsInDomain{
				Domain:         "test",
				ResultsPerPage: 1,
				Offset:         2,
			},
			Handler: queryAccountsInDomainHandler,
			WantErr: nil,
			PtrExpectedResponse: QueryAccountsInDomainResponse{
				Accounts: []*types.Account{{Domain: "test", Name: tutils.StrPtr("2"), Owner: bobAddr}},
			},
		},
	}

	runQueryTests(t, testCases)
}

func Test_queryGetAccountsWithOwner(t *testing.T) {

	testCases := map[string]subTest{
		"success": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				as := k.AccountStore(ctx)
				ds := k.DomainStore(ctx)
				ds.Create(&types.Domain{Name: "test", Admin: bobAddr})
				as.Create(&types.Account{Domain: "test", Name: tutils.StrPtr("1"), Owner: aliceAddr})
				as.Create(&types.Account{Domain: "test", Name: tutils.StrPtr("2"), Owner: aliceAddr})
				as.Create(&types.Account{Domain: "test", Name: tutils.StrPtr("3"), Owner: aliceAddr})
				as.Create(&types.Account{Domain: "test", Name: tutils.StrPtr("4"), Owner: aliceAddr})

			},
			Request: &QueryAccountsWithOwner{
				Owner:          aliceAddr,
				ResultsPerPage: 0,
				Offset:         0,
			},
			Handler: queryAccountsWithOwnerHandler,
			WantErr: nil,
			PtrExpectedResponse: &QueryAccountsWithOwnerResponse{
				Accounts: []types.Account{
					{Domain: "test", Name: tutils.StrPtr("1"), Owner: aliceAddr},
					{Domain: "test", Name: tutils.StrPtr("2"), Owner: aliceAddr},
					{Domain: "test", Name: tutils.StrPtr("3"), Owner: aliceAddr},
					{Domain: "test", Name: tutils.StrPtr("4"), Owner: aliceAddr},
				},
			},
		},
	}
	runQueryTests(t, testCases)
}

func Test_queryGetDomainsWithOwner(t *testing.T) {
	testCases := map[string]subTest{
		"success": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				ds := k.DomainStore(ctx)
				ds.Create(&types.Domain{Name: "test", Admin: aliceAddr})
				ds.Create(&types.Domain{Name: "test2", Admin: aliceAddr})
			},
			Request: &QueryDomainsWithOwner{
				Owner:          aliceAddr,
				ResultsPerPage: 0,
				Offset:         0,
			},
			Handler: queryDomainsWithOwnerHandler,
			WantErr: nil,
			PtrExpectedResponse: &QueryDomainsWithOwnerResponse{
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
				as := k.AccountStore(ctx)
				as.Create(&types.Account{
					Domain: "test",
					Name:   tutils.StrPtr("1"),
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
				Name:   tutils.StrPtr("1"),
				Owner:  bobAddr,
			}},
		},
		"success starname": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				k.AccountStore(ctx).Create(&types.Account{
					Domain: "test",
					Name:   tutils.StrPtr("1"),
					Owner:  bobAddr,
				})
			},
			Request: &QueryResolveAccount{
				Starname: "1*test",
				Name:     "",
				Domain:   "",
			},
			Handler: queryResolveAccountHandler,
			WantErr: nil,
			PtrExpectedResponse: &QueryResolveAccountResponse{Account: types.Account{
				Domain: "test",
				Name:   tutils.StrPtr("1"),
				Owner:  bobAddr,
			}},
		},
		"failure provide only one param starname": {
			Request: &QueryResolveAccount{
				Domain:   "test",
				Name:     "1",
				Starname: "1*test",
			},
			Handler: queryResolveAccountHandler,
			WantErr: types.ErrProvideStarnameOrDomainName,
		},
		"failure provide only one param starname 2": {
			Request: &QueryResolveAccount{
				Domain:   "test",
				Name:     "",
				Starname: "1*test",
			},
			Handler: queryResolveAccountHandler,
			WantErr: types.ErrProvideStarnameOrDomainName,
		},
		"failure provide only one param starname 3": {
			Request: &QueryResolveAccount{
				Name:     "test",
				Starname: "1*test",
			},
			Handler: queryResolveAccountHandler,
			WantErr: types.ErrProvideStarnameOrDomainName,
		},
		"starname must contain separator": {
			Request: &QueryResolveAccount{
				Starname: "1test",
			},
			Handler: queryResolveAccountHandler,
			WantErr: types.ErrStarnameNotContainSep,
		},
		"starname must contain single separator": {
			Request: &QueryResolveAccount{
				Starname: "1*te*st",
			},
			Handler: queryResolveAccountHandler,
			WantErr: types.ErrStarnameMultipleSeparator,
		},
	}
	runQueryTests(t, testCases)
}

func Test_queryResolveDomainHandler(t *testing.T) {
	testCases := map[string]subTest{
		"success": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				k.DomainStore(ctx).Create(&types.Domain{
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

func Test_queryResourceAccountsHandler(t *testing.T) {
	testCases := map[string]subTest{
		"success": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				as := k.AccountStore(ctx)
				resource := types.Resource{
					URI:      "id-1",
					Resource: "addr-1",
				}
				as.Create(&types.Account{
					Domain:     "test",
					Name:       tutils.StrPtr("1"),
					Owner:      bobAddr,
					ValidUntil: 0,
					Resources:  []types.Resource{resource},
				})
				as.Create(&types.Account{
					Domain:     "test",
					Name:       tutils.StrPtr("2"),
					Owner:      bobAddr,
					ValidUntil: 0,
					Resources:  []types.Resource{resource},
				})
			},
			Request: &QueryResolveResource{
				Resource: types.Resource{
					URI:      "id-1",
					Resource: "addr-1",
				},
			},
			Handler: queryResourceAccountHandler,
			WantErr: nil,
			PtrExpectedResponse: &QueryResolveResourceResponse{
				Accounts: []types.Account{
					{
						Domain:     "test",
						Name:       tutils.StrPtr("1"),
						Owner:      bobAddr,
						ValidUntil: 0,
						Resources: []types.Resource{{
							URI:      "id-1",
							Resource: "addr-1",
						}},
					},
					{
						Domain:     "test",
						Name:       tutils.StrPtr("2"),
						Owner:      bobAddr,
						ValidUntil: 0,
						Resources: []types.Resource{{
							URI:      "id-1",
							Resource: "addr-1",
						}},
					},
				}},
		},
	}

	runQueryTests(t, testCases)
}

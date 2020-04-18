package keeper

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/domain/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

type subTest struct {
	// BeforeTest are the action to perform before the test
	BeforeTest func(t *testing.T, ctx sdk.Context, k Keeper)
	// Request is the query request
	Request interface{}
	// Handler is the handler function of the query
	Handler func(ctx sdk.Context, args []string, req abci.RequestQuery, k Keeper) ([]byte, error)
	// WantErr is the error we expect, if != from nil it will be matched with errors.Is
	WantErr error
	// ExpectedResponse is the response we want that will be marshalled and checked agains the response
	ExpectedResponse interface{}
}

func Test_queryGetAccountsInDomain(t *testing.T) {
	testCases := map[string]subTest{
		"domain name empty": {
			Request: &QueryAccountsInDomain{
				Domain:         "",
				ResultsPerPage: 0,
				Offset:         0,
			},
			Handler:          queryGetAccountsInDomain,
			WantErr:          types.ErrInvalidDomainName,
			ExpectedResponse: nil,
		},
		"no accounts": {
			Request: &QueryAccountsInDomain{
				Domain:         "does not exist",
				ResultsPerPage: 0,
				Offset:         0,
			},
			Handler:          queryGetAccountsInDomain,
			WantErr:          nil,
			ExpectedResponse: QueryAccountsInDomainResponse{},
		},
		"success default": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				k.SetDomain(ctx, types.Domain{Name: "test"})
				k.SetAccount(ctx, types.Account{Domain: "test", Name: "1"})
				k.SetAccount(ctx, types.Account{Domain: "test", Name: "2"})
			},
			Request: &QueryAccountsInDomain{
				Domain:         "test",
				ResultsPerPage: 0,
				Offset:         0,
			},
			Handler: queryGetAccountsInDomain,
			WantErr: nil,
			ExpectedResponse: QueryAccountsInDomainResponse{
				Accounts: []types.Account{{Domain: "test", Name: "1"}, {Domain: "test", Name: "2"}},
			},
		},
		"success with paging": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				k.SetDomain(ctx, types.Domain{Name: "test"})
				k.SetAccount(ctx, types.Account{Domain: "test", Name: "1"})
				k.SetAccount(ctx, types.Account{Domain: "test", Name: "2"})
				k.SetAccount(ctx, types.Account{Domain: "test", Name: "3"})
			},
			Request: &QueryAccountsInDomain{
				Domain:         "test",
				ResultsPerPage: 1,
				Offset:         2,
			},
			Handler: queryGetAccountsInDomain,
			WantErr: nil,
			ExpectedResponse: QueryAccountsInDomainResponse{
				Accounts: []types.Account{{Domain: "test", Name: "2"}},
			},
		},
		"invalid paging": {
			BeforeTest: func(t *testing.T, ctx sdk.Context, k Keeper) {
				k.SetDomain(ctx, types.Domain{Name: "test"})
				k.SetAccount(ctx, types.Account{Domain: "test", Name: "1"})
				k.SetAccount(ctx, types.Account{Domain: "test", Name: "2"})
				k.SetAccount(ctx, types.Account{Domain: "test", Name: "3"})
			},
			Request: &QueryAccountsInDomain{
				Domain:         "test",
				ResultsPerPage: 1,
				Offset:         4,
			},
			Handler:          queryGetAccountsInDomain,
			WantErr:          sdkerrors.ErrInvalidRequest,
			ExpectedResponse: nil,
		},
	}

	runQueryTests(t, testCases)
}

func runQueryTests(t *testing.T, cases map[string]subTest) {
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			k, ctx := NewTestKeeper(t, true)
			if test.BeforeTest != nil {
				test.BeforeTest(t, ctx, k)
			}
			// marshal request
			reqBody, err := json.Marshal(test.Request)
			if err != nil {
				t.Fatalf("unable to marshal request: %s", err)
			}
			// do test
			got, err := test.Handler(ctx, nil, abci.RequestQuery{Data: reqBody}, k)
			if !errors.Is(err, test.WantErr) {
				t.Fatalf("wanted err: %s, got: %s", test.WantErr, err)
			}
			// check if expected response should be nil to avoid
			// false positives of marshaling as "null"
			if got == nil && test.ExpectedResponse == nil {
				// success
				return
			}
			// marshal expected response and compare with what we've got
			expectedBytes, err := codec.MarshalJSONIndent(k.cdc, test.ExpectedResponse)
			if err != nil {
				t.Fatalf("marshal error: %s", err)
			}
			if !bytes.Equal(got, expectedBytes) {
				t.Fatalf("unexpected response: \nwant:\t %s, \ngot:\t %s", expectedBytes, got)
			}
		})
	}
}

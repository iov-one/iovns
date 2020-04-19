package keeper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"os"
	"testing"
)

var aliceAddr, bobAddr sdk.AccAddress

func TestMain(t *testing.M) {
	aliceAddr, bobAddr = genTestAddress()
	os.Exit(t.Run())
}

func genTestAddress() (sdk.AccAddress, sdk.AccAddress) {
	keyBase := keys.NewInMemory()
	addr1, _, err := keyBase.CreateMnemonic("alice", keys.English, "", keys.Secp256k1)
	if err != nil {
		fmt.Println("unable to generate mock addresses " + err.Error())
		os.Exit(1)
	}
	addr2, _, err := keyBase.CreateMnemonic("bob", keys.English, "", keys.Secp256k1)
	if err != nil {
		fmt.Println("unable to generate mock addresses " + err.Error())
		os.Exit(1)
	}
	return addr1.GetAddress(), addr2.GetAddress()
}

type subTest struct {
	// BeforeTest are the action to perform before the test
	BeforeTest func(t *testing.T, ctx sdk.Context, k Keeper)
	// Request is the query request
	Request interface{ Validate() error } // represents aliceAddr query
	// Handler is the handler function of the query
	Handler func(ctx sdk.Context, args []string, req abci.RequestQuery, k Keeper) ([]byte, error)
	// WantErr is the error we expect, if != from nil it will be matched with errors.Is
	WantErr error
	// PtrExpectedResponse is the response we want that will be marshalled and checked against the response (pointer expected)
	PtrExpectedResponse interface{}
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
			if got == nil && test.PtrExpectedResponse == nil {
				// success
				return
			}
			// marshal expected response and compare with what we've got
			expectedBytes, err := codec.MarshalJSONIndent(k.cdc, test.PtrExpectedResponse)
			if err != nil {
				t.Fatalf("marshal error: %s", err)
			}
			if !bytes.Equal(got, expectedBytes) {
				t.Fatalf("unexpected response: \nwant:\t %s, \ngot:\t %s", expectedBytes, got)
			}
		})
	}
}

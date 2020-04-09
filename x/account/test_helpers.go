package account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
)

type subTest struct {
	// BeforeTest is the function run before doing the test,
	// used for example to store state, like configurations etc.
	// Ignored if nil
	BeforeTest func(t *testing.T, k Keeper, ctx sdk.Context)
	// Test is the function that runs the actual test
	Test func(t *testing.T, k Keeper, ctx sdk.Context)
	// AfterTest performs actions after the test is run, it can
	// be used to check if the state after Test is run matches
	// the result we expect.
	// Ignored if nil
	AfterTest func(t *testing.T, k Keeper, ctx sdk.Context)
}

/*
// runTests run tests cases after generating a new keeper and context for each test case
func runTests(t *testing.T, tests map[string]subTest) {
	for name, test := range tests {
		keeper, ctx := newTestKeeper(t, true)
		// run sub subTest
		t.Run(name, func(t *testing.T) {
			// run before subTest
			if test.BeforeTest != nil {
				test.BeforeTest(t, keeper, ctx)
			}
			// run actual subTest
			test.Test(t, keeper, ctx)
			// run after subTest
			if test.AfterTest != nil {
				test.AfterTest(t, keeper, ctx)
			}
		})
	}
}


*/
func newTestKeeper() {

}

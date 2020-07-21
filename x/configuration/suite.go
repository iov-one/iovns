package configuration

import (
	"testing"
	"time"

	"github.com/iov-one/iovns/pkg/utils"

	"github.com/cosmos/cosmos-sdk/types"
)

var ks, addrs = utils.GeneratePrivKeyAddressPairs(3)
var AliceKey types.AccAddress = addrs[0]
var BobKey types.AccAddress = addrs[1]
var CharlieKey types.AccAddress = addrs[2]

const RegexMatchAll = "^(.*?)?"
const RegexMatchNothing = "$^"

// subTest defines a test runner
type SubTest struct {
	// BeforeTestBlockTime is the block time during before test in unix seconds
	// WARNING: if block time is given 0, it will be accepted as time.Now()
	BeforeTestBlockTime int64
	// BeforeTest is the function run before doing the test,
	// used for example to store state, like configurations etc.
	// Ignored if nil
	BeforeTest func(t *testing.T, k Keeper, ctx types.Context)
	// TestBlockTime is the block time during test in unix seconds
	// WARNING: if block time is given 0, it will be accepted as time.Now()
	TestBlockTime int64
	// Test is the function that runs the actual test
	Test func(t *testing.T, k Keeper, ctx types.Context)
	// AfterTestBlockTime is the block time during after test in unix seconds
	// WARNING: if block time is given 0, it will be accepted as time.Now()
	AfterTestBlockTime int64
	// AfterTest performs actions after the test is run, it can
	// be used to check if the state after Test is run matches
	// the result we expect.
	// Ignored if nil
	AfterTest func(t *testing.T, k Keeper, ctx types.Context)
}

// runTests run tests cases after generating a new keeper and context for each test case
func RunTests(t *testing.T, tests map[string]SubTest) {
	for name, test := range tests {
		keeper, ctx := NewTestKeeper(t, true)
		// set default fees
		setFees := keeper.SetFees
		fees := NewFees()
		fees.SetDefaults("testcoin")
		setFees(ctx, fees)
		// run sub SubTest
		t.Run(name, func(t *testing.T) {
			// run before SubTest
			if test.BeforeTest != nil {
				if test.BeforeTestBlockTime != 0 {
					t := time.Unix(test.BeforeTestBlockTime, 0)
					ctx = ctx.WithBlockTime(t)
				}
				test.BeforeTest(t, keeper, ctx)
			}

			if test.TestBlockTime != 0 {
				t := time.Unix(test.TestBlockTime, 0)
				ctx = ctx.WithBlockTime(t)
			}
			// run actual SubTest
			test.Test(t, keeper, ctx)

			// run after SubTest
			if test.AfterTest != nil {
				if test.AfterTestBlockTime != 0 {
					t := time.Unix(test.AfterTestBlockTime, 0)
					ctx = ctx.WithBlockTime(t)
				}
				test.AfterTest(t, keeper, ctx)
			}
		})
	}
}

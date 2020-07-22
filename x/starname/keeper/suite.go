package keeper

import (
	"fmt"
	"testing"
	"time"

	"github.com/iov-one/iovns/pkg/utils"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration"
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
	BeforeTest func(t *testing.T, k Keeper, ctx types.Context, mocks *Mocks)
	// TestBlockTime is the block time during test in unix seconds
	// WARNING: if block time is given 0, it will be accepted as time.Now()
	TestBlockTime int64
	// Test is the function that runs the actual test
	Test func(t *testing.T, k Keeper, ctx types.Context, mocks *Mocks)
	// AfterTestBlockTime is the block time during after test in unix seconds
	// WARNING: if block time is given 0, it will be accepted as time.Now()
	AfterTestBlockTime int64
	// AfterTest performs actions after the test is run, it can
	// be used to check if the state after Test is run matches
	// the result we expect.
	// Ignored if nil
	AfterTest func(t *testing.T, k Keeper, ctx types.Context, mocks *Mocks)
}

// runTests run tests cases after generating a new keeper and context for each test case
func RunTests(t *testing.T, tests map[string]SubTest) {
	for name, test := range tests {
		domainKeeper, ctx, mocks := NewTestKeeper(t, true)
		// set default mock.Supply not to fail
		mocks.Supply.SetSendCoinsFromAccountToModule(func(ctx types.Context, addr types.AccAddress, moduleName string, coins types.Coins) error {
			return nil
		})
		// set default fees
		setFees := domainKeeper.ConfigurationKeeper.(ConfigurationSetter).SetFees
		fees := configuration.NewFees()
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
				test.BeforeTest(t, domainKeeper, ctx, mocks)
			}

			if test.TestBlockTime != 0 {
				t := time.Unix(test.TestBlockTime, 0)
				ctx = ctx.WithBlockTime(t)
			}
			// run actual SubTest
			test.Test(t, domainKeeper, ctx, mocks)

			// run after SubTest
			if test.AfterTest != nil {
				if test.AfterTestBlockTime != 0 {
					t := time.Unix(test.AfterTestBlockTime, 0)
					ctx = ctx.WithBlockTime(t)
				}
				test.AfterTest(t, domainKeeper, ctx, mocks)
			}
		})
	}
}

// since the exposed interface for configuration keeper
// does not include set config as the domain module should
// not be able to change configuration state, then only
// in test cases we expose this method
type ConfigurationSetter interface {
	SetConfig(ctx types.Context, config configuration.Config)
	SetFees(ctx types.Context, fees *configuration.Fees)
}

// getConfigSetter exposes the configurationSetter interface
// allowing the module to set configuration state, this should only
// be used for tests and will panic if the keeper provided can not
// be cast to configurationSetter
func GetConfigSetter(keeper ConfigurationKeeper) ConfigurationSetter {
	// check if the configuration keeper is also a config setter
	configSetter, ok := keeper.(ConfigurationSetter)
	if !ok {
		panic(fmt.Sprintf("cannot cast configuration keeper to configuration setter: got uncastable type: %T", keeper))
	}
	return configSetter
}

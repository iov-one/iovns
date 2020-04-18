package domain

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/keeper"
	"os"
	"testing"
)

var aliceKey keys.Info
var bobKey keys.Info

const regexMatchAll = "^(.*?)?"
const regexMatchNothing = "$^"

// TestMain is going to init test addresses
func TestMain(t *testing.M) {
	keyBase := keys.NewInMemory()
	addr1, _, err := keyBase.CreateMnemonic("alice", keys.English, "", keys.Secp256k1)
	if err != nil {
		fmt.Println("unable to generate mock addresses " + err.Error())
		os.Exit(1)
	}
	aliceKey = addr1
	addr2, _, err := keyBase.CreateMnemonic("bob", keys.English, "", keys.Secp256k1)
	if err != nil {
		fmt.Println("unable to generate mock addresses " + err.Error())
		os.Exit(1)
	}
	bobKey = addr2
	// run and exit
	os.Exit(t.Run())
}

// subTest defines a test runner
type subTest struct {
	// BeforeTest is the function run before doing the test,
	// used for example to store state, like configurations etc.
	// Ignored if nil
	BeforeTest func(t *testing.T, k keeper.Keeper, ctx sdk.Context)
	// Test is the function that runs the actual test
	Test func(t *testing.T, k keeper.Keeper, ctx sdk.Context)
	// AfterTest performs actions after the test is run, it can
	// be used to check if the state after Test is run matches
	// the result we expect.
	// Ignored if nil
	AfterTest func(t *testing.T, k keeper.Keeper, ctx sdk.Context)
}

// runTests run tests cases after generating a new keeper and context for each test case
func runTests(t *testing.T, tests map[string]subTest) {
	for name, test := range tests {
		domainKeeper, ctx := keeper.NewTestKeeper(t, true)
		// run sub subTest
		t.Run(name, func(t *testing.T) {
			// run before subTest
			if test.BeforeTest != nil {
				test.BeforeTest(t, domainKeeper, ctx)
			}
			// run actual subTest
			test.Test(t, domainKeeper, ctx)
			// run after subTest
			if test.AfterTest != nil {
				test.AfterTest(t, domainKeeper, ctx)
			}
		})
	}
}

// since the exposed interface for configuration keeper
// does not include set config as the domain module should
// not be able to change configuration state, then only
// in test cases we expose this method
type configurationSetter interface {
	SetConfig(ctx sdk.Context, config configuration.Config)
}

// getConfigSetter exposes the configurationSetter interface
// allowing the module to set configuration state, this should only
// be used for tests and will panic if the keeper provided can not
// be cast to configurationSetter
func getConfigSetter(keeper keeper.ConfigurationKeeper) configurationSetter {
	// check if the configuration keeper is also a config setter
	configSetter, ok := keeper.(configurationSetter)
	if !ok {
		panic(fmt.Sprintf("cannot cast configuration keeper to configuration setter: got uncastable type: %T", keeper))
	}
	return configSetter
}

package tutils

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/keeper"
)

// SubTest defines a test runner
type SubTest struct {
	// BeforeTestBlockTime is the block time during before test in unix seconds
	BeforeTestBlockTime int64
	// BeforeTest is the function run before doing the test,
	// used for example to store state, like configurations etc.
	// Ignored if nil
	BeforeTest func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks)
	// TestBlockTime is the block time during test in unix seconds
	TestBlockTime int64
	// Test is the function that runs the actual test
	Test func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks)
	// AfterTestBlockTime is the block time during after test in unix seconds
	AfterTestBlockTime int64
	// AfterTest performs actions after the test is run, it can
	// be used to check if the state after Test is run matches
	// the result we expect.
	// Ignored if nil
	AfterTest func(t *testing.T, k keeper.Keeper, ctx types.Context, mocks *keeper.Mocks)
}

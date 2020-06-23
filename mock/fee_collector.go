package mock

import sdk "github.com/cosmos/cosmos-sdk/types"

type FeeCollector interface {
	CollectFee(ctx sdk.Context, fee sdk.Coin, addr sdk.AccAddress) error
}

type feeCollector struct {
	collectFee func(ctx sdk.Context, fee sdk.Coin, addr sdk.AccAddress) error
}

func (s *feeCollector) CollectFee(ctx sdk.Context, fee sdk.Coin, addr sdk.AccAddress) error {
	return s.collectFee(ctx, fee, addr)
}

type FeeCollectorMock struct {
	s *feeCollector
}

func (s *FeeCollectorMock) CollectFee(f func(ctx sdk.Context, fee sdk.Coin, addr sdk.AccAddress) error) {
	s.s.collectFee = f
}

func (s *FeeCollectorMock) Mock() FeeCollector {
	return s.s
}

func NewFeeCollector() *FeeCollectorMock {
	mock := &FeeCollectorMock{s: &feeCollector{}}
	// set default
	mock.CollectFee(func(ctx sdk.Context, fee sdk.Coin, addr sdk.AccAddress) error {
		return nil
	})
	return mock
}

package mock

import sdk "github.com/cosmos/cosmos-sdk/types"

type SupplyKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error
}

type supplyKeeper struct {
	sendCoinsFromAccountToModule func(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error
}

func (s *supplyKeeper) SendCoinsFromAccountToModule(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error {
	return s.sendCoinsFromAccountToModule(ctx, addr, moduleName, coins)
}

type SupplyKeeperMock struct {
	s *supplyKeeper
}

func (s *SupplyKeeperMock) SetSendCoinsFromAccountToModule(f func(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error) {
	s.s.sendCoinsFromAccountToModule = f
}

func (s *SupplyKeeperMock) Mock() SupplyKeeper {
	return s.s
}

func NewSupplyKeeper() *SupplyKeeperMock {
	mock := &SupplyKeeperMock{s: &supplyKeeper{}}
	// set default
	mock.SetSendCoinsFromAccountToModule(func(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error {
		return nil
	})
	return mock
}

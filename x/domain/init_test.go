package domain

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovnsd/x/configuration"
	"github.com/iov-one/iovnsd/x/domain/types"
)

// since the exposed interface for configuration keeper
// does not include set config as the domain module should
// not be able to change configuration state, then only
// in test cases we expose this method
type configurationSetter interface {
	SetConfig(ctx sdk.Context, config configuration.Config)
}

// getConfigSetter exposes the configurationSetter interface
// allowing yhe module to set configuration state, this should only
// be used for tests and will panic if the keeper provided can not
// be cast to configurationSetter
func getConfigSetter(keeper types.ConfigurationKeeper) configurationSetter {
	// check if the configuration keeper is also a config setter
	configSetter, ok := keeper.(configurationSetter)
	if !ok {
		panic(fmt.Sprintf("handleMsgRegisterDomain() cannot cast configuration keeper to configuration setter: got uncastable type: %T", keeper))
	}
	return configSetter
}

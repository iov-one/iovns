package configuration

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

// configKey defines the key used for the configuration
// since the configuration is only one the key will always be one
const configKey = ""

type paramSubspace interface {
}

type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramspace paramSubspace // TODO define what this is
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramspace paramSubspace) Keeper {
	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		paramspace: paramspace,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf(ModuleName))
}

// GetConfiguration returns the configuration of the blockchain
func (k Keeper) GetConfiguration(ctx sdk.Context) Config {
	store := ctx.KVStore(k.storeKey)
	confBytes := store.Get([]byte(configKey))
	if confBytes == nil {
		panic("no configuration available")
	}
	var conf Config
	k.cdc.MustUnmarshalBinaryBare(confBytes, &conf)
	// success
	return conf
}

// SetConfig updates or saves a new config in the store
func (k Keeper) SetConfig(ctx sdk.Context, conf Config) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(configKey), k.cdc.MustMarshalBinaryBare(conf))
}

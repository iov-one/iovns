package configuration

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/iov-one/iovns/x/configuration/client/cli"
	"github.com/iov-one/iovns/x/configuration/client/rest"
	"github.com/iov-one/iovns/x/configuration/types"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// nolint
// - - - FILL APP MODULE BASIC -- //
// AppModuleBasic implements the AppModuleBasic interface of the cosmos-sdk
type AppModuleBasic struct{}

func (AppModuleBasic) Name() string                   { return types.ModuleName }
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) { types.RegisterCodec(cdc) }
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}
func (AppModuleBasic) ValidateGenesis(b json.RawMessage) (err error) {
	var data GenesisState
	err = types.ModuleCdc.UnmarshalJSON(b, &data)
	if err != nil {
		return
	}
	return ValidateGenesis(data)
}

func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, router *mux.Router) {
	rest.RegisterRoutes(ctx, router, types.ModuleName, AvailableQueries())
}
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(types.StoreKey, cdc)
}
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(types.StoreKey, cdc)
}

// - - FILL APP MODULE - -
type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

func NewAppModule(k Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
	}
}
func (AppModule) Name() string                                       { return types.ModuleName }
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry)         {}
func (AppModule) Route() string                                      { return types.RouterKey }
func (AppModule) QuerierRoute() string                               { return types.QuerierRoute }
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
func (a AppModule) NewHandler() sdk.Handler {
	return NewHandler(a.keeper)
}
func (a AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(a.keeper)
}

func (a AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, a.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

func (a AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	genesisState := ExportGenesis(ctx, a.keeper)
	return types.ModuleCdc.MustMarshalJSON(genesisState)
}

package configuration

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/iov-one/iovnsd/x/configuration/client/cli"
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

func (AppModuleBasic) Name() string                   { return ModuleName }
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) { RegisterCodec(cdc) }
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState)
}
func (AppModuleBasic) ValidateGenesis(b json.RawMessage) (err error) {
	var data GenesisState
	err = ModuleCdc.UnmarshalJSON(b, data)
	if err != nil {
		return
	}
	return ValidateGenesis(data)
}

func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, router *mux.Router) {
	// TODO fill
}
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command { return cli.GetTxCmd(cdc) }
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
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
func (AppModule) Name() string                                       { return ModuleName }
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry)         {}
func (AppModule) Route() string                                      { return RouterKey }
func (AppModule) QuerierRoute() string                               { return QuerierRoute }
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
	ModuleCdc.MustUnmarshalJSON(data, genesisState)
	InitGenesis(ctx, a.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

func (a AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	genesisState := ExportGenesis(ctx, a.keeper)
	return ModuleCdc.MustMarshalJSON(genesisState)
}

package cli

import (
	"fmt"
	"github.com/iov-one/iovns/pkg/queries"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/iov-one/iovns/x/configuration/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd builds all the query commands for the module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// group config queries under a sub-command
	configQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	// add queries
	configQueryCmd.AddCommand(
		flags.GetCommands(
			getCmdQueryConfig(queryRoute, cdc),
			getCmdQueryFees(queryRoute, cdc),
		)...,
	)
	// return cmd list
	return configQueryCmd
}

// getCmdQueryConfig returns the command to get the configuration
func getCmdQueryConfig(route string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-config",
		Short: "gets the current configuration",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			path := fmt.Sprintf("custom/%s/%s", route, types.QueryConfig)
			resp, _, err := cliCtx.Query(path)
			if err != nil {
				return err
			}
			var jsonResp types.QueryConfigResponse
			err = queries.DefaultQueryDecode(resp, &jsonResp)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(jsonResp)
		},
	}
}

func getCmdQueryFees(route string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-fees",
		Short: "gets the current fees",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			path := fmt.Sprintf("custom/%s/%s", route, types.QueryFees)
			resp, _, err := cliCtx.Query(path)
			if err != nil {
				return err
			}
			var jsonResp types.QueryFeesResponse
			err = queries.DefaultQueryDecode(resp, &jsonResp)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(jsonResp)
		},
	}
}

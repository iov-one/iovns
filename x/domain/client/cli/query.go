package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/iov-one/iovnsd/x/domain/keeper"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/iov-one/iovnsd/x/domain/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group domain queries under a subcommand
	domainQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	domainQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdQueryDomain(queryRoute, cdc),
		)...,
	)

	return domainQueryCmd
}

// TODO: Add Query Commands

// GetCmdQueryDomain is the command that returns
func GetCmdQueryDomain(route string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "query [domain name]",
		Short: "query domain name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]
			path := fmt.Sprintf("custom/%s/%s/%s", route, keeper.QueryDomain, name)
			resp, _, err := cliCtx.Query(path)
			if err != nil {
				return cliCtx.PrintOutput(fmt.Sprintf("could not get domain information for %s: %s", name, err))
			}
			var result types.QueryResultDomain
			cdc.MustUnmarshalJSON(resp, &result)
			return cliCtx.PrintOutput(result)
		},
	}
}

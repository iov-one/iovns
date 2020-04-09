package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/iov-one/iovnsd/x/domain/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd builds the commands for queries in the domain module
func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	domainQueryCmd := &cobra.Command{
		Use:                        storeKey, // store key is same as module name
		Short:                      "querying commands for the domain module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	domainQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdQueryDomain(storeKey, cdc), // add query domain command
		)...,
	)

	return domainQueryCmd
}

// GetQueryDomain is the command used to query a domain by its name
func GetCmdQueryDomain(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:                "get [name]",
		Short:              "get a domain by its name",
		DisableFlagParsing: false,
		Args:               cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryDomain, name), nil)
			if err != nil {
				return err
			}
			// print output
			return cliCtx.PrintOutput(fmt.Sprintf("%s", res))
		},
	}
}

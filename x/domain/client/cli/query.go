package cli

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/tutils"
	"github.com/spf13/cobra"
)

// GetQueryCmd builds the commands for queries in the domain module
func GetQueryCmd(moduleQueryPath string, cdc *codec.Codec, queries []iovns.QueryCommand) *cobra.Command {
	domainQueryCmd := &cobra.Command{
		Use:                        moduleQueryPath, // store key is same as module name
		Short:                      "querying commands for the domain module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	domainQueryCmd.AddCommand(
		flags.GetCommands(
			generateQueryCommands(moduleQueryPath, cdc, queries)...,
		)...,
	)
	return domainQueryCmd
}

// generateQueryCommands generate the query commands from each iovns.QueryCommand type provided.
func generateQueryCommands(moduleQueryPath string, cdc *codec.Codec, queryCommands []iovns.QueryCommand) []*cobra.Command {
	cmds := make([]*cobra.Command, len(queryCommands))
	// generate commands
	for i, queryInterface := range queryCommands {
		// get query type so we can clone it
		// when we want to unmarshal data
		typ := tutils.GetPtrType(queryInterface)
		cmd := &cobra.Command{
			Use:   queryInterface.Use(),
			Short: queryInterface.Description(),
			Long:  queryInterface.Description(),
			RunE: func(cmd *cobra.Command, args []string) (err error) {
				// clone query type
				query := tutils.CloneFromType(typ).(iovns.QueryCommand)
				argParser, err := arg.NewParser(arg.Config{Program: query.Use()}, query)
				if err != nil {
					return
				}
				// parse args
				err = argParser.Parse(args)
				if err != nil {
					return
				}
				// generate cli CTX
				cliCtx := context.NewCLIContext().WithCodec(cdc)
				// set path
				path := fmt.Sprintf("custom/%s/%s", moduleQueryPath, query.QueryPath())
				// get request bytes
				b, err := iovns.DefaultQueryEncode(query)
				if err != nil {
					return
				}
				// do query
				res, _, err := cliCtx.QueryWithData(path, b)
				if err != nil {
					return err
				}
				// print output
				err = cliCtx.PrintOutput(res)
				if err != nil {
					return
				}
				// success
				return nil
			},
		}
		cmds[i] = cmd
	}
	return cmds
}

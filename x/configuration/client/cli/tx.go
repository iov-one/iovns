package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/iov-one/iovns/x/configuration/types"
	"github.com/spf13/cobra"
)

// GetTxCmd builds all the transaction commands for the configuration module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	domainTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	domainTxCmd.AddCommand(flags.PostCommands(
	// TODO: Add tx based commands
	// GetCmd<Action>(cdc)
	)...)

	return domainTxCmd
}

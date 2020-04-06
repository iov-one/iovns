package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/iov-one/iovnsd/x/configuration"
	"github.com/spf13/cobra"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	domainTxCmd := &cobra.Command{
		Use:                        configuration.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", configuration.ModuleName),
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

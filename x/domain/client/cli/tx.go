package cli

import (
	"bufio"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/iov-one/bnsd/x/domain/internal/types"
)

// GetTxCmd returns the transaction commands for this module
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

// Example:
//
// GetCmd<Action> is the CLI command for doing <Action>
// func GetCmd<Action>(cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "/* Describe your action cmd */",
// 		Short: "/* Provide a short description on the cmd */",
// 		Args:  cobra.ExactArgs(2), // Does your request require arguments
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)
// 			inBuf := bufio.NewReader(cmd.InOrStdin())
// 			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

// 			msg := types.NewMsg<Action>(/* Action params */)
// 			err = msg.ValidateBasic()
// 			if err != nil {
// 				return err
// 			}

// 			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
// 		},
// 	}
// }

func BroadcastRegisterDomain(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "register [name] ",
		Short: "register registers a new domain",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.MsgRegisterDomain{
				Name:         "",
				Admin:        nil,
				HasSuperuser: false,
				Broker:       nil,
				AccountRenew: 0,
			}
			if err = msg.ValidateBasic(); err != nil {
				return
			}
			// broadcast tx
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBuilder, []sdk.Msg{msg})
		},
	}
}

package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client/context"
	types3 "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/iov-one/iovns/x/fee/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	feeTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	feeTxCmd.AddCommand(flags.PostCommands(
		getCmdUpdateFeeConfiguration(cdc),
	)...)

	return feeTxCmd
}

func getCmdUpdateFeeConfiguration(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-fee-config",
		Short: "update fees configuration using a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithBroadcastMode(flags.BroadcastBlock)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get fees file
			feeFile, err := cmd.Flags().GetString("fees-file")
			if err != nil {
				return err
			}
			f, err := os.Open(feeFile)
			if err != nil {
				return fmt.Errorf("unable to open fee file: %s", err)
			}
			defer f.Close()
			newFees := new(types.FeeConfiguration)
			err = json.NewDecoder(f).Decode(newFees)
			if err != nil {
				return err
			}
			msg := types.MsgUpdateConfiguration{
				Fees:       newFees,
				Configurer: cliCtx.GetFromAddress(),
			}
			if err := msg.ValidateBasic(); err != nil {
				return fmt.Errorf("invalid tx: %w", err)
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBuilder, []types3.Msg{msg})
		},
	}
	cmd.Flags().String("fees-file", "fees.json", "fees file in json format")
	return cmd
}

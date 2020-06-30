package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/iov-one/iovns/x/configuration/types"
	"github.com/spf13/cobra"
)

// GetTxCmd clubs together all the CLI tx commands
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	configTxCmd := &cobra.Command{
		Use:                        storeKey,
		Short:                      fmt.Sprintf("%s transactions subcommands", storeKey),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	configTxCmd.AddCommand(flags.PostCommands(
		getCmdUpdateConfig(cdc),
		getCmdUpdateFees(cdc),
	)...)
	return configTxCmd
}

var defaultDuration, _ = time.ParseDuration("1h")

const defaultRegex = "^(.*?)?"
const defaultNumber = 1

func getCmdUpdateFees(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-fees",
		Short: "update fees using a file",
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
			newFees := new(types.Fees)
			err = json.NewDecoder(f).Decode(newFees)
			if err != nil {
				return err
			}
			msg := types.MsgUpdateFees{
				Fees:       newFees,
				Configurer: cliCtx.GetFromAddress(),
			}
			if err := msg.ValidateBasic(); err != nil {
				return fmt.Errorf("invalid tx: %w", err)
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBuilder, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String("fees-file", "fees.json", "fees file in json format")
	return cmd
}

func getCmdUpdateConfig(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-config",
		Short: "update domain configuration, provide file",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithBroadcastMode(flags.BroadcastBlock)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			configFile, err := cmd.Flags().GetString("config-file")
			if err != nil {
				return err
			}
			f, err := os.Open(configFile)
			if err != nil {
				return fmt.Errorf("unable to open config file: %s", err)
			}
			defer f.Close()
			newConfig := new(types.Config)
			err = json.NewDecoder(f).Decode(newConfig)
			if err != nil {
				return fmt.Errorf("unable to decode json: %s", err)
			}
			msg := types.MsgUpdateConfig{
				Signer:           cliCtx.GetFromAddress(),
				NewConfiguration: *newConfig,
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBuilder, []sdk.Msg{msg})
		},
	}
	// add flags
	cmd.Flags().String("config-file", "config.json", "config file in json format")
	return cmd
}

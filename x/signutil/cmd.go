package signutil

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
)

const DefaultChainID = "signed-message-v1"
const DefaultAccountNumber uint64 = 0
const DefaultSequence uint64 = 0

// getTxCmd clubs together all the CLI tx commands
func getTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	configTxCmd := &cobra.Command{
		Use:                        storeKey,
		Short:                      fmt.Sprintf("%s transactions subcommands", storeKey),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	configTxCmd.AddCommand(flags.PostCommands(
		signCmd(cdc),
		verifyCmd(cdc),
	)...)
	return configTxCmd
}

func signCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "outputs the json string to signCmd",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get file
			file, err := cmd.Flags().GetString("file")
			if err != nil {
				return err
			}
			// get text
			text, err := cmd.Flags().GetString("text")
			if err != nil {
				return err
			}
			pairs, err := cmd.Flags().GetStringArray("pair")
			if err != nil {
				return err
			}
			if (text != "") && (file != "" || len(pairs) != 0) || (file != "" && len(pairs) != 0) {
				return fmt.Errorf("only one of text, file, pairs can be specified")
			}
			msg := MsgSignText{
				Message: "",
				Pairs:   nil,
				Signer:  cliCtx.GetFromAddress(),
			}
			switch true {
			case text != "":
				msg.Message = text
			case file != "":
				f, err := os.Open(file)
				if err != nil {
					return err
				}
				defer f.Close()
				buf := new(bytes.Buffer)
				_, err = io.Copy(buf, f)
				if err != nil {
					return err
				}
				msg.Message = buf.String()
			case len(pairs) != 0:
				kv := make([]Pair, len(pairs))
				for i, raw := range pairs {
					split := strings.Split(raw, "=")
					if len(split) < 1 {
						return fmt.Errorf("invalid formatted value: %s", raw)
					}
					key := split[0]
					value := strings.Join(split[1:], "=")
					kv[i] = Pair{
						Key:   key,
						Value: value,
					}
				}
				msg.Pairs = kv
			default:
				return fmt.Errorf("either file or text flag must be specified")
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBuilder, []sdk.Msg{msg})
		},
	}
	cmd.Flags().StringP("file", "f", "", "file to signCmd")
	cmd.Flags().StringP("text", "t", "", "string to signCmd")
	cmd.Flags().StringArrayP("pair", "p", nil, "key value pairs, specified as key=value")
	return cmd
}

func verifyCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "verify a signature from a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := cmd.Flags().GetString("file")
			if err != nil {
				return err
			}
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			b, err := ioutil.ReadAll(f)
			var tx auth.StdTx
			err = cdc.UnmarshalJSON(b, &tx)
			if err != nil {
				return err
			}
			chainID, err := cmd.Flags().GetString(flags.FlagChainID)
			if err != nil {
				return err
			}
			if chainID == "" {
				chainID = DefaultChainID
			}
			accountNumber, err := cmd.Flags().GetUint64(flags.FlagAccountNumber)
			if err != nil {
				return err
			}
			if accountNumber == 0 {
				accountNumber = DefaultAccountNumber
			}
			sequence, err := cmd.Flags().GetUint64(flags.FlagSequence)
			if err != nil {
				return err
			}
			if sequence == 0 {
				sequence = DefaultSequence
			}
			if err = Verify(tx, chainID, accountNumber, sequence); err != nil {
				return err
			}
			msgs := tx.GetMsgs()
			if len(msgs) != 1 {
				return fmt.Errorf("Expected 1 msg but got %d.", len(msgs))
			}
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			if cliCtx.OutputFormat == "json" {
				var bin []byte
				var err error
				if cliCtx.Indent {
					bin, err = json.MarshalIndent(msgs[0], "", "  ")
				} else {
					bin, err = json.Marshal(msgs[0])
				}
				if err != nil {
					return err
				}
				fmt.Fprint(cmd.OutOrStdout(), string(bin))
			} else {
				var msg MsgSignText
				err = cdc.UnmarshalJSON(msgs[0].GetSignBytes(), &msg)
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "signer: %s\nmessage: %s\n", msg.Signer, msg.Message)
			}
			return nil
		},
	}
	cmd.Flags().StringP("file", "f", "", "signed transaction file")
	return cmd
}

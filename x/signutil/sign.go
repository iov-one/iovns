package signutil

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

type Pair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MsgArbitrarySignature struct {
	Message string         `json:"message,omitempty"`
	Pairs   []Pair         `json:"pairs,omitempty"`
	Signer  sdk.AccAddress `json:"signer"`
}

func (m MsgArbitrarySignature) Route() string {
	return ModuleName
}

func (m MsgArbitrarySignature) Type() string {
	return "arbitrary_signature"
}

func (m MsgArbitrarySignature) ValidateBasic() error {
	if len(m.Message) == 0 && len(m.Pairs) == 0 {
		return fmt.Errorf("empty msg and pairs")
	}
	if m.Signer.Empty() {
		return fmt.Errorf("missing signer")
	}
	return nil
}

// GetSignBytes implements sdk.Message
func (m MsgArbitrarySignature) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Message
func (m MsgArbitrarySignature) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{m.Signer} }

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
		sign(cdc),
		verify(cdc),
	)...)
	return configTxCmd
}

func sign(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "outputs the json string to sign",
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
			msg := MsgArbitrarySignature{
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
	cmd.Flags().StringP("file", "f", "", "file to sign")
	cmd.Flags().StringP("text", "t", "", "string to sign")
	cmd.Flags().StringArrayP("pair", "p", nil, "key value pairs, specified as key=value")
	return cmd
}

func verify(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "verify",
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("not implemented")
		},
	}
}

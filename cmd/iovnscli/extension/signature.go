package extension

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"os"
)

type SignatureSchema struct {
	ChanID  string         `json:"@chain_id"`
	Type    string         `json:"@type"`
	Message []byte         `json:"text"`
	Sig     string         `json:"sig"`
	PubKey  sdk.AccAddress `json:"address"`
}

type sigCommand struct {
	file string
}

func (s *sigCommand) applyFlags(flag *pflag.FlagSet) {
	flag.StringP("file", "f", "", "")
}

func (s *sigCommand) extractFlags(flag *pflag.FlagSet) (err error) {
	s.file, err = flag.GetString("file")
	if err != nil {
		return err
	}
	return
}

func SignatureCommand() *cobra.Command {
	req := new(sigCommand)
	cmd := &cobra.Command{
		Use: "sign",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// extract flags
			err = req.extractFlags(cmd.Flags())
			if err != nil {
				return
			}
			// retrieve file
			f, err := os.Open(req.file)
			if err != nil {
				return
			}
			defer f.Close()
			buf := &bytes.Buffer{}
			_, err = io.Copy(buf, f)
			if err != nil {
				return
			}
			cliCtx := context.NewCLIContext()
			kb, err := keys.NewKeyring(sdk.KeyringServiceName(), viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), cmd.InOrStdin())
			if err != nil {
				return err
			}
			sig, _, err := kb.Sign(cliCtx.GetFromName(), keys.DefaultBIP39Passphrase, buf.Bytes())
			if err != nil {
				return
			}
			messageJSON, err := json.Marshal(&SignatureSchema{
				ChanID:  cliCtx.ChainID,
				Type:    "message",
				Sig:     string(sig),
				Message: buf.Bytes(),
				PubKey:  cliCtx.GetFromAddress(),
			})
			if err != nil {
				return
			}
			cmd.Println(fmt.Sprintf("%s", messageJSON))
			return nil
		},
	}
	req.applyFlags(cmd.Flags())
	return cmd
}

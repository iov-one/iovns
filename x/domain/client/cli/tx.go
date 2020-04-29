package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/types"
	"github.com/spf13/cobra"
	"os"
)

// GetTxCmd clubs together all the CLI tx commands
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	domainTxCmd := &cobra.Command{
		Use:                        storeKey,
		Short:                      fmt.Sprintf("%s transactions subcommands", storeKey),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	domainTxCmd.AddCommand(flags.PostCommands(
		getCmdRegisterDomain(cdc),
		getCmdAddAccountCerts(cdc),
		getCmdFlushDomain(cdc),
		getCmdTransferAccount(cdc),
		getCmdTransferDomain(cdc),
		getCmdReplaceAccountTargets(cdc),
		getCmdDelDomain(cdc),
		getCmdDelAccount(cdc),
		getCmdRenewDomain(cdc),
		getCmdRenewAccount(cdc),
		getCmdDelAccountCerts(cdc),
		getCmdRegisterAccount(cdc),
	)...)
	return domainTxCmd
}

func getCmdTransferDomain(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-domain",
		Short: "transfer a domain",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return
			}
			newOwner, err := cmd.Flags().GetString("new-owner")
			if err != nil {
				return err
			}
			// get sdk.AccAddress from string
			newOwnerAddr, err := sdk.AccAddressFromBech32(newOwner)
			if err != nil {
				return
			}
			// build msg
			msg := &types.MsgTransferDomain{
				Domain:   domain,
				Owner:    cliCtx.GetFromAddress(),
				NewAdmin: newOwnerAddr,
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
	cmd.Flags().String("domain", "", "the domain name to transfer")
	cmd.Flags().String("new-owner", "", "the new owner address in bech32 format")
	//
	return cmd
}

func getCmdTransferAccount(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-account",
		Short: "transfer an account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return
			}
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return
			}
			newOwner, err := cmd.Flags().GetString("new-owner")
			if err != nil {
				return err
			}
			// get sdk.AccAddress from string
			newOwnerAddr, err := sdk.AccAddressFromBech32(newOwner)
			if err != nil {
				return
			}
			// build msg
			msg := &types.MsgTransferAccount{
				Domain:   domain,
				Name:     name,
				Owner:    cliCtx.GetFromAddress(),
				NewOwner: newOwnerAddr,
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
	cmd.Flags().String("domain", "", "the domain name of account")
	cmd.Flags().String("name", "", "the name of the account you want to transfer")
	cmd.Flags().String("new-owner", "", "the new owner address in bech32 format")
	//
	return cmd
}

func getCmdReplaceAccountTargets(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replace-targets",
		Short: "replace account targets",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return
			}
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return
			}
			targetsPath, err := cmd.Flags().GetString("src")
			if err != nil {
				return err
			}
			// open targets file
			f, err := os.Open(targetsPath)
			if err != nil {
				return err
			}
			defer f.Close()
			// unmarshal targets
			var targets []iovns.BlockchainAddress
			err = json.NewDecoder(f).Decode(&targets)
			if err != nil {
				return
			}
			// build msg
			msg := &types.MsgReplaceAccountTargets{
				Domain:     domain,
				Name:       name,
				NewTargets: targets,
				Owner:      cliCtx.GetFromAddress(),
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
	cmd.Flags().String("domain", "", "the domain name of account")
	cmd.Flags().String("name", "", "the name of the account whose targets you want to replace")
	cmd.Flags().String("src", "targets.json", "the file containing the new targets in json format")
	// return cmd
	return cmd
}

func getCmdFlushDomain(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flush-domain",
		Short: "flush a domain",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return
			}
			// build msg
			msg := &types.MsgFlushDomain{
				Domain: domain,
				Owner:  cliCtx.GetFromAddress(),
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
	cmd.Flags().String("domain", "", "name of the domain you want to flush")
	//
	return cmd
}

func getCmdDelDomain(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "del-domain",
		Short: "delete a domain",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return
			}
			// build msg
			msg := &types.MsgDeleteDomain{
				Domain: domain,
				Owner:  cliCtx.GetFromAddress(),
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
	cmd.Flags().String("domain", "", "name of the domain you want to delete")
	//
	return cmd
}

func getCmdDelAccount(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "del-account",
		Short: "delete an account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return
			}
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return
			}
			// build msg
			msg := &types.MsgDeleteAccount{
				Domain: domain,
				Name:   name,
				Owner:  cliCtx.GetFromAddress(),
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
	cmd.Flags().String("domain", "", "the domain name of account")
	cmd.Flags().String("name", "", "the name of the account you want to delete")
	//
	return cmd
}

func getCmdRenewDomain(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "renew-domain",
		Short: "renew a domain",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return
			}
			// build msg
			msg := &types.MsgRenewDomain{
				Domain: domain,
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
	cmd.Flags().String("domain", "", "name of the domain you want to renew")
	// return
	return cmd
}

func getCmdRenewAccount(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "renew-account",
		Short: "renew an account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return
			}
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return
			}
			// build msg
			msg := &types.MsgRenewAccount{
				Domain: domain,
				Name:   name,
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
	cmd.Flags().String("domain", "", "domain name of the account")
	cmd.Flags().String("name", "", "account name you want to renew")
	// return
	return cmd
}

func getCmdDelAccountCerts(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "del-certs",
		Short: "delete certificates of an account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return
			}
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return
			}
			newCert, err := cmd.Flags().GetBytesHex("cert")
			if err != nil {
				return
			}
			// build msg
			msg := &types.MsgDeleteAccountCertificate{
				Domain:            domain,
				Name:              name,
				Owner:             cliCtx.GetFromAddress(),
				DeleteCertificate: newCert,
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
	cmd.Flags().String("domain", "", "domain name of the account")
	cmd.Flags().String("name", "", "account name")
	cmd.Flags().BytesHex("cert", []byte{}, "hex bytes of the certificate you want to delete")
	// return cmd
	return cmd
}

func getCmdAddAccountCerts(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-certs",
		Short: "add certificates to account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return
			}
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return
			}
			newCert, err := cmd.Flags().GetBytesHex("cert")
			if err != nil {
				return
			}
			// build msg
			msg := &types.MsgAddAccountCertificates{
				Domain:         domain,
				Name:           name,
				Owner:          cliCtx.GetFromAddress(),
				NewCertificate: newCert,
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
	cmd.Flags().String("domain", "", "domain of the account")
	cmd.Flags().String("name", "", "name of the account")
	cmd.Flags().BytesHex("cert", []byte{}, "hex bytes of the certificate you want to add")
	// return cmd
	return cmd
}

// getCmdRegisterAccount is the cli command to register accounts
func getCmdRegisterAccount(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "register-account",
		Short:                      "register an account",
		SuggestionsMinimumDistance: 2,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return
			}
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return
			}
			// build msg
			msg := &types.MsgRegisterAccount{
				Domain: domain,
				Name:   name,
				Owner:  cliCtx.GetFromAddress(),
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBuilder, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String("domain", "", "the existing domain name for your account")
	cmd.Flags().String("name", "", "the name of your account")
	return cmd
}

func getCmdRegisterDomain(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-domain",
		Short: "register a domain",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			hasSuperUser, err := cmd.Flags().GetBool("has-superuser")
			if err != nil {
				return err
			}
			accountRenew, err := cmd.Flags().GetInt64("account-renew")
			if err != nil {
				return err
			}
			msg := &types.MsgRegisterDomain{
				Name:         domain,
				Admin:        cliCtx.GetFromAddress(),
				HasSuperuser: hasSuperUser,
				Broker:       nil,
				AccountRenew: accountRenew,
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
	cmd.Flags().String("domain", "", "name of the domain you want to register")
	cmd.Flags().Bool("has-superuser", true, "define if this domain has a superuser or not")
	cmd.Flags().Int64("account-renew", 10000000, "account duration in seconds before expiration")
	return cmd
}

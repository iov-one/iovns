package cli

import (
	"bufio"
	"fmt"
	"regexp"

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
	)...)
	return configTxCmd
}

func getCmdUpdateConfig(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-config",
		Short: "update domain configuration",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get flags
			ownerStr, err := cmd.Flags().GetStringSlice("owners")
			if err != nil {
				return
			}
			var owners []sdk.AccAddress
			for _, o := range ownerStr {
				addr, err := sdk.AccAddressFromBech32(o)
				if err != nil {
					return err
				}
				owners = append(owners, addr)
			}

			validDomain, err := cmd.Flags().GetString("valid-domain")
			if err != nil {
				return err
			}
			_, err = regexp.Compile(validDomain)
			if err != nil {
				return err
			}

			validName, err := cmd.Flags().GetString("valid-name")
			if err != nil {
				return err
			}
			_, err = regexp.Compile(validName)
			if err != nil {
				return err
			}

			validBlockchainID, err := cmd.Flags().GetString("valid-blockchain-id")
			if err != nil {
				return err
			}
			_, err = regexp.Compile(validBlockchainID)
			if err != nil {
				return err
			}

			validBlockchainAddress, err := cmd.Flags().GetString("valid-blockchain-address")
			if err != nil {
				return err
			}
			_, err = regexp.Compile(validBlockchainAddress)
			if err != nil {
				return err
			}

			domainRenew, err := cmd.Flags().GetUint64("domain-renew")
			if err != nil {
				return err
			}

			domainGracePeriod, err := cmd.Flags().GetUint64("domain-grace-period")
			if err != nil {
				return err
			}

			config := types.Config{
				Owners:                 owners,
				ValidDomain:            validDomain,
				ValidName:              "",
				ValidBlockchainID:      "",
				ValidBlockchainAddress: "",
				DomainRenew:            0,
				DomainGracePeriod:      0,
			}
			// build msg
			msg := &types.MsgUpdateConfig{
				Signers:          nil,
				NewConfiguration: types.Config{},
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
	cmd.Flags().StringSlice("owners", []string{}, "configuration owners in bech32 format. Separated by comma(,)")
	cmd.Flags().String("valid-domain", "", "regexp that determines if domain name is valid or not")
	cmd.Flags().String("valid-name", "", "regexp that determines if account name is valid or not")
	cmd.Flags().String("valid-blockchain-id", "", "regexp that determines if blockchain id is valid or not")
	cmd.Flags().String("valid-blockchain-address", "", "regexp that determines if blockchain address is valid or not")
	cmd.Flags().Uint64("domain-renew", 10000000, "domain renewal duration in seconds before expiration")
	cmd.Flags().Uint64("domain-grace-period", 10000000, "domain grace period duration in seconds")
	return cmd
}

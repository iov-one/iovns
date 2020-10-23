package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/iov-one/iovns/pkg/queries"
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
		getCmdUpdateFee(cdc),
	)...)
	return configTxCmd
}

var defaultDuration, _ = time.ParseDuration("1h")

const defaultRegex = "^(.*?)?"
const defaultNumber = 1

func getCmdUpdateFee(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update-fee",
		Short:   "update fee parameter(s) using key=value pair(s) noting that all numeric values should be specified in euros",
		Example: "iovnscli tx configuration update-fee --from iovSAS --parameter fee_coin_price=0.29 --parameter register_domain_default=12",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithBroadcastMode(flags.BroadcastBlock)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// query the existing fee configuration
			path := fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryFees)
			resp, _, err := cliCtx.Query(path)
			if err != nil {
				return err
			}
			var jsonResp types.QueryFeesResponse
			err = queries.DefaultQueryDecode(resp, &jsonResp)
			if err != nil {
				return err
			}
			fees := jsonResp.Fees
			// set new values via json tag https://gist.github.com/lelandbatey/a5c957b537bed39d1d6fb202c3b8de06
			v := reflect.ValueOf(fees).Elem()
			findJsonName := func(t reflect.StructTag) string {
				if jt, ok := t.Lookup("json"); ok {
					return strings.Split(jt, ",")[0]
				}
				panic(fmt.Errorf("tag provided does not define a json tag"))
			}
			fieldNames := map[string]int{}
			for i := 0; i < v.NumField(); i++ {
				typeField := v.Type().Field(i)
				tag := typeField.Tag
				jname := findJsonName(tag)
				fieldNames[jname] = i
			}
			// read the fee parameter pair(s)
			pairs, err := cmd.Flags().GetStringArray("parameter")
			if err != nil {
				return err
			}
			// iterate over the pair(s) and set the fee parameter value
			for _, raw := range pairs {
				split := strings.Split(raw, "=")
				if len(split) != 2 {
					return fmt.Errorf("invalid pair: %s", raw)
				}
				key := split[0]
				value := split[1]
				fieldNum, ok := fieldNames[key]
				if !ok {
					return fmt.Errorf("%s is not a valid fee configuration variable", key)
				}
				if key == "fee_coin_denom" { // special case of a string value
					fees.FeeCoinDenom = value
				} else { // decimal values
					dec, err := sdk.NewDecFromStr(value)
					if err != nil {
						return fmt.Errorf("failed to make %s a decimal value for key %s", value, key)
					}
					if key == "fee_coin_price" { // special case of converting euros to megaeuros so that fees can be specified in euros, not uiov
						million, err := sdk.NewDecFromStr("1000000")
						if err != nil {
							return fmt.Errorf("failed to make '1000000' a decimal value")
						}
						dec = dec.Quo(million)
					}
					fieldVal := v.Field(fieldNum)
					fieldVal.Set(reflect.ValueOf(dec))
				}
			}
			// submit the tx
			msg := types.MsgUpdateFees{
				Fees:       fees,
				Configurer: cliCtx.GetFromAddress(),
			}
			if err := msg.ValidateBasic(); err != nil {
				return fmt.Errorf("invalid tx: %w", err)
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBuilder, []sdk.Msg{msg})
		},
	}
	cmd.Flags().StringArrayP("parameter", "p", nil, "key/value pairs, specified as key=value")
	return cmd
}

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
		Short: "update domain configuration, provide the values you want to override in current configuration",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithBroadcastMode(flags.BroadcastBlock)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			config := &types.Config{}
			if !cliCtx.GenerateOnly {
				rawCfg, _, err := cliCtx.QueryStore([]byte(types.ConfigKey), types.StoreKey)
				if err != nil {
					return err
				}
				cdc.MustUnmarshalBinaryBare(rawCfg, config)
			}
			var signer sdk.AccAddress
			// if tx is not generate only, use --from flag as signer, otherwise get it from signer flag
			if !cliCtx.GenerateOnly {
				signer = cliCtx.FromAddress
			} else {
				signerStr, err := cmd.Flags().GetString("signer")
				if err != nil {
					return err
				}
				signer, err = sdk.AccAddressFromBech32(signerStr)
				if err != nil {
					return err
				}
			}
			// get flags
			configurerStr, err := cmd.Flags().GetString("configurer")
			if err != nil {
				return
			}
			if configurerStr != "" {
				configurer, err := sdk.AccAddressFromBech32(configurerStr)
				if err != nil {
					return err
				}
				config.Configurer = configurer
			}
			validDomainName, err := cmd.Flags().GetString("valid-domain-name")
			if err != nil {
				return err
			}
			if validDomainName != defaultRegex {
				config.ValidDomainName = validDomainName
			}
			validAccountName, err := cmd.Flags().GetString("valid-account-name")
			if err != nil {
				return err
			}
			if validAccountName != defaultRegex {
				config.ValidAccountName = validAccountName
			}
			validURI, err := cmd.Flags().GetString("valid-uri")
			if err != nil {
				return err
			}
			if validURI != defaultRegex {
				config.ValidURI = validURI
			}
			validResource, err := cmd.Flags().GetString("valid-resource")
			if err != nil {
				return err
			}
			if validResource != defaultRegex {
				config.ValidResource = validResource
			}
			domainRenew, err := cmd.Flags().GetDuration("domain-renew-period")
			if err != nil {
				return err
			}
			if domainRenew != defaultDuration {
				config.DomainRenewalPeriod = domainRenew
			}
			domainRenewCountMax, err := cmd.Flags().GetUint32("domain-renew-count-max")
			if err != nil {
				return err
			}
			if domainRenewCountMax != defaultNumber {
				config.DomainRenewalCountMax = domainRenewCountMax
			}
			domainGracePeriod, err := cmd.Flags().GetDuration("domain-grace-period")
			if err != nil {
				return err
			}
			if domainGracePeriod != defaultNumber {
				config.DomainGracePeriod = domainGracePeriod
			}
			accountRenewPeriod, err := cmd.Flags().GetDuration("account-renew-period")
			if err != nil {
				return err
			}
			if accountRenewPeriod != defaultNumber {
				config.AccountRenewalPeriod = accountRenewPeriod
			}
			accountRenewCountMax, err := cmd.Flags().GetUint32("account-renew-count-max")
			if err != nil {
				return err
			}
			if accountRenewCountMax != defaultNumber {
				config.AccountRenewalCountMax = accountRenewCountMax
			}
			accountGracePeriod, err := cmd.Flags().GetDuration("account-grace-period")
			if err != nil {
				return err
			}
			if accountGracePeriod != defaultDuration {
				config.AccountGracePeriod = accountGracePeriod
			}
			resourceMax, err := cmd.Flags().GetUint32("resource-max")
			if err != nil {
				return err
			}
			if resourceMax != defaultNumber {
				config.ResourcesMax = resourceMax
			}
			certificateSizeMax, err := cmd.Flags().GetUint64("certificate-size-max")
			if err != nil {
				return err
			}
			if certificateSizeMax != defaultNumber {
				config.CertificateSizeMax = certificateSizeMax
			}
			certificateCountMax, err := cmd.Flags().GetUint32("certificate-count-max")
			if err != nil {
				return err
			}
			if certificateCountMax != defaultNumber {
				config.CertificateCountMax = certificateCountMax
			}
			metadataSizeMax, err := cmd.Flags().GetUint64("metadata-size-max")
			if err != nil {
				return err
			}
			if metadataSizeMax != defaultNumber {
				config.MetadataSizeMax = metadataSizeMax
			}

			if err := config.Validate(); err != nil {
				return err
			}
			// build msg
			msg := &types.MsgUpdateConfig{
				Signer:           signer,
				NewConfiguration: *config,
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
	cmd.Flags().String("signer", "", "current configuration owner, for offline usage, otherwise --from is used")
	cmd.Flags().String("configurer", "", "new configuration owner")
	cmd.Flags().String("valid-domain-name", defaultRegex, "regexp that determines if domain name is valid or not")
	cmd.Flags().String("valid-account-name", defaultRegex, "regexp that determines if account name is valid or not")
	cmd.Flags().String("valid-uri", defaultRegex, "regexp that determines if uri is valid or not")
	cmd.Flags().String("valid-resource", defaultRegex, "regexp that determines if resource is valid or not")

	cmd.Flags().Duration("domain-renew-period", defaultDuration, "domain renewal duration in seconds before expiration")
	cmd.Flags().Uint32("domain-renew-count-max", uint32(defaultNumber), "maximum number of applicable domain renewals")
	cmd.Flags().Duration("domain-grace-period", defaultDuration, "domain grace period duration in seconds")

	cmd.Flags().Duration("account-renew-period", defaultDuration, "domain renewal duration in seconds before expiration")
	cmd.Flags().Uint32("account-renew-count-max", uint32(defaultNumber), "maximum number of applicable account renewals")
	cmd.Flags().Duration("account-grace-period", defaultDuration, "account grace period duration in seconds")

	cmd.Flags().Uint32("resource-max", uint32(defaultNumber), "maximum number of resources could be saved under an account")
	cmd.Flags().Uint64("certificate-size-max", uint64(defaultNumber), "maximum size of a certificate that could be saved under an account")
	cmd.Flags().Uint32("certificate-count-max", uint32(defaultNumber), "maximum number of certificates that could be saved under an account")
	cmd.Flags().Uint64("metadata-size-max", uint64(defaultNumber), "maximum size of metadata that could be saved under an account")
	return cmd
}

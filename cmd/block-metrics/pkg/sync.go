package pkg

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/prometheus/common/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/starname/types"

	"github.com/pkg/errors"
)

const syncRetryTimeout = 3 * time.Second

// Sync uploads to local store all blocks that are not present yet, starting
// with the blocks with the lowest height first. It always returns the number of
// blocks inserted, even if returning an error.
func Sync(ctx context.Context, tmc *TendermintClient, st *Store, denom string, urlLCD string) (uint, error) {
	var (
		inserted        uint
		syncedHeight    int64
		lastKnownHeight int64
	)

	switch block, err := st.LatestBlock(ctx); {
	case ErrNotFound.Is(err):
		syncedHeight = 0
		if err := st.InsertGenesis(ctx, tmc); err != nil {
			return inserted, errors.Wrapf(err, "fetch genesis")
		}
	case err == nil:
		syncedHeight = block.Height
	default:
		return inserted, errors.Wrap(err, "latest block")
	}

	// begin a database transaction
	if err := st.BatchBegin(ctx); err != nil {
		return inserted, errors.Wrap(err, "st.BatchBegin()")
	}

	for {
		nextHeight := syncedHeight + 1
		if lastKnownHeight < nextHeight {
			info, err := AbciInfo(tmc)
			if err != nil {
				return inserted, errors.Wrap(err, "info")
			}

			lastKnownHeight = info.LastBlockHeight
		}

		if lastKnownHeight < nextHeight {
			select {
			case <-ctx.Done():
				return inserted, ctx.Err()
			case <-time.After(syncRetryTimeout):
			}
			// make sure we don't run into the bug where we try to retrieve a commit for non-existent height
			continue
		}

		c, err := Commit(ctx, tmc, nextHeight)
		if err != nil {
			// BUG this can happen when the commit does not exist.
			// There is no sane way to distinguish this case from
			// any other tendermint API error.
			return inserted, errors.Wrapf(err, "blocks for %d", syncedHeight+1)
		}
		syncedHeight = c.Height

		tmblock, err := FetchBlock(ctx, tmc, nextHeight)
		if err != nil {
			return inserted, errors.Wrapf(err, "blocks for %d", syncedHeight+1)
		}

		fee := sdk.ZeroInt()
		block := Block{
			Height:  c.Height,
			Hash:    hex.EncodeToString(c.Hash),
			Time:    c.Time.UTC(),
			FeeFrac: fee.Uint64(),
		}
		if err := st.InsertBlock(ctx, block); err != nil {
			return inserted, errors.Wrapf(err, "insert block %d", c.Height)
		}

		for _, tx := range tmblock.Transactions {
			coins := tx.Fee.Amount
			for _, c := range coins {
				if c.Denom != denom {
					return 1, errors.Wrapf(ErrDenom, "not supported denom: %s, expected %s", c.Denom, denom)
				}
				fee = fee.Add(c.Amount)
			}
			if err := routeMsgs(ctx, st, tx.Msgs, c.Height, denom, urlLCD); err != nil {
				log.Error(errors.Wrapf(err, "height %d", c.Height))
			}
		}

		// commit the database transaction, potentially in a batch
		if lastKnownHeight-nextHeight == 0 || inserted%100 == 0 {
			if err = st.BatchCommit(ctx); err != nil {
				return inserted, errors.Wrapf(err, "inserted %d; failed at block %d", inserted, c.Height)
			}
		}

		inserted++
	}
}

func routeMsgs(ctx context.Context, st *Store, msgs []sdk.Msg, height int64, denom string, urlLCD string) error {
	// allocate a slice with the maximum needed capacity
	queries := make([]*LcdRequestData, 0, len(msgs))[:]

	for _, msg := range msgs {
		var accountID int64
		params := make(map[string]string)

		switch m := msg.(type) {
		case *types.MsgRegisterDomain:
			if id, err := st.RegisterDomain(ctx, m, height); err != nil {
				return errors.Wrap(err, "register domain message")
			} else {
				accountID = id
				params["action"] = "register_domain"
				params["owner"] = m.Admin.String()
				params["domain_name"] = m.Name
				params["domain_type"] = string(m.DomainType)
			}
		case *types.MsgDeleteDomain:
			if id, err := st.DeleteDomain(ctx, m, height); err != nil {
				return errors.Wrapf(err, "delete domain message, domain name: %s", m.Domain)
			} else {
				accountID = id
				params["action"] = "delete_domain"
				params["owner"] = m.Owner.String()
				params["domain_name"] = m.Domain
			}
		case *types.MsgTransferDomain:
			if id, err := st.TransferDomain(ctx, m, height); err != nil {
				return errors.Wrapf(err, "transfer domain message, domain name: %s", m.Domain)
			} else {
				accountID = id
				params["action"] = "transfer_domain"
				params["owner"] = m.Owner.String()
				params["domain_name"] = m.Domain
				params["new_domain_owner"] = m.NewAdmin.String()
			}
		case *types.MsgRegisterAccount:
			if id, err := st.RegisterAccount(ctx, m, height); err != nil {
				return errors.Wrapf(err, "register account message, domain name: %s, account name: %s", m.Domain, m.Name)
			} else {
				accountID = id
				params["action"] = "register_account"
				params["owner"] = m.Owner.String()
				params["domain_name"] = m.Domain
				params["account_name"] = m.Name
			}
		case *types.MsgDeleteAccount:
			if id, err := st.DeleteAccount(ctx, m, height); err != nil {
				return errors.Wrapf(err, "delete account message, domain name: %s, account name: %s", m.Domain, m.Name)
			} else {
				accountID = id
				params["action"] = "delete_account"
				params["owner"] = m.Owner.String()
				params["domain_name"] = m.Domain
				params["account_name"] = m.Name
			}
		case *types.MsgTransferAccount:
			if id, err := st.TransferAccount(ctx, m, height); err != nil {
				return errors.Wrapf(err, "transfer account message, domain name: %s, account name: %s", m.Domain, m.Name)
			} else {
				accountID = id
				params["action"] = "transfer_account"
				params["owner"] = m.Owner.String()
				params["domain_name"] = m.Domain
				params["account_name"] = m.Name
				params["new_account_owner"] = m.NewOwner.String()
			}
		case *types.MsgReplaceAccountResources:
			if id, err := st.ReplaceAccountResources(ctx, m, height); err != nil {
				return errors.Wrapf(err, "replace account resources message, domain name: %s, account name: %s", m.Domain, m.Name)
			} else {
				accountID = id
				params["action"] = "replace_account_resources"
				params["owner"] = m.Owner.String()
				params["domain_name"] = m.Domain
				params["account_name"] = m.Name
			}
		case *types.MsgReplaceAccountMetadata:
			if id, err := st.ReplaceAccountMetadata(ctx, m, height); err != nil {
				return errors.Wrapf(err, "replace account metadata message, domain name: %s, account name: %s", m.Domain, m.Name)
			} else {
				accountID = id
				params["action"] = "set_account_metadata"
				params["owner"] = m.Owner.String()
				params["domain_name"] = m.Domain
				params["account_name"] = m.Name
				params["new_metadata"] = m.NewMetadataURI
			}
		case *types.MsgAddAccountCertificates:
			if id, err := st.AddAccountCertificates(ctx, m, height); err != nil {
				return errors.Wrapf(err, "add account certificates message, domain name: %s, account name: %s", m.Domain, m.Name)
			} else {
				accountID = id
				params["action"] = "add_certificates_account"
				params["owner"] = m.Owner.String()
				params["domain_name"] = m.Domain
				params["account_name"] = m.Name
				params["new_certificate"] = hex.EncodeToString(m.NewCertificate)
			}
		case *types.MsgDeleteAccountCertificate:
			if id, err := st.DeleteAccountCerts(ctx, m, height); err != nil {
				return errors.Wrapf(err, "delete account certificates message, domain name: %s, account name: %s", m.Domain, m.Name)
			} else {
				accountID = id
				params["action"] = "delete_certificate_account"
				params["owner"] = m.Owner.String()
				params["domain_name"] = m.Domain
				params["account_name"] = m.Name
				params["deleted_certificate"] = hex.EncodeToString(m.DeleteCertificate)
			}
		case *types.MsgRenewDomain:
			if id, err := st.RenewDomain(ctx, m, height); err != nil {
				return errors.Wrapf(err, "renew domain message, domain name: %s", m.Domain)
			} else {
				accountID = id
				params["action"] = "renew_domain"
				params["sender"] = m.Signer.String()
				params["domain_name"] = m.Domain
			}
		case *types.MsgRenewAccount:
			if id, err := st.RenewAccount(ctx, m, height); err != nil {
				return errors.Wrapf(err, "renew account message, domain name: %s, account name: %s", m.Domain, m.Name)
			} else {
				accountID = id
				params["action"] = "renew_account"
				params["sender"] = m.Signer.String()
				params["domain_name"] = m.Domain
				params["account_name"] = m.Name
			}
		}

		if len(params) > 0 {
			queries = append(queries, &LcdRequestData{AccountID: accountID, Params: params})
		}
	}

	if len(queries) > 0 {
		if responses, err := FetchLcdData(ctx, urlLCD, &queries, height); err != nil {
			return errors.Wrapf(err, "FetchLcdData() failed")
		} else if err = st.HandleLcdData(ctx, &queries, responses, height, denom); err != nil {
			return errors.Wrapf(err, "HandleLcdData() failed")
		}
	}

	return nil
}

package pkg

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/prometheus/common/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/types"

	"github.com/pkg/errors"
)

const syncRetryTimeout = 3 * time.Second

// Sync uploads to local store all blocks that are not present yet, starting
// with the blocks with the lowest hight first. It always returns the number of
// blocks inserted, even if returning an error.
func Sync(ctx context.Context, tmc *TendermintClient, st *Store, hrp string) (uint, error) {
	var (
		inserted        uint
		syncedHeight    int64
		lastKnownHeight int64
	)

	switch block, err := st.LatestBlock(ctx); {
	case ErrNotFound.Is(err):
		syncedHeight = 0
	case err == nil:
		syncedHeight = block.Height
	default:
		return inserted, errors.Wrap(err, "latest block")
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
		for _, tx := range tmblock.Transactions {
			coins := tx.Fee.Amount
			for _, c := range coins {
				if c.Denom != hrp {
					return 1, errors.Wrapf(ErrDenom, "not supported denom: %s, expected %s", c.Denom, hrp)
				}
				fee = fee.Add(c.Amount)
			}

			if err := routeMsgs(ctx, st, tx.Msgs); err != nil {
				log.Error(errors.Wrapf(err, "height", c.Height))
			}
		}

		block := Block{
			Height:  c.Height,
			Hash:    hex.EncodeToString(c.Hash),
			Time:    c.Time.UTC(),
			FeeFrac: fee.Uint64(),
		}
		if err := st.InsertBlock(ctx, block); err != nil {
			return inserted, errors.Wrapf(err, "insert block %d", c.Height)
		}
		inserted++
	}
}

// Domain/Account valid until field is skipped, maybe could be implemented via
// extra query calls on specific height
func routeMsgs(ctx context.Context, st *Store, msgs []sdk.Msg) error {
	for _, msg := range msgs {
		switch m := msg.(type) {
		case *types.MsgRegisterDomain:
			if _, err := st.RegisterDomain(ctx, m); err != nil {
				return errors.Wrap(err, "register domain message")
			}
		case *types.MsgDeleteDomain:
			if err := st.DeleteDomain(ctx, m); err != nil {
				return errors.Wrapf(err, "delete domain message, domain name: %s", m.Domain)
			}
		case *types.MsgTransferDomain:
			if err := st.TransferDomain(ctx, m); err != nil {
				return errors.Wrapf(err, "transfer domain message, domain name: %s", m.Domain)
			}
		case *types.MsgRegisterAccount:
			if _, err := st.RegisterAccount(ctx, m); err != nil {
				return errors.Wrapf(err, "register account message, domain name: %s, account name: %s", m.Domain, m.Name)
			}
		case *types.MsgDeleteAccount:
			if err := st.DeleteAccount(ctx, m); err != nil {
				return errors.Wrapf(err, "delete account message, domain name: %s, account name: %s", m.Domain, m.Name)
			}
		case *types.MsgReplaceAccountResources:
			if _, err := st.ReplaceAccountResources(ctx, m); err != nil {
				return errors.Wrapf(err, "replace account resources msg, domain name: %s, account name: %s", m.Domain, m.Name)
			}
		}
	}
	return nil
}

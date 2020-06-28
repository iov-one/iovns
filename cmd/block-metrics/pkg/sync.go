package pkg

import (
	"context"
	"time"

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

			for _, msg := range tx.Msgs {
				switch m := msg.(type) {
				case *types.MsgRegisterDomain:
					if _, err := st.RegisterDomain(ctx, m); err != nil {
						return inserted, errors.Wrapf(err, "register domain message %d", c.Height)
					}
				}
			}
		}

		inserted++
	}
}

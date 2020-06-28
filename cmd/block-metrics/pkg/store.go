package pkg

import (
	"context"
	"database/sql"
	"time"

	"github.com/iov-one/iovns/x/domain/types"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// NewStore returns a store that provides an access to our database.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

type Store struct {
	db *sql.DB
}

type Block struct {
	Height       int64        `json:"height"`
	Hash         string       `json:"hash"`
	Time         time.Time    `json:"time"`
	ProposerID   int64        `json:"-"`
	ProposerName string       `json:"proposer_name"`
	MissingIDs   []int64      `json:"-"`
	Messages     []string     `json:"messages,omitempty"`
	FeeFrac      uint64       `json:"fee_frac"`
	Transactions []auth.StdTx `json:"transactions"`
}

func (s *Store) RegisterDomain(ctx context.Context, msg *types.MsgRegisterDomain) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO domains (name, admin, type, broker, fee_payer_addr)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, msg.Name, msg.Admin, msg.DomainType, msg.Broker, msg.FeePayerAddr).Scan(&id)
	return id, castPgErr(err)
}

// LatestBlock returns the block with the greatest high value. This method
// returns ErrNotFound if no block exist.
// Note that it doesn't load the validators by default
func (s *Store) LatestBlock(ctx context.Context) (*Block, error) {
	blocks, err := s.LastNBlock(ctx, 1, 0)
	if err != nil {
		return nil, err
	}
	return blocks[0], nil
}

// LoadLastNBlock returns the last blocks with given count.
// ErrNotFound is returned if no blocks exist.
// ErrLimit is returned if allowed limit is exceeded
// Note that it doesn't load the validators by default
func (s *Store) LastNBlock(ctx context.Context, limit, after int) ([]*Block, error) {
	// max number of blocks that is allowed to retrieved is 100
	if limit > 100 {
		return nil, errors.Wrapf(ErrLimit, "limit exceeded")
	}

	var rows *sql.Rows
	var err error
	if after == 0 {
		rows, err = s.db.QueryContext(ctx, `
		SELECT block_height, block_hash, block_time, proposer_id, messages, fee_frac
		FROM blocks
		ORDER BY block_height DESC
		LIMIT $1
	`, limit)
	} else {
		rows, err = s.db.QueryContext(ctx, `
		SELECT block_height, block_hash, block_time, proposer_id, messages, fee_frac
		FROM blocks
		WHERE block_height < $1
		ORDER BY block_height DESC
		LIMIT $2
	`, after, limit)
	}
	defer rows.Close()

	if err != nil {
		err = castPgErr(err)
		if ErrNotFound.Is(err) {
			return nil, errors.Wrap(err, "no blocks")
		}
		return nil, errors.Wrap(castPgErr(err), "cannot select block")
	}

	var blocks []*Block

	for rows.Next() {
		var b Block
		err := rows.Scan(&b.Height, &b.Hash, &b.Time, &b.ProposerID, pq.Array(&b.Messages), &b.FeeFrac)
		if err != nil {
			err = castPgErr(err)
			if ErrNotFound.Is(err) {
				return nil, errors.Wrap(err, "no blocks")
			}
			return nil, errors.Wrap(castPgErr(err), "cannot select block")

		}
		txs, err := s.LoadTxsInBlock(ctx, b.Height)
		if err != nil && !ErrNotFound.Is(err) {
			return nil, err
		}
		b.Transactions = txs

		// normalize it here, as not always stored like this in the db
		b.Time = b.Time.UTC()
		blocks = append(blocks, &b)
	}
	if len(blocks) == 0 {
		return nil, errors.Wrap(ErrNotFound, "no blocks")
	}
	return blocks, nil
}

func (s *Store) LoadTxsInBlock(ctx context.Context, blockHeight int64) ([]auth.StdTx, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT messages, fee, signatures, memo
		FROM transactions
		WHERE block_id=$1
	`, blockHeight)
	defer rows.Close()

	if err != nil {
		err = castPgErr(err)
		if ErrNotFound.Is(err) {
			return nil, errors.Wrap(err, "no txs")
		}
		return nil, errors.Wrap(castPgErr(err), "cannot select txs")
	}

	var txs []auth.StdTx

	for rows.Next() {
		var tx auth.StdTx
		err := rows.Scan(&tx.Msgs, &tx.Fee, &tx.Signatures, &tx.Memo)
		if err != nil {
			err = castPgErr(err)
			if ErrNotFound.Is(err) {
				return nil, errors.Wrap(err, "no tx")
			}
			return nil, errors.Wrap(castPgErr(err), "cannot select tx")
		}
		txs = append(txs, tx)
	}

	if len(txs) == 0 {
		return nil, errors.Wrap(ErrNotFound, "no txs")
	}

	return txs, nil
}

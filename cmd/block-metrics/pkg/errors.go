package pkg

import (
	"database/sql"

	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/lib/pq"
)

const ModuleName = "block-metrics"

var (
	ErrNotImplemented = errors.Register(ModuleName, 1, "not implemented")
	ErrFailedResponse = errors.Register(ModuleName, 2, "failed response")
	ErrConflict       = errors.Register(ModuleName, 3, "conflict")
	ErrLimit          = errors.Register(ModuleName, 4, "limit")
	ErrNotFound       = errors.Register(ModuleName, 5, "not found")
	ErrDenom          = errors.Register(ModuleName, 6, "denomination not supported")
)

func wrapPgErr(err error, msg string) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(castPgErr(err), msg)
}

func castPgErr(err error) error {
	if err == nil {
		return nil
	}

	if err == sql.ErrNoRows {
		return ErrNotFound
	}

	if e, ok := err.(*pq.Error); ok {
		switch prefix := e.Code[:2]; prefix {
		case "20":
			return errors.Wrap(ErrNotFound, e.Message)
		case "23":
			return errors.Wrap(ErrConflict, e.Message)
		}
		err = errors.Wrap(err, string(e.Code))
	}

	return err
}

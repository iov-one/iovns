package pkg

import (
	"encoding/hex"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
)

// Tentermint is using hex encoded binary data. Provide a type that will do the
// conversion as a part of JSON unmarshaling.
type hexstring []byte

func (h *hexstring) UnmarshalJSON(raw []byte) error {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return errors.Wrap(err, "invalid JSON string")
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return errors.Wrap(err, "hex decode")
	}
	*h = b
	return nil
}

// Tentermint is using strings where a number is expected. Provide a type that
// will do the conversion as a part of JSON unmarshaling.
type sint64 int64

func (i sint64) Int64() int64 {
	return int64(i)
}

func (i *sint64) UnmarshalJSON(raw []byte) error {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return errors.Wrap(err, "invalid JSON string")
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid number")
	}
	*i = sint64(n)
	return nil
}

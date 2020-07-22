package queries

import (
	"encoding/json"
)

// QueryEncoder defines a function that encodes query models to bytes
type QueryEncoder func(queryModel interface{}) ([]byte, error)

// QueryDecoder defines a function that decodes query bytes to query models
type QueryDecoder func(data []byte, ptrTargetModel interface{}) error

// DefaultQueryEncode is the default function used
// to marshal query models into bytes
var DefaultQueryEncode QueryEncoder = json.Marshal

// DefaultQueryDecode is the default function used to
// decode query bytes to query models
var DefaultQueryDecode QueryDecoder = json.Unmarshal

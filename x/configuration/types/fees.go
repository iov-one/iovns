package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dt "github.com/iov-one/iovns/x/domain/types"
)

// LevelFeeMapper maps fees based on level
type LevelFeeMapper map[string]sdk.Coin

// MarshalJSON marshals the map in a deterministic way
func (m LevelFeeMapper) MarshalJSON() ([]byte, error) {
	// golang marshals deterministically
	// maps keys are ordered and structs
	// follow order of their fields

	// use this subtype to make sure the
	// order will be the same even in case
	// of changes on the type from cosmos-sdk
	type coin struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	}

	jsonMap := make(map[string]coin, len(m))
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	for _, key := range keys {
		c := m[key]
		jsonMap[key] = coin{
			Denom:  c.Denom,
			Amount: c.Amount.String(),
		}
	}
	result, err := json.Marshal(jsonMap)
	if err != nil {
		panic(err)
	}
	return result, nil
}

func (m *LevelFeeMapper) UnmarshalJSON(b []byte) error {
	// make map if it is has not been initialized
	if *m == nil {
		*m = make(LevelFeeMapper)
	}
	// use this subtype to make sure the
	// order will be the same even in case
	// of changes on the type from cosmos-sdk
	type coin struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	}
	var x map[string]coin
	err := json.Unmarshal(b, &x)
	if err != nil {
		return err
	}
	for k, v := range x {
		sdkInt, ok := sdk.NewIntFromString(v.Amount)
		if !ok {
			return fmt.Errorf("invalid sdk.Int: %s", v.Amount)
		}
		(*m)[k] = sdk.NewCoin(v.Denom, sdkInt)
	}
	return nil
}

// LevelFees contains different type of fees
// to calculate coins to detract when
// processing different messages
type Fees struct {
	// LevelFees maps msg fees to their level
	LevelFees map[dt.MsgUniqueID]LevelFeeMapper
	// DefaultFees maps the default fees for a msg
	DefaultFees map[dt.MsgUniqueID]sdk.Coin
}

// MarshalJSON makes sure the map is ordered deterministically
func (f *Fees) MarshalJSON() ([]byte, error) {
	// do not edit this or
	// there will be undeterministic
	// behaviour with the current state
	type coin struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	}
	type fee struct {
		LevelFees   map[dt.MsgUniqueID]LevelFeeMapper `json:"level_fees"`
		DefaultFees map[dt.MsgUniqueID]coin           `json:"default_fees"`
	}
	var x = fee{
		LevelFees:   f.LevelFees,
		DefaultFees: make(map[dt.MsgUniqueID]coin, len(f.DefaultFees)),
	}
	for k, v := range f.DefaultFees {
		x.DefaultFees[k] = coin{
			Denom:  v.Denom,
			Amount: v.Amount.String(),
		}
	}
	return json.Marshal(x)
}

func (f *Fees) UnmarshalJSON(b []byte) error {
	// init fees if nil
	if f == nil {
		*f = Fees{}
	}
	// init maps if nil
	if f.DefaultFees == nil {
		f.DefaultFees = make(map[dt.MsgUniqueID]sdk.Coin)
	}
	if f.LevelFees == nil {
		f.LevelFees = make(map[dt.MsgUniqueID]LevelFeeMapper)
	}
	// re-use types used for marshalling
	type coin struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	}
	var x struct {
		DefaultFees map[string]coin           `json:"default_fees"`
		LevelFees   map[string]LevelFeeMapper `json:"level_fees"`
	}
	x.DefaultFees = make(map[string]coin)
	x.LevelFees = make(map[string]LevelFeeMapper)
	// unmarshal
	err := json.Unmarshal(b, &x)
	if err != nil {
		return err
	}
	// set default fees
	for k, v := range x.DefaultFees {
		sdkInt, ok := sdk.NewIntFromString(v.Amount)
		if !ok {
			return fmt.Errorf("invalid sdk.Int: %s", v.Amount)
		}
		f.DefaultFees[dt.MsgUniqueID(k)] = sdk.NewCoin(v.Denom, sdkInt)
	}
	for k, v := range x.LevelFees {

		f.LevelFees[dt.MsgUniqueID(k)] = v
	}
	return nil
}

// NewFees is Fees constructor
func NewFees() *Fees {
	return &Fees{
		LevelFees:   make(map[dt.MsgUniqueID]LevelFeeMapper),
		DefaultFees: make(map[dt.MsgUniqueID]sdk.Coin),
	}
}

// CalculateLevelFees calculates fees based on message type and level
// if there is no level fee then it retreats to the default fees for msg
// false is returned only in the case in which no fee was found or can be applied.
func (f *Fees) CalculateLevelFees(msg dt.Feeable, level int) (sdk.Coin, bool) {
	sdkIntLevel := sdk.NewInt(int64(level))
	msgID := msg.ID()
	// get fees per message type
	msgFees, ok := f.LevelFees[msgID]
	// if fees based on sdkIntLevel are not found
	// return the default fee
	if !ok {
		// if the fee was not found then
		// apply the default fees for the msg
		fee, ok := f.DefaultFees[msgID]
		if !ok {
			// if not found return nothing
			return sdk.Coin{}, false
		}
		// if found return the default fee
		return fee, true
	}
	// get fees based on sdkIntLevel
	fee, ok := msgFees[sdkIntLevel.String()]
	if !ok {
		// if not found return the default level fee
		defaultFee, ok := f.DefaultFees[msgID]
		if !ok {
			// no fees found
			return sdk.Coin{}, false
		}
		// return default fee
		return defaultFee, true
	}
	// return fee
	return fee, true
}

// UpsertLevelFees updates or sets the level fees for the message
func (f *Fees) UpsertLevelFees(msg dt.Feeable, level int, coin sdk.Coin) {
	sdkIntLevel := sdk.NewInt(int64(level))
	msgID := msg.ID()
	feesMap, ok := f.LevelFees[msgID]
	// if fee map for that msg type does not exist create it
	if !ok {
		f.LevelFees[msgID] = make(LevelFeeMapper)
		feesMap = f.LevelFees[msgID]
	}
	// update fees
	feesMap[sdkIntLevel.String()] = coin
}

// UpsertDefaultFees updates or sets the default fees for sdk.Msg
func (f *Fees) UpsertDefaultFees(msg dt.Feeable, coin sdk.Coin) {
	f.DefaultFees[msg.ID()] = coin
}

func (f *Fees) DeleteLevelFee(msg dt.Feeable, level int) {
	sdkIntLevel := sdk.NewInt(int64(level))
	feeMap, ok := f.LevelFees[msg.ID()]
	if !ok {
		return
	}
	delete(feeMap, sdkIntLevel.String())
}

func (f *Fees) DeleteDefaultFee(msg dt.Feeable) {
	delete(f.DefaultFees, msg.ID())
}

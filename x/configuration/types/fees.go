package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// msgUniqueID exists to make sure
// that sdk.Msg are parsed into unique IDs
type msgUniqueID string

// LengthFeeMapper maps fees based on length
type LengthFeeMapper map[int]sdk.Coin

// LengthFees contains different type of fees
// to calculate coins to detract when
// processing different messages
type Fees struct {
	// LengthFees maps msg fees to their length
	LengthFees map[msgUniqueID]LengthFeeMapper
	// DefaultFees maps the default fees for a msg
	DefaultFees map[msgUniqueID]sdk.Coin
}

// NewFees is Fees constructor
func NewFees() *Fees {
	return &Fees{
		LengthFees:  make(map[msgUniqueID]LengthFeeMapper),
		DefaultFees: make(map[msgUniqueID]sdk.Coin),
	}
}

// CalculateLengthFees calculates fees based on message type and length
// if there is no length fee then it retreats to the default fees for msg
// false is returned only in the case in which no fee was found or can be applied.
func (f *Fees) CalculateLengthFees(msg sdk.Msg, length int) (sdk.Coin, bool) {
	msgID := f.getMsgID(msg)
	// get fees per message type
	msgFees, ok := f.LengthFees[msgID]
	// if fees based on length are not found
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
	// get fees based on length
	fee, ok := msgFees[length]
	if !ok {
		// if not found return the default length fee
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

// getMsgID returns the unique id for the message to apply fees on
func (f *Fees) getMsgID(msg sdk.Msg) msgUniqueID {
	return msgUniqueID(fmt.Sprintf("%s/%s", msg.Route(), msg.Type()))
}

// UpsertLengthFees updates or sets the length fees for the message
func (f *Fees) UpsertLengthFees(msg sdk.Msg, length int, coin sdk.Coin) {
	msgID := f.getMsgID(msg)
	feesMap, ok := f.LengthFees[msgID]
	// if fee map for that msg type does not exist create it
	if !ok {
		f.LengthFees[msgID] = make(LengthFeeMapper)
		feesMap = f.LengthFees[msgID]
	}
	// update fees
	feesMap[length] = coin
}

// UpsertDefaultFees updates or sets the default fees for sdk.Msg
func (f *Fees) UpsertDefaultFees(msg sdk.Msg, coin sdk.Coin) {
	f.DefaultFees[f.getMsgID(msg)] = coin
}

func (f *Fees) DeleteLengthFee(msg sdk.Msg, length int) {
	feeMap, ok := f.LengthFees[f.getMsgID(msg)]
	if !ok {
		return
	}
	delete(feeMap, length)
}

func (f *Fees) DeleteDefaultFee(msg sdk.Msg) {
	delete(f.DefaultFees, f.getMsgID(msg))
}

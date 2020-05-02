package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
)

func TestFeesCases(t *testing.T) {
	t.Run("cases", func(t *testing.T) {
		f := NewFees()
		_, ok := f.CalculateLengthFees(MsgUpdateConfig{}, 0)
		// no fees expected
		if ok {
			t.Fatalf("CalculateLengthFees() unexpected result: %v", ok)
		}
		// get default fees
		coinFee := sdk.NewCoin("test", sdk.NewInt(10))
		f.UpsertDefaultFees(MsgUpdateConfig{}, coinFee)
		res, ok := f.CalculateLengthFees(MsgUpdateConfig{}, 0)
		if !ok {
			t.Fatalf("CalculateLengthFees() result expected")
		}
		if !res.IsEqual(coinFee) {
			t.Fatalf("CalculateLengthFees() wanted: %s, got: %s", coinFee, res)
		}
		// get level fees
		levelFee := sdk.NewCoin("test", sdk.NewInt(15))
		f.UpsertLengthFees(MsgUpdateConfig{}, 10, levelFee)
		res, ok = f.CalculateLengthFees(MsgUpdateConfig{}, 10)
		if !ok {
			t.Fatalf("CalculateLengthFees() result expected")
		}
		if !res.IsEqual(levelFee) {
			t.Fatalf("CalculatELengthFees() wanted: %s, got: %s", levelFee, res)
		}
		// get default fee because level does not exist
		res, ok = f.CalculateLengthFees(MsgUpdateConfig{}, 11)
		if !ok {
			t.Fatalf("CalculateLengthFees() result expected")
		}
		if !res.IsEqual(coinFee) {
			t.Fatalf("CalculateLengthFees() want: %s, got: %s", coinFee, res)
		}
	})
}

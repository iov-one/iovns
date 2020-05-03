package types

import (
	"bytes"
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/types"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"time"
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

func TestLengthFeeMapper_MarshalUnmarshal(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	x := make(LengthFeeMapper)
	for i := 0; i < 500; i++ {
		x[strconv.Itoa(rand.Int())] = sdk.NewCoin("idk", sdk.NewInt(int64(rand.Int())))
	}
	b, err := json.Marshal(x)
	// test deterministic marshalling
	for i := 0; i < 500; i++ {
		x, err := json.Marshal(x)
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
		if !bytes.Equal(x, b) {
			t.Fatalf("undeterministic behaviour")
		}
	}
	if err != nil {
		t.Fatalf("got error: %s", err)
	}
	var y = make(LengthFeeMapper)
	err = json.Unmarshal(b, &y)
	if err != nil {
		t.Fatalf("got error: %s", err)
	}
	if !reflect.DeepEqual(y, x) {
		t.Fatalf("results do not match")
	}
}

func TestFees_MarshalUnmarshalJSON(t *testing.T) {
	fees := NewFees()
	// insert default fees
	fees.UpsertDefaultFees(MsgUpdateConfig{}, sdk.NewCoin("test", sdk.NewInt(2)))
	fees.UpsertDefaultFees(&types.MsgRegisterDomain{}, sdk.NewCoin("test", sdk.NewInt(1)))
	// insert length fees
	fees.UpsertLengthFees(MsgUpdateConfig{}, 10, sdk.NewCoin("test", sdk.NewInt(10)))
	fees.UpsertLengthFees(&types.MsgRegisterDomain{}, 4, sdk.NewCoin("test", sdk.NewInt(15)))
	x, err := json.Marshal(fees)
	if err != nil {
		t.Fatalf("got error: %s", err)
	}
	var y = new(Fees)
	err = json.Unmarshal(x, y)
	if err != nil {
		t.Fatalf("got error: %s", err)
	}
	// check equality
	if !reflect.DeepEqual(fees, y) {
		t.Fatalf("results mismatch: got %+v, want: %+v", y, fees)
	}
}

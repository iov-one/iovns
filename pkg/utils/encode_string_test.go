package utils

import (
	"math"
	"testing"
)

func TestEncode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		x := "hilokkkkkkkkowekdfokewfokew?WER?=£$)=£\")R=\"WEKDFSFDSkkkkkkkkk"
		b := Base64Encode(x)
		resp, err := Base64Decode(b)
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
		if resp != x {
			t.Fatalf("decoding mismatch: want: %s, got: %s", x, b)
		}
	})
	t.Logf("%s", []byte{uint8(math.MaxUint8)})
}

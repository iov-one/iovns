package types

import (
	"bytes"
	"testing"
)

func TestNewPrimaryKey(t *testing.T) {
	e := []byte("test")
	x := NewPrimaryKey(e)
	if !bytes.Equal(x.Key(), e) {
		t.Fatal("mismatch")
	}
}

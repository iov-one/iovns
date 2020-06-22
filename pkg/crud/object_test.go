package crud

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

type testObject struct{}

func (t testObject) PrimaryKey() PrimaryKey {
	return []byte{0x0}
}

func (t testObject) SecondaryKeys() []SecondaryKey {
	return []SecondaryKey{
		{
			Key:         []byte{0x1},
			StorePrefix: []byte{0x1},
		},
	}
}

type TestIndex struct{}

func (t TestIndex) SecondaryKey() SecondaryKey {
	return SecondaryKey{StorePrefix: []byte{0x1}, Key: []byte("key")}
}

func Benchmark_inspect(b *testing.B) {
	type cmplx struct {
		PK  string `crud:"primaryKey"`
		SK1 string `crud:"secondaryKey,01"`
		SK2 string `crud:"secondaryKey,02"`
		SK3 string `crud:"secondaryKey,03"`
		SK4 string `crud:"secondaryKey,04"`
		SK5 string `crud:"secondaryKey,05"`
	}
	x := &cmplx{
		PK:  "1",
		SK1: "2",
		SK2: "3",
		SK3: "4",
		SK4: "5",
		SK5: "6",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := inspect(x)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_getKeys(b *testing.B) {
	type cmplx struct {
		PK  string `crud:"primaryKey"`
		SK1 string `crud:"secondaryKey,01"`
		SK2 string `crud:"secondaryKey,02"`
		SK3 string `crud:"secondaryKey,03"`
		SK4 string `crud:"secondaryKey,04"`
		SK5 string `crud:"secondaryKey,05"`
	}
	x := &cmplx{
		PK:  "1",
		SK1: "2",
		SK2: "3",
		SK3: "4",
		SK4: "5",
		SK5: "6",
	}
	v := reflect.ValueOf(x).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := getKeys(v)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Test_inspect(t *testing.T) {
	t.Run("not a pointer", func(t *testing.T) {
		_, _, err := inspect(0)
		if !errors.Is(err, errNotAPointer) {
			t.Fatalf("unexpected error: %s", err)
		}
	})
	t.Run("not a struct", func(t *testing.T) {
		_, _, err := inspect(new(int))
		if !errors.Is(err, errPointerToStruct) {
			t.Fatalf("unexpected error: %s", err)
		}
	})
	t.Run("implements object", func(t *testing.T) {
		obj := &testObject{}
		pk, sk, err := inspect(obj)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(pk, obj.PrimaryKey()) {
			t.Fatal("unexpected primary key")
		}
		if !reflect.DeepEqual(sk, obj.SecondaryKeys()) {
			t.Fatal("unexpected secondary key")
		}
	})
}
func Test_marshalValue(t *testing.T) {
	x := "hello"
	b, err := marshalValue(reflect.ValueOf(x))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%x", b)
}

func Test_marshalSlice(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		x := []byte{0x01}
		got, err := marshalSlice(reflect.ValueOf(x))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, x) {
			t.Fatal("unexpected result")
		}
	})
	t.Run("invalid type", func(t *testing.T) {
		x := []string{"hi"}
		_, err := marshalSlice(reflect.ValueOf(x))
		if err == nil {
			t.Fatal("expected error on invalid type")
		}

	})
}

func Test_getKeys(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		type obj struct {
			PK string `crud:"primaryKey"`
			SK string `crud:"secondaryKey,01"`
		}
		pk, sk, err := getKeys(reflect.ValueOf(obj{
			PK: "test1",
			SK: "test2",
		}))
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%#v %#v", pk, sk)
	})
	t.Run("multiple primary keys", func(t *testing.T) {
		type obj struct {
			PK  string `crud:"primaryKey"`
			PK2 string `crud:"primaryKey"`
		}
		_, _, err := getKeys(reflect.ValueOf(obj{
			PK:  "test1",
			PK2: "test2",
		}))
		if !errors.Is(err, errMultiplePrimaryKeys) {
			t.Fatalf("unexpected error: %s", err)
		}
	})
	t.Run("no primary key", func(t *testing.T) {
		type obj struct {
			SK  string `crud:"secondaryKey,02"`
			SK2 string `crud:"secondaryKey,01"`
		}
		_, _, err := getKeys(reflect.ValueOf(obj{
			SK:  "test1",
			SK2: "test2",
		}))
		if !errors.Is(err, errNoPrimaryKey) {
			t.Fatalf("unexpected error: %s", err)
		}
	})
	t.Run("invalid hex", func(t *testing.T) {
		type obj struct {
			PK string `crud:"primaryKey"`
			SK string `crud:"secondaryKey,0x1"`
		}
		_, _, err := getKeys(reflect.ValueOf(obj{
			PK: "test1",
			SK: "test2",
		}))
		if err == nil {
			t.Fatal("error expected")
		}
	})
	t.Run("with index type", func(t *testing.T) {
		type obj struct {
			PK string `crud:"primaryKey"`
			SK TestIndex
		}
		pk, sk, err := getKeys(reflect.ValueOf(&obj{
			PK: "pk",
			SK: TestIndex{},
		}).Elem())
		if err != nil {
			t.Fatal(err)
		}
		t.Log()
		t.Logf("%#v, %#v", pk, sk)
	})
	t.Run("with index slice type", func(t *testing.T) {
		type obj struct {
			PK string `crud:"primaryKey"`
			SK []TestIndex
		}
		pk, sk, err := getKeys(reflect.ValueOf(&obj{
			PK: "pk",
			SK: []TestIndex{
				{}, {},
			},
		}).Elem())
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%#v, %#v", pk, sk)
	})
}

func Test_isValidSecondaryKey(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		err := isValidSecondaryKey(SecondaryKey{
			Key:         []byte("not empty"),
			StorePrefix: []byte{0x1},
		})
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("invalid secondary key prefix", func(t *testing.T) {
		err := isValidSecondaryKey(SecondaryKey{
			Key:         nil,
			StorePrefix: []byte{PrimaryKeyPrefix},
		})
		if !errors.Is(err, errIsPrimaryKeyPrefix) {
			t.Fatalf("unexpected error: %s", err)
		}
	})
	t.Run("invalid secondary key key", func(t *testing.T) {
		err := isValidSecondaryKey(SecondaryKey{
			Key:         nil,
			StorePrefix: []byte{0x1},
		})
		if !errors.Is(err, errEmptyKey) {
			t.Fatalf("unexpected error: %s", err)
		}
	})
}

func Test_validateKeys(t *testing.T) {
	t.Run("empty primary key", func(t *testing.T) {
		err := validateKeys(nil, nil)
		if !errors.Is(err, errEmptyKey) {
			t.Fatalf("unexpected error: %s", err)
		}
	})
	t.Run("duplicate secondary key", func(t *testing.T) {
		err := validateKeys([]byte("valid"), []SecondaryKey{
			{
				Key:         []byte("same"),
				StorePrefix: []byte("pfx"),
			},
			{
				Key:         []byte("same"),
				StorePrefix: []byte("pfx"),
			},
		})
		if !errors.Is(err, errDuplicateKey) {
			t.Fatalf("unexpected error: %s", err)
		}
	})
}

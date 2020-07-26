package types

// PrimaryKey defines a primary key, which is a secondary key, under the hood, but with a fixed 0x0 prefix
type PrimaryKey interface {
	// KeyCopy makes a copy of the primary key
	KeyCopy() []byte
	// Key returns the primary key
	Key() []byte
}

type primaryKey []byte

func (p primaryKey) KeyCopy() []byte {
	cpy := make([]byte, len(p))
	copy(cpy, p)
	return cpy
}

func (p primaryKey) Key() []byte {
	return p
}


func NewPrimaryKey(value []byte) PrimaryKey {
	v := make(primaryKey, len(value))
	copy(v, value)
	return v
}

func NewPrimaryKeyFromString(s string) PrimaryKey {
	return primaryKey(s)
}


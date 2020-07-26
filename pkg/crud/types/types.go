package types

// Object defines an object in which we can do crud operations
type Object interface {
	// PrimaryKey returns the unique key of the object
	PrimaryKey() PrimaryKey
	// SecondaryKeys returns the secondary keys used to index the object
	SecondaryKeys() []SecondaryKey
}

// Store defines the crud store behaviour
// interface
type Store interface {
	// Create creates a new object and its indexes
	Create(o Object)
	// Read reads to the given target using the provided primary key
	// returns false if no object exists
	Read(key PrimaryKey, target Object) bool
	// Update updates the object and its indexes, it will panic if the object does not exist
	Update(o Object)
	// Delete deletes the object given its primary key and remove its indexes,
	// it will panic if the provided primary key does not exist in the kv store
	Delete(key PrimaryKey)
	// Filter returns a filter given an object whose fields are filters
	Filter(filter Object) Filter
	IterateKeys(func (pk PrimaryKey) bool)
}

// Filter defines the behaviour of the store filter
type Filter interface {
	Read(target Object)
	Valid() bool
	Next()
	Delete()
}
package store

import (
	"github.com/iov-one/iovns/pkg/crud/types"
	"sort"
)

type hash string

// TODO this is an inefficient implementation that should be changed asap.
type keySet map[hash]struct{}

func (k keySet) Insert(b types.PrimaryKey) {
	key := b.Key()
	k[hash(key)] = struct{}{}
}

func (k keySet) Has(b types.PrimaryKey) bool {
	key := b.Key()
	_, ok := k[hash(key)]
	return ok
}

func (k keySet) Keys() []types.PrimaryKey {
	primaryKeys := make([]types.PrimaryKey, 0, len(k))
	for key := range k {
		primaryKeys = append(primaryKeys, types.NewPrimaryKeyFromString(string(key)))
	}
	return primaryKeys
}

func (k keySet) Len() int {
	return len(k)
}

type set interface {
	Has(b types.PrimaryKey) bool
	Keys() []types.PrimaryKey
	Len() int
}

func primaryKeysFromSets(sets []set) []types.PrimaryKey {
	if len(sets) == 0 {
		return nil
	}
	if len(sets) == 1 {
		return sets[0].Keys()
	}
	// determine the smallest set, as the final filter
	// will have, at best, all the keys in the smallest one
	sort.Slice(sets, func(i, j int) bool {
		return sets[i].Len() < sets[j].Len()
	})
	smallestSet := sets[0]
	smallerLen := smallestSet.Len()
	// if smallest is zero then return nothing
	if smallerLen == 0 {
		return nil
	}
	// nice now start filtering
	primaryKeys := make([]types.PrimaryKey, 0, smallestSet.Len())
	for _, key := range smallestSet.Keys() {
		if !isInAll(key, sets) {
			continue
		}
		primaryKeys = append(primaryKeys, key)
	}
	// success
	return primaryKeys
}

// isInAll verifies that the given primary key is present in all sets
func isInAll(key types.PrimaryKey, sets []set) bool {
	for _, set := range sets {
		if !set.Has(key) {
			return false
		}
	}
	// key is in all sets
	return true
}

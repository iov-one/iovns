package iovns

import "github.com/cosmos/cosmos-sdk/store/types"

// TODO test these functions, unless tested they're unusable
func DoInBatches(store types.KVStore, batchNumber int, do func(key []byte)) {
	doAll := func(keys [][]byte) {
		for _, key := range keys {
			do(key)
		}
	}
	var currKey []byte
	for {
		keys := DoStartToEnd(store, currKey, batchNumber)
		// process keys
		doAll(keys)
		// check if keys are empty or if last key is nil (end of iterator)
		if len(keys) == 0 || keys[len(keys)-1] == nil {
			// we're done
			return
		}
		// if we're not done set last key as current key and keep iterating
		currKey = keys[len(keys)-1]
	}
}

func DoStartToEnd(store types.KVStore, start []byte, times int) [][]byte {
	keys := make([][]byte, 0, times)
	iterator := store.Iterator(start, nil)
	defer iterator.Close()
	for i := 0; i < times; i++ {
		if !iterator.Valid() {
			break
		}
		// set current key
		keys = append(keys, iterator.Key())
		// move to next key
		iterator.Next()
	}
	// return
	return keys
}

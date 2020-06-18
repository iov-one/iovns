// Package index contains utilities to create single key indexes in key-value stores
// An index is formed in the following way: a KVStore is provided, a prefix and an
// Indexer object.
//
// The KVStore is then prefixed using the provided prefix, this helps to identify the
// index from other ones, for example: prefix 0x01 identifies the index of addresses
// to object Items stored in another part of the provided KVStore.
//
// The Indexer serves the purpose of being able to create unique byte keys that
// identify the object of which we want to create relations to other objects.
// For example: we want to create relationships between account addresses
// and objects owned by the addresses. So we generate an unique byte key
// for the address. Example: address 0x00000000000000 creates the unique index
// index = byte("0x00000000000000")
// Index keys are then suffixed with byte 0xFF, if an index key contains the reserved
// separator byte it is base64 encoded.
//
// The suffixing comes necessary because in order to iterate keys in the KVStore we simply
// prefix our base KVStore, since KVStores are trees, if we had the following index keys:
// IndexA  = byte("123")
// IndexB = byte("1234")
// While iterating the prefixed store index <index_identifier><IndexA>
// we would end up iterating over IndexB keys too, as they're part of prefix <index_identifier><IndexA> [ + IndexB[3] ]
// With the reserved separator byte, we put an end to the index key and make sure that, while iterating, shorter index keys
// contained in other larger index keys do not end up being mixed together.
// For this same reason we can not have index keys containing the separator, and as a drastic measure keys containing it
// are base64 encoded to guarantee the isolation and safety of the index.
//
// SecondaryKeys save values, those values are identified as Indexed, because they have the ability to turn themselves into
// unique byte keys, the package index provides some basic packing functionality.
// The mentioned values must also possess the ability to uniquely identify themselves back from this unique key.
package index

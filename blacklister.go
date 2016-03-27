package catego

import "github.com/Workiva/go-datastructures/bitarray"

// Blacklister is the structure that will answer
// if an ID is banned or not. It use a BitArray to be as
// fast as possible
type Blacklister struct {
	store bitarray.BitArray
}

// Is return if the given id is banned or not
func (b *Blacklister) Is(id ID) bool {
	var ok bool
	var err error
	ok, err = b.store.GetBit(uint64(id))
	if err != nil {
		return false
	}
	return ok
}

// GetStorage returns the underlying bitarray, useful to do some intersection, and or with the blacklisted node
func (b *Blacklister) GetStorage() bitarray.BitArray {
	return b.store
}

package lib

type Bitset uint64

type ID uint64

type ComponentID ID

type EntityID ID

func (b Bitset) HasID(id ComponentID) bool {
	return b&(1<<id) != 0
}

func (b Bitset) AddID(id ComponentID) Bitset {
	return b | (1 << id)
}

func (b Bitset) RemoveID(id ComponentID) Bitset {
	return b & ^(1 << id)
}

func (b Bitset) Has(bitset Bitset) bool {
	return b&bitset == bitset
}

func (b Bitset) DoesNotHave(bitset Bitset) bool {
	return b&bitset == 0
}

func (b Bitset) Without(bitset Bitset) Bitset {
	return b & ^bitset
}

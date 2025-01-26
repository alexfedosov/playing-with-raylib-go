package lib

type Archetype struct {
	bitset                Bitset
	components            map[ComponentID][]interface{}
	entities              []EntityID
	disabledMaskPerEntity map[EntityID]Bitset
}

func NewArchetype(bitset Bitset, entityCapacity int, componentsCapacity int) *Archetype {
	return &Archetype{
		bitset:                bitset,
		components:            make(map[ComponentID][]interface{}, componentsCapacity),
		entities:              make([]EntityID, entityCapacity),
		disabledMaskPerEntity: make(map[EntityID]Bitset, entityCapacity),
	}
}

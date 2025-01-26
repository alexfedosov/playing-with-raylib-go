package lib

import "fmt"

func (w *World) Log() {
	fmt.Printf("----------------------------------------------------------------")
	for bitset, archetype := range w.archetypes {
		fmt.Printf("Archetype (bitset: %b):\n", bitset)
		for componentID, components := range archetype.components {
			fmt.Printf("  Component ComponentID %d:\n", componentID)
			for _, component := range components {
				fmt.Printf("	%v\n", component)
			}
		}
	}
}

func (w *World) GetEntityCount() int {
	return len(w.entities)
}

package lib

import (
	"testing"
)

type CharacterComponent struct {
	name string
}

type PositionComponent struct {
	x, y float64
}

type IsEnabledComponent struct{}

var _ = RegisterComponent[CharacterComponent]()
var _ = RegisterComponent[PositionComponent]()
var _ = RegisterComponent[IsEnabledComponent]()

func TestWorld_AddComponent(t *testing.T) {
	w := NewWorld()

	for i := 0; i < 10000000; i++ {
		entity := w.CreateEntity()
		w.AddComponents(entity, CharacterComponent{name: "test"}, PositionComponent{x: 10, y: 10}, IsEnabledComponent{})
	}

	w.Query().With(GetComponentID[CharacterComponent]()).Each(func(id EntityID, m map[ComponentID]interface{}) {

	})
}

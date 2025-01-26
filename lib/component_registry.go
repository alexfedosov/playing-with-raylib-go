package lib

import (
	"fmt"
	"reflect"
)

type ComponentRegistry struct {
	nextID   ComponentID
	typeToID map[reflect.Type]ComponentID
}

var Registry = ComponentRegistry{
	nextID:   1,
	typeToID: make(map[reflect.Type]ComponentID),
}

func RegisterComponent[T any]() ComponentID {
	var component T
	componentType := reflect.TypeOf(component)

	if id, exists := Registry.typeToID[componentType]; exists {
		return id
	}

	id := Registry.nextID
	Registry.typeToID[componentType] = id
	Registry.nextID++
	fmt.Printf("registered component %v with id %v\n", componentType, id)
	return id
}

func GetComponentID[T any]() ComponentID {
	var component T
	return GetComponentIDOf(component)
}

func GetComponentIDOf[T any](component T) ComponentID {
	componentType := reflect.TypeOf(component)

	id, exists := Registry.typeToID[componentType]
	if !exists {
		panic(fmt.Sprintf("component type %v not registered", componentType))
	}
	return id
}

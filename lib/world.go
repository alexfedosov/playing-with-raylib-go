package lib

type World struct {
	archetypes       map[Bitset]*Archetype
	entityArchetypes map[EntityID]Bitset

	entities     map[EntityID]bool
	nextEntityID EntityID

	systems []System

	queryCache *QueryCache
}

func NewWorld() *World {
	return &World{
		archetypes:       make(map[Bitset]*Archetype),
		entityArchetypes: make(map[EntityID]Bitset),
		entities:         make(map[EntityID]bool),
		systems:          make([]System, 0),
		queryCache:       NewQueryCache(),
	}
}

func (w *World) CreateEntity() EntityID {
	id := w.nextEntityID
	w.nextEntityID++
	w.entities[id] = true
	w.entityArchetypes[id] = 0
	return id
}

func (w *World) DestroyEntity(entity EntityID) {
	bitset, exist := w.entityArchetypes[entity]
	if !exist {
		return
	}
	archetype, exists := w.archetypes[bitset]
	if !exists {
		return
	}

	entityIdx := -1
	for idx, id := range archetype.entities {
		if id == entity {
			entityIdx = idx
			break
		}
	}

	if entityIdx == -1 {
		return
	}

	lastIdx := len(archetype.entities) - 1
	archetype.entities[entityIdx] = archetype.entities[lastIdx]
	archetype.entities = archetype.entities[:lastIdx]
	delete(archetype.disabledMaskPerEntity, entity)

	if len(archetype.entities) == 0 {
		delete(w.archetypes, bitset)
	} else {
		for componentId, components := range archetype.components {
			components[entityIdx] = components[lastIdx]
			archetype.components[componentId] = components[:lastIdx]
		}
	}
	w.archetypes[bitset] = archetype

	delete(w.entityArchetypes, entity)
	delete(w.entities, entity)
}

func (w *World) AddComponents(entity EntityID, components ...interface{}) {
	oldBitset := w.entityArchetypes[entity]
	newBitset := oldBitset
	for _, component := range components {
		componentID := GetComponentIDOf(component)
		newBitset = newBitset.AddID(componentID)
		w.queryCache.Invalidate(componentID)
	}
	if oldBitset == newBitset {
		w.updateEntityComponent(entity, newBitset, components...)
	} else {
		w.moveEntityToArchetype(entity, oldBitset, newBitset, components)
	}
}

func (w *World) DisableComponent(entity EntityID, componentID ComponentID) {
	bitset := w.entityArchetypes[entity]

	disabledMask, exists := w.archetypes[bitset].disabledMaskPerEntity[entity]
	if !exists {
		disabledMask = 0
	}
	w.archetypes[bitset].disabledMaskPerEntity[entity] = disabledMask.AddID(componentID)
	w.queryCache.Invalidate(componentID)
}

func (w *World) EnableComponent(entity EntityID, componentID ComponentID) {
	bitset := w.entityArchetypes[entity]
	disabledMask, exists := w.archetypes[bitset].disabledMaskPerEntity[entity]
	if !exists {
		disabledMask = 0
	}
	w.archetypes[bitset].disabledMaskPerEntity[entity] = disabledMask.RemoveID(componentID)
	w.queryCache.Invalidate(componentID)
}

func (w *World) RemoveComponent(entity EntityID, componentID ComponentID) {
	oldBitset := w.entityArchetypes[entity]
	newBitset := oldBitset.RemoveID(componentID)
	w.moveEntityToArchetype(entity, oldBitset, newBitset, nil)
	w.queryCache.Invalidate(componentID)
}

func (w *World) AddSystem(system System) {
	w.systems = append(w.systems, system)
}

func (w *World) Update(deltaTime float64) {
	for _, system := range w.systems {
		system.Update(w, deltaTime)
	}
}

func (w *World) moveEntityToArchetype(entity EntityID, oldBitset, newBitset Bitset, components []interface{}) {
	oldArchetype, oldExists := w.archetypes[oldBitset]
	newArchetype, exists := w.archetypes[newBitset]
	if !exists {
		newArchetype = NewArchetype(newBitset, 0, len(components))
		w.archetypes[newBitset] = newArchetype
	}

	// If entity was in an old archetype, move its Components
	if oldExists {
		entityIndex := -1
		for i, e := range oldArchetype.entities {
			if e == entity {
				entityIndex = i
				break
			}
		}

		if entityIndex != -1 {
			newArchetype.disabledMaskPerEntity[entity] = oldArchetype.disabledMaskPerEntity[entity]

			// Copy Components to new archetype
			for id, components := range oldArchetype.components {
				if newBitset.HasID(id) { // If component should exist in new archetype
					newArchetype.components[id] = append(
						newArchetype.components[id],
						components[entityIndex],
					)
				}
			}

			// Remove from old archetype
			for id, components := range oldArchetype.components {
				oldArchetype.components[id] = append(
					components[:entityIndex],
					components[entityIndex+1:]...,
				)
			}
			oldArchetype.entities = append(
				oldArchetype.entities[:entityIndex],
				oldArchetype.entities[entityIndex+1:]...,
			)
		}
	}

	// Add new component if provided
	if components != nil {
		for _, component := range components {
			componentID := GetComponentIDOf(component)
			components := newArchetype.components[componentID]
			newArchetype.components[componentID] = append(components, component)
		}
	}

	// Add entity to new archetype
	newArchetype.entities = append(newArchetype.entities, entity)
	w.entityArchetypes[entity] = newBitset
}

func (w *World) updateEntityComponent(entity EntityID, bitset Bitset, components ...interface{}) {
	archetype := w.archetypes[bitset]
	for _, component := range components {
		id := GetComponentIDOf(component)
		for i, e := range archetype.entities {
			if e == entity {
				archetype.components[id][i] = component
			}
		}
	}
}

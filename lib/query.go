package lib

import "fmt"

type Query struct {
	world     *World
	required  Bitset
	forbidden Bitset
}

type QueryEntity struct {
	ID         EntityID
	Components map[ComponentID]interface{}
}

type QueryResult struct {
	world    *World
	Entities []QueryEntity
}

func (w *World) Query() *Query {
	return &Query{world: w}
}

func (q *Query) With(componentIDs ...ComponentID) *Query {
	for _, id := range componentIDs {
		q.required = q.required.AddID(id)
	}
	return q
}

func (q *Query) Without(componentIDs ...ComponentID) *Query {
	for _, id := range componentIDs {
		q.forbidden = q.forbidden.AddID(id)
	}
	return q
}

func (q *Query) CacheKey() QueryCacheKey {
	return QueryCacheKey{required: q.required, forbidden: q.forbidden}
}

func (q *Query) Get() QueryResult {
	cacheKey := q.CacheKey()
	if result := q.world.queryCache.Get(cacheKey); result != nil {
		fmt.Printf("Cache hit for query:\r\n")
		return *result
	}
	entities := make([]QueryEntity, 0)
	for bitset, archetype := range q.world.archetypes {
		if !bitset.Has(q.required) || !bitset.DoesNotHave(q.forbidden) {
			continue
		}

		for entityIndex, entity := range archetype.entities {
			components := make(map[ComponentID]interface{})
			disabledMask, exists := archetype.disabledMaskPerEntity[entity]
			if !exists {
				disabledMask = Bitset(0)
			}
			oneOfTheQueriedComponentsIsDisabled := false
			for componentID, componentArray := range archetype.components {
				componentWasQueried := q.required.HasID(componentID)
				disabled := disabledMask.HasID(componentID)
				if componentWasQueried && disabled {
					oneOfTheQueriedComponentsIsDisabled = true
					break
				}
				if componentWasQueried && len(componentArray) > entityIndex {
					components[componentID] = componentArray[entityIndex]
				}
			}
			if !oneOfTheQueriedComponentsIsDisabled {
				entities = append(entities, QueryEntity{ID: entity, Components: components})
			}
		}
	}
	result := QueryResult{Entities: entities, world: q.world}
	q.world.queryCache.Set(cacheKey, result)
	return result
}

func (q *Query) Each(fn func(EntityID, map[ComponentID]interface{})) {
	result := q.Get()
	for _, entity := range result.Entities {
		fn(entity.ID, entity.Components)
	}
}

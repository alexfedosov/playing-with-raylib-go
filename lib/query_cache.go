package lib

type QueryCacheKey struct {
	required  Bitset
	forbidden Bitset
}

type QueryCache struct {
	cache map[QueryCacheKey]QueryResult
}

func NewQueryCache() *QueryCache {
	return &QueryCache{
		cache: make(map[QueryCacheKey]QueryResult),
	}
}

func (q *QueryCache) Get(key QueryCacheKey) *QueryResult {
	value, exist := q.cache[key]
	if exist {
		return &value
	}
	return nil
}

func (q *QueryCache) Set(key QueryCacheKey, value QueryResult) {
	q.cache[key] = value
}

func (q *QueryCache) Invalidate(componentID ComponentID) {
	for key := range q.cache {
		if key.required.HasID(componentID) || key.forbidden.HasID(componentID) {
			delete(q.cache, key)
		}
	}
}

package cache

import "container/list"

type PigeonCache struct {
	maxMemory  int64
	usedMemory int64
	list       *list.List
	cache      map[string]*list.Element // *list.Element is a point to a node in PigeonCache.list
}

type entity struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxMemory int64) *PigeonCache {
	return &PigeonCache{
		maxMemory:  maxMemory,
		usedMemory: 0,
		list:       list.New(),
		cache:      make(map[string]*list.Element),
	}
}

func (pigeon *PigeonCache) Get(key string) (value Value, ok bool) {
	element, ok := pigeon.cache[key]
	if ok {
		pigeon.list.MoveToFront(element)
		return element.Value.(*entity).value, true
	} else {
		return nil, false
	}
}



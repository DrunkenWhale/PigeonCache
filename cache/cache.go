package cache

import "container/list"

type PigeonCache struct {
	maxMemory  int64
	usedMemory int64
	list       *list.List
	cache      map[string]*list.Element // *list.Element is a point to a node in PigeonCache.list
	feedback   func(key string, value Value)
}

type entity struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxMemory int64, feedback func(string, Value)) *PigeonCache {
	return &PigeonCache{
		maxMemory:  maxMemory,
		usedMemory: 0,
		list:       list.New(),
		cache:      make(map[string]*list.Element),
		feedback:   feedback,
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

func (pigeon *PigeonCache) RemoveOldest() {
	element := pigeon.list.Back()
	if element != nil {
		pigeon.list.Remove(element)
		kv := element.Value.(*entity)
		delete(pigeon.cache, kv.key)
		pigeon.usedMemory -= int64(kv.value.Len()) + int64(len(kv.key))
		if pigeon.feedback != nil {
			pigeon.feedback(kv.key, kv.value)
		}
	}
}


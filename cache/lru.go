package cache

import "container/list"

type PigeonCache struct {
	maxMemory  int64
	usedMemory int64
	list       *list.List
	Cache      map[string]*list.Element // *list.Element is a point to a node in PigeonCache.list
	feedback   func(key string, value Value)
}

type entity struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func NewPigeonCache(maxMemory int64, feedback func(string, Value)) *PigeonCache {
	return &PigeonCache{
		maxMemory:  maxMemory,
		usedMemory: 0,
		list:       list.New(),
		Cache:      make(map[string]*list.Element),
		feedback:   feedback,
	}
}

func (pigeon *PigeonCache) Get(key string) (value Value, ok bool) {
	element, ok := pigeon.Cache[key]
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
		delete(pigeon.Cache, kv.key)
		pigeon.usedMemory -= int64(kv.value.Len()) + int64(len(kv.key))
		if pigeon.feedback != nil {
			pigeon.feedback(kv.key, kv.value)
		}
	}
}

// Len how many element in this Cache list

func (pigeon *PigeonCache) Len() int {
	return pigeon.list.Len()
}

func (pigeon *PigeonCache) Put(key string, value Value) {
	element, ok := pigeon.Cache[key]
	if ok {
		// element exist
		pigeon.list.MoveToFront(element)
		kv := element.Value.(*entity)
		pigeon.usedMemory += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		element := pigeon.list.PushBack(&entity{key: key, value: value})
		pigeon.Cache[key] = element
		pigeon.usedMemory += int64(len(key)) + int64(value.Len())
	}
	for pigeon.maxMemory > 0 && pigeon.maxMemory < pigeon.usedMemory {
		pigeon.RemoveOldest()
	}
}

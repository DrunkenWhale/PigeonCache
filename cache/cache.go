package cache

import "container/list"

type PigeonCache struct {
	maxMemory  int64
	usedMemory int64
	list       *list.List
	cache      map[string]*list.Element // *list.Element is a point to a node in PigeonCache.list
}

func New(maxMemory int64) *PigeonCache {
	return &PigeonCache{
		maxMemory:  maxMemory,
		usedMemory: 0,
		list:       list.New(),
		cache:      make(map[string]*list.Element),
	}
}

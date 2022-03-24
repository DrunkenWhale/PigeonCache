package main

import (
	"PigeonCache/cache"
)

type String string

func (d String) Len() int {
	return len(d)
}


func main() {
	c := cache.NewPigeonCache(1024, nil)
	c.Put("114", String("114514"))
	c.Put("1141", String("114514"))
	c.Put("1142", String("114514"))
	c.Put("1143", String("114514"))
	c.Put("1144", String("114514"))
	c.Put("1145", String("114514"))
	c.Put("1146", String("114514"))
	c.Put("1147", String("114514"))
	c.Put("1148", String("114514"))
	c.Put("1140", String("114514"))
	c.Put("1149", String("114514"))
	c.Put("11465", String("114514"))
	c.Put("11445", String("114514"))

}

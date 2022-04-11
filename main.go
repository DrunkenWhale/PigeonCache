package main

import (
	"PigeonCache/pigeoncache"
	"fmt"
	"log"
)

type String string

func (d String) Len() int {
	return len(d)
}

func main() {
	TestGet()

}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet() {
	loadCounts := make(map[string]int, len(db))
	test := pigeoncache.NewGroup("test", 2<<7, pigeoncache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	for k, _ := range db {
		value, err := test.Get(k)
		if err == nil && value.Len() == 0 {
			fmt.Println("miss, loading now")
			c, _ := test.Get(k)
			fmt.Println(c.String())
		} else {
			fmt.Println(value)
		}
	}

}

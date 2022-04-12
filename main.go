package main

import (
	"PigeonCache/pigeoncache"
	"PigeonCache/pigeoncache/consistenthash"
	"fmt"
	"log"
	"net/http"
)

type String string

func (d String) Len() int {
	return len(d)
}

func main() {
	testBetaVersion()
}

func testBetaVersion() {
	g := pigeoncache.NewGroup("pigeon", 2<<7, pigeoncache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			return nil, nil
		}))
	adds := []string{
		"http://localhost:8001",
		"http://localhost:8002",
		"http://localhost:8003",
	}

	for _, add := range adds {
		go startApiServer(add, g)
	}
	startCacheServer("http://localhost:7070", adds, g)
}

func startApiServer(apiAddr string, pigeon *pigeoncache.Group) {
	http.Handle(apiAddr+"/api", http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			key := request.URL.Query().Get("key")
			view, err := pigeon.Get(key)
			if err != nil {
				http.Error(writer, err.Error(), 500)
				return
			}
			writer.Header().Set("Content-Type", "application/octet-stream")
			_, _ = writer.Write(view.ByteSlice())
		},
	))
	log.Println("api server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func startCacheServer(add string, adds []string, pigeon *pigeoncache.Group) {
	peers := pigeoncache.NewHttpPool(add)
	peers.Set(adds...)
	pigeon.RegisterPeers(peers)
	log.Println("[Pigeon Cache] running at ", add)
	log.Fatal(http.ListenAndServe(add[7:], peers))
}

func testConsistentHash() {
	hash := consistenthash.New(3, nil)
	hash.Add("2", "4", "6")
	fmt.Println(hash.Get("8"))
	hash.Add("7")
	fmt.Println(hash.Get("8"))
}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func testGet() {
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
	address := "localhost:9999"
	peers := pigeoncache.NewHttpPool(address)
	log.Println("pigeoncache")
	log.Fatal(http.ListenAndServe(address, peers))

}

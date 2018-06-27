package main

import (
	"encoding/json"
	"fmt"

	"github.com/kpbird/pcache"
)

type student struct {
	Name string
}

func main() {
	cache := pcache.NewPCache("./cache", 50, 20)

	// Simple Key Value store
	cache.PSet("key1", []byte("value1"))
	v1, _ := cache.PGet("key1")
	fmt.Println("Value of Key 1", string(v1))

	// Store Structure in value
	s := student{Name: "ketan"}
	b, _ := json.Marshal(s)
	fmt.Println("set key2 ", b)
	cache.PSet("key2", b)

	value, _ := cache.PGet("key2")
	fmt.Println("get key2 ", value)
	var s1 student
	json.Unmarshal(value, &s1)
	fmt.Println("value key2 ", s1)

}

# PCache
Persistence Cache for Go Lang. PCache is caching library which support In-memory and Disk cache. You can configure number keys In-Memory and In-Disk.

PCache maintain In-Memory cache. When In-Memory limit reached. It move key-value to Disk. PCache move least frequently used key to disk.

# Feature

- [x] Move key to disk when memory limit reached

- [x] Retrive key from disk if key don't exist In-memory

- [x] Delete key from disk when limit reached

- [x] Goroutine safe

- [ ] Expiration support

- [ ] Iterator support

# Install

```
$ go get github.com/kpbird/pcache
```

# Example 1

Store string as value. 

```
cache := pcache.NewPCache("./cache", 50, 20)

cache.PSet("key1", []byte("value1"))
v1, _ := cache.PGet("key1")
fmt.Println("Value of Key 1", string(v1))

```

# Example 2

Store structure as value

```
type student struct {
	Name string
}

cache := pcache.NewPCache("./cache", 50, 20)
// Store Structure in value
s := student{Name: "ketan"}
b, _ := json.Marshal(s)
fmt.Println("set key2 ", b)
cache.PSet("key2", b)

// Retrive Structure from value
value, _ := cache.PGet("key2")
fmt.Println("get key2 ", value)
var s1 student
json.Unmarshal(value, &s1)
fmt.Println("value key2 ", s1)
```

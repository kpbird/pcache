package pcache

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

//PCache structure
type PCache struct {
	filePath     string
	maxKeyInMem  int
	maxKeyInDisk int
	cacheCount   map[string]int
	cache        map[string][]byte
	sync.RWMutex
}

//NewPCache - create new object of Pcache
func NewPCache(filePath string, maxKeyInMem int, maxKeyInDisk int) *PCache {
	m := &PCache{filePath: filePath,
		maxKeyInMem:  maxKeyInMem,
		maxKeyInDisk: maxKeyInDisk,
		cacheCount:   make(map[string]int),
		cache:        make(map[string][]byte)}
	return m
}

//PSet - to set value
func (pcache *PCache) PSet(key string, value []byte) {
	pcache.Lock()
	if len(pcache.cache) >= pcache.maxKeyInMem {
		pcache.saveToFile()
	}
	pcache.cache[key] = value
	pcache.cacheCount[key] = 0
	pcache.Unlock()
}

//PGet - to get value
func (pcache *PCache) PGet(key string) ([]byte, error) {
	pcache.RLock()
	// find key in memory
	value, ok := pcache.cache[key]
	if ok {
		pcache.cacheCount[key] = pcache.cacheCount[key] + 1
		pcache.RUnlock()
		return value, nil
		// find key in disk
	} else if pcache.isFileExist(key) {
		value := pcache.loadFromFile(key)
		pcache.RUnlock()
		return value, nil
	}
	pcache.RUnlock()
	// return error
	return nil, errors.New("Key not found")
}

//PRemove - to delete key
func (pcache *PCache) PRemove(key string) {
	pcache.Lock()
	delete(pcache.cache, key)
	delete(pcache.cacheCount, key)
	pcache.Unlock()
}

//// Private functions

func (pcache *PCache) saveToFile() {
	key := pcache.findLeastUsedKey()
	file := pcache.filePath + "/" + pcache.md5(key) + ".cache"
	err := ioutil.WriteFile(file, pcache.cache[key], 0644)
	if err != nil {
		fmt.Println(err)
	}
	delete(pcache.cache, key)
	delete(pcache.cacheCount, key)
	pcache.removeCachedFile()
}

func (pcache *PCache) loadFromFile(key string) []byte {

	file := pcache.filePath + "/" + pcache.md5(key) + ".cache"
	value, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}
	pcache.cache[key] = value
	pcache.cacheCount[key] = 0

	if len(pcache.cache) >= pcache.maxKeyInMem {
		pcache.saveToFile()
	}

	return value

}

func (pcache *PCache) removeCachedFile() {
	files, err := ioutil.ReadDir(pcache.filePath)
	if err != nil {
		log.Fatal(err)
	}
	if len(files) <= pcache.maxKeyInDisk {
		return
	}
	var oldestFile os.FileInfo
	oldestTime := time.Now()
	for _, file := range files {
		if file.Mode().IsRegular() && file.ModTime().Before(oldestTime) {
			oldestFile = file
			oldestTime = file.ModTime()
		}
	}

	if oldestFile != nil {
		os.Remove(pcache.filePath + "/" + oldestFile.Name())
	}

}

func (pcache *PCache) findLeastUsedKey() string {
	minV := -1
	minK := ""
	for k, v := range pcache.cacheCount {
		if minV == -1 {
			minV = v
			minK = k
		}

		if v < minV {
			minV = v
			minK = k
		}
	}
	return minK
}

func (pcache *PCache) isFileExist(key string) bool {
	path := pcache.filePath + "/" + key + ".gob"
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (pcache *PCache) md5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

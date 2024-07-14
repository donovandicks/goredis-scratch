package interpreter

import (
	"sync"

	"golang.org/x/exp/maps"
)

type HashStore struct {
	sync.RWMutex
	data map[string]map[string]string
}

func (hs *HashStore) Set(hashName, key, val string) {
	hs.Lock()
	_, ok := hs.data[hashName]
	if !ok {
		hs.data[hashName] = make(map[string]string)
	}

	hs.data[hashName][key] = val
	hs.Unlock()
}

func (hs *HashStore) Get(hashName, key string) (string, bool) {
	hs.RLock()
	val, ok := hs.data[hashName][key]
	hs.RUnlock()
	return val, ok
}

func (hs *HashStore) GetAll(hashName string) []string {

	hs.RLock()
	hash, ok := hs.data[hashName]
	if !ok {
		return nil
	}

	vals := make([]string, 0)
	keys := maps.Keys(hash)
	for _, key := range keys {
		val, _ := hash[key]
		vals = append(vals, val)
	}
	hs.RUnlock()
	return vals
}

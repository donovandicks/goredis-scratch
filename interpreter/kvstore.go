package interpreter

import "sync"

type KVStore struct {
	sync.RWMutex
	data map[string]string
}

func (kv *KVStore) Set(key, val string) {
	kv.Lock()
	kv.data[key] = val
	kv.Unlock()
}

func (kv *KVStore) Get(key string) (string, bool) {
	kv.RLock()
	val, ok := kv.data[key]
	kv.RUnlock()
	return val, ok
}

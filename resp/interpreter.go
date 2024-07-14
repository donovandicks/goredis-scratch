package resp

import (
	"fmt"
	"sync"
)

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

type Interpreter struct {
	handlers map[string]func([]Value) Value
	kvStore  KVStore
}

func NewInterpreter() *Interpreter {
	interp := &Interpreter{
		kvStore: KVStore{data: make(map[string]string)},
	}

	interp.handlers = map[string]func([]Value) Value{
		"PING": interp.ping,
		"GET":  interp.get,
		"SET":  interp.set,
	}

	return interp
}

// ping returns the first argument after the command, or PONG if none provided.
func (i *Interpreter) ping(args []Value) Value {
	if len(args) == 0 {
		return NewStringValue("PONG")
	}

	return NewStringValue(args[0].Bulk)
}

func (i *Interpreter) set(args []Value) Value {
	if len(args) != 2 {
		return NewErrorValue("ERR wrong number of arguments for 'set' command")
	}

	key := args[0].Bulk
	val := args[1].Bulk

	i.kvStore.Set(key, val)

	return NewStringValue("OK")
}

func (i *Interpreter) get(args []Value) Value {
	if len(args) != 1 {
		return NewErrorValue("ERR wrong number of arguments for 'get' command")
	}

	val, ok := i.kvStore.Get(args[0].Bulk)
	if !ok {
		return NewNullValue()
	}

	return NewBulkValue(val)
}

func (i *Interpreter) Interpret(command string, args []Value) Value {
	handler, ok := i.handlers[command]
	if !ok {
		return NewErrorValue(fmt.Sprintf("Command %v not supported", command))
	}

	return handler(args)
}

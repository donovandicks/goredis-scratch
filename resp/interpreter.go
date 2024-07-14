package resp

import (
	"fmt"
	"sync"
)

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"GET":  get,
	"SET":  set,
}

type KVStore struct {
	sync.RWMutex
	data map[string]string
}

var Sets = KVStore{
	data: make(map[string]string),
}

// ping returns the first argument after the command, or PONG if none provided.
func ping(args []Value) Value {
	if len(args) == 0 {
		return NewStringValue("PONG")
	}

	return NewStringValue(args[0].Bulk)
}

func set(args []Value) Value {
	if len(args) != 2 {
		return NewErrorValue("ERR wrong number of arguments for 'set' command")
	}

	key := args[0].Bulk
	val := args[1].Bulk

	Sets.Lock()
	Sets.data[key] = val
	Sets.Unlock()

	return NewStringValue("OK")
}

func get(args []Value) Value {
	if len(args) != 1 {
		return NewErrorValue("ERR wrong number of arguments for 'get' command")
	}

	Sets.RLock()
	val, ok := Sets.data[args[0].Bulk]
	Sets.RUnlock()
	if !ok {
		return NewNullValue()
	}

	return NewBulkValue(val)
}

type Interpreter struct{}

func (i *Interpreter) Interpret(command string, args []Value) Value {
	handler, ok := Handlers[command]
	if !ok {
		return NewErrorValue(fmt.Sprintf("Command %v not supported", command))
	}

	return handler(args)
}

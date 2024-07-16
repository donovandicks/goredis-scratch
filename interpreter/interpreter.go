package interpreter

import (
	"fmt"

	"github.com/donovandicks/goredis/resp"
)

var (
	WriteCommands = map[string]struct{}{
		"SET":  {},
		"HSET": {},
	}
)

// Check if the given command writes values
func CommandIsWrite(command string) bool {
	_, ok := WriteCommands[command]
	return ok
}

type Interpreter struct {
	handlers  map[string]func([]resp.Value) resp.Value
	kvStore   KVStore
	hashStore HashStore
}

func NewInterpreter() *Interpreter {
	interp := &Interpreter{
		kvStore:   KVStore{data: make(map[string]string)},
		hashStore: HashStore{data: make(map[string]map[string]string)},
	}

	interp.registerHandlers()

	return interp
}

func (i *Interpreter) registerHandlers() {
	i.handlers = map[string]func([]resp.Value) resp.Value{
		"PING":    i.ping,
		"GET":     i.get,
		"SET":     i.set,
		"HSET":    i.hset,
		"HGET":    i.hget,
		"HGETALL": i.hgetall,
	}
}

// ping returns the first argument after the command, or PONG if none provided.
func (i *Interpreter) ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.NewStringValue("PONG")
	}

	return resp.NewStringValue(args[0].Bulk)
}

func (i *Interpreter) set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.NewErrorValue("ERR wrong number of arguments for 'set' command")
	}

	key := args[0].Bulk
	val := args[1].Bulk

	i.kvStore.Set(key, val)

	return resp.NewOKValue()
}

func (i *Interpreter) get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewErrorValue("ERR wrong number of arguments for 'get' command")
	}

	val, ok := i.kvStore.Get(args[0].Bulk)
	if !ok {
		return resp.NewNullValue()
	}

	return resp.NewBulkValue(val)
}

func (i *Interpreter) hset(args []resp.Value) resp.Value {
	if len(args) != 3 {
		return resp.NewErrorValue("ERR wrong number of arguments for 'hset' command")
	}

	hashName := args[0].Bulk
	key := args[1].Bulk
	val := args[2].Bulk

	i.hashStore.Set(hashName, key, val)
	return resp.NewOKValue()
}

func (i *Interpreter) hget(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.NewStringValue("ERR wrong number of arguments for 'hget' command")
	}

	hashName := args[0].Bulk
	key := args[1].Bulk

	val, ok := i.hashStore.Get(hashName, key)
	if !ok {
		return resp.NewNullValue()
	}

	return resp.NewStringValue(val)
}

func (i *Interpreter) hgetall(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewStringValue("ERR wrong number of arguments for 'hgetall' command")
	}

	hashName := args[0].Bulk

	vals := i.hashStore.GetAll(hashName)
	if vals == nil {
		return resp.NewNullValue()
	}

	out := resp.Value{Typ: resp.Array}
	for _, val := range vals {
		out.Array = append(out.Array, resp.NewStringValue(val))
	}

	return out
}

func (i *Interpreter) Interpret(command string, args []resp.Value) resp.Value {
	handler, ok := i.handlers[command]
	if !ok {
		return resp.NewErrorValue(fmt.Sprintf("Command %v not supported", command))
	}

	return handler(args)
}

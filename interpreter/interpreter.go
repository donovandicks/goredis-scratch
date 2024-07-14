package interpreter

import (
	"fmt"

	"github.com/donovandicks/goredis/resp"
)

type Interpreter struct {
	handlers map[string]func([]resp.Value) resp.Value
	kvStore  KVStore
}

func NewInterpreter() *Interpreter {
	interp := &Interpreter{
		kvStore: KVStore{data: make(map[string]string)},
	}

	interp.registerHandlers()

	return interp
}

func (i *Interpreter) registerHandlers() {
	i.handlers = map[string]func([]resp.Value) resp.Value{
		"PING": i.ping,
		"GET":  i.get,
		"SET":  i.set,
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

	return resp.NewStringValue("OK")
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

func (i *Interpreter) Interpret(command string, args []resp.Value) resp.Value {
	handler, ok := i.handlers[command]
	if !ok {
		return resp.NewErrorValue(fmt.Sprintf("Command %v not supported", command))
	}

	return handler(args)
}

package persistence

import (
	"github.com/donovandicks/goredis/interpreter"
	"github.com/donovandicks/goredis/resp"
)

type Strategy interface {
	Read(interp *interpreter.Interpreter)
	Write(val resp.Value) error
	Close() error
}

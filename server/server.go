package server

import (
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/donovandicks/goredis/interpreter"
	"github.com/donovandicks/goredis/persistence"
	"github.com/donovandicks/goredis/resp"
)

const (
	BufferSize = 1024
)

type Server struct {
	listener net.Listener
	interp   *interpreter.Interpreter
	persist  persistence.Strategy
}

func NewServer(persist persistence.Strategy) *Server {
	interp := interpreter.NewInterpreter()

	persist.Read(interp)

	return &Server{
		interp:  interp,
		persist: persist,
	}
}

func (s *Server) Run() {
	s.start()
	s.recv()
}

func (s *Server) start() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("listening on port 6379...")
	s.listener = listener
}

func (s *Server) recv() {
	conn, err := s.listener.Accept()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	defer conn.Close()

	for {
		parser := resp.NewParser(conn)
		value, err := parser.Read()
		if err != nil {
			if err.Error() == "EOF" {
				return
			}

			slog.Error(err.Error())
			return
		}

		if value.Typ != resp.Array {
			slog.Error(fmt.Sprintf("Invalid request type %v, expected array", value.Typ))
			continue
		}

		slog.Info(fmt.Sprintf("Received array: %+v", value.Array))
		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		writer := NewWriter(conn)

		out := s.interp.Interpret(command, args)
		writer.Write(out)

		if out.Typ != resp.Error && interpreter.CommandIsWrite(command) {
			s.persist.Write(value)
		}
	}
}

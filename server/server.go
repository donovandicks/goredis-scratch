package server

import (
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/donovandicks/goredis/interpreter"
	"github.com/donovandicks/goredis/resp"
)

const (
	BufferSize = 1024
)

type Server struct {
	listener net.Listener
}

func NewServer() *Server {
	return &Server{}
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
	interpreter := interpreter.NewInterpreter()

	for {
		parser := resp.NewParer(conn)
		value, err := parser.Read()
		if err != nil {
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

		writer := resp.NewWriter(conn)

		out := interpreter.Interpret(command, args)
		writer.Write(out)
	}
}

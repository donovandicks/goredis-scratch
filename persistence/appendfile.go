package persistence

import (
	"bufio"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/donovandicks/goredis/interpreter"
	"github.com/donovandicks/goredis/resp"
)

type AppendOnlyFile struct {
	sync.Mutex
	file   *os.File
	reader *bufio.Reader
}

// Create a new AppendOnlyFile with a file object at the given path.
//
// Will use any file found at the existing path, otherwise will create
// a new one.
func NewFile(path string) (*AppendOnlyFile, error) {
	// Create or open the file
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &AppendOnlyFile{
		file:   f,
		reader: bufio.NewReader(f),
	}

	// Sync the file to disk every second
	go func() {
		for {
			aof.Lock()
			aof.file.Sync()
			aof.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

// Close the underlying file object on the AoF.
func (aof *AppendOnlyFile) Close() error {
	aof.Lock()
	defer aof.Unlock()

	return aof.file.Close()
}

// Write the marshalled value into the underyling file object.
func (aof *AppendOnlyFile) Write(val resp.Value) error {
	aof.Lock()
	defer aof.Unlock()

	_, err := aof.file.Write(val.Marshal())
	return err
}

func (aof *AppendOnlyFile) Read(interp *interpreter.Interpreter) {
	parser := resp.NewParser(aof.file)
	values, _ := parser.ReadMany()
	for _, val := range values {
		command := strings.ToUpper(val.Array[0].Bulk)
		args := val.Array[1:]
		interp.Interpret(command, args)
	}
}

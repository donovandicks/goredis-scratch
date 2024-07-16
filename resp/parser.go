package resp

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"strconv"
)

const (
	STR   = '+'
	ERR   = '-'
	INT   = ':'
	BULK  = '$'
	ARRAY = '*'
)

type Parser struct {
	reader *bufio.Reader
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{
		reader: bufio.NewReader(reader),
	}
}

// readBulk parses a bulk string message
//
// An example:
//
//	$4\r\nUser\r\n
//
// Where '$' indicates the data type (bulk string) and '4' indicates the length.
// The following line (User\r\n) is the data itself with return carriage.
func (p *Parser) readBulk() (Value, error) {
	v := Value{}
	v.Typ = Bulk

	len, _, err := p.readInt()
	if err != nil {
		return v, nil
	}

	bulk := make([]byte, len)
	p.reader.Read(bulk)

	v.Bulk = string(bulk)
	p.readLine() // read the remaining carriage return
	return v, nil
}

// readLine reads a single line from the current buffer.
//
// It returns (bytesInLine, numBytesInLine, error)
func (p *Parser) readLine() ([]byte, int, error) {
	line := []byte{}
	n := 0
	for {
		b, err := p.reader.ReadByte()
		if err != nil {
			return nil, n, err
		}

		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			// Exit at the end of input
			break
		}
	}

	return line[:len(line)-2], n, nil
}

// readInt parses the the current line as an integer.
//
// It returns (theNumber, numBytesInLine, error)
func (p *Parser) readInt() (int, int, error) {
	line, n, err := p.readLine()
	if err != nil {
		return 0, 0, err
	}

	num, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}

	return int(num), n, nil
}

// readArray reads an input array of commands/
//
// An array begins with '*' followed by a number indicating the
// length of the array. Array elements are pairs of lines separated by
// '\r\n' and of the form <type><length>\r\n<value>
func (p *Parser) readArray() (Value, error) {
	v := Value{}
	v.Typ = Array

	len, _, err := p.readInt()
	if err != nil {
		return v, err
	}

	v.Array = make([]Value, 0)
	for i := 0; i < len; i++ {
		val, err := p.Read()
		if err != nil {
			return val, err
		}

		v.Array = append(v.Array, val)
	}

	return v, nil
}

func (p *Parser) Read() (Value, error) {
	typ, err := p.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch typ {
	case ARRAY:
		return p.readArray()
	case BULK:
		return p.readBulk()
	default:
		slog.Info(fmt.Sprintf("Unsupported type: %v", typ))
		return Value{}, nil
	}
}

func (p *Parser) ReadMany() ([]Value, error) {
	vals := make([]Value, 0)
	for {
		val, err := p.Read()
		if err != nil {
			if err.Error() == "EOF" {
				return vals, nil
			}

			return nil, err
		}

		vals = append(vals, val)
	}
}

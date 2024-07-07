// File: protocol/orsp.go

package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	SimpleString = '+'
	Error        = '-'
	Integer      = ':'
	BulkString   = '$'
	Array        = '*'
)

type ORSPValue interface {
	Marshal() string
}

type SimpleStringValue string
type ErrorValue string
type IntegerValue int64
type BulkStringValue string
type ArrayValue []ORSPValue

func (v SimpleStringValue) Marshal() string {
	return fmt.Sprintf("%c%s\r\n", SimpleString, string(v))
}

func (v ErrorValue) Marshal() string {
	return fmt.Sprintf("%c%s\r\n", Error, string(v))
}

func (v IntegerValue) Marshal() string {
	return fmt.Sprintf("%c%d\r\n", Integer, int64(v))
}

func (v BulkStringValue) Marshal() string {
	if v == "" {
		return "$-1\r\n"
	}
	return fmt.Sprintf("%c%d\r\n%s\r\n", BulkString, len(v), string(v))
}

func (v ArrayValue) Marshal() string {
	if v == nil {
		return "*-1\r\n"
	}
	s := fmt.Sprintf("%c%d\r\n", Array, len(v))
	for _, item := range v {
		s += item.Marshal()
	}
	return s
}

func Unmarshal(reader *bufio.Reader) (ORSPValue, error) {
	typeChar, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch typeChar {
	case SimpleString:
		line, err := readLine(reader)
		return SimpleStringValue(line), err
	case Error:
		line, err := readLine(reader)
		return ErrorValue(line), err
	case Integer:
		line, err := readLine(reader)
		if err != nil {
			return nil, err
		}
		n, err := strconv.ParseInt(line, 10, 64)
		return IntegerValue(n), err
	case BulkString:
		return readBulkString(reader)
	case Array:
		return readArray(reader)
	default:
		return nil, fmt.Errorf("unknown type: %c", typeChar)
	}
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	return strings.TrimRight(line, "\r\n"), err
}

func readBulkString(reader *bufio.Reader) (BulkStringValue, error) {
	lenStr, err := readLine(reader)
	if err != nil {
		return "", err
	}
	length, err := strconv.Atoi(lenStr)
	if err != nil {
		return "", err
	}
	if length == -1 {
		return "", nil
	}
	bulk := make([]byte, length+2) // +2 for \r\n
	_, err = io.ReadFull(reader, bulk)
	return BulkStringValue(bulk[:length]), err
}

func readArray(reader *bufio.Reader) (ArrayValue, error) {
	lenStr, err := readLine(reader)
	if err != nil {
		return nil, err
	}
	length, err := strconv.Atoi(lenStr)
	if err != nil {
		return nil, err
	}
	if length == -1 {
		return nil, nil
	}
	array := make(ArrayValue, length)
	for i := 0; i < length; i++ {
		array[i], err = Unmarshal(reader)
		if err != nil {
			return nil, err
		}
	}
	return array, nil
}

package protocol

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
)

const (
	SimpleString   = '+'
	Error          = '-'
	Integer        = ':'
	BulkString     = '$'
	Array          = '*'
	Null           = '_'
	Boolean        = '#'
	Double         = ','
	BigNumber      = '('
	BulkError      = '!'
	VerbatimString = '='
	Map            = '%'
	Set            = '~'
	Push           = '>'
)

type ORSPValue interface {
	Marshal() string
}

type SimpleStringValue string
type ErrorValue string
type IntegerValue int64
type BulkStringValue string
type ArrayValue []ORSPValue
type NullValue struct{}
type BooleanValue bool
type DoubleValue float64
type BigNumberValue struct{ big.Int }
type BulkErrorValue struct {
	Code    string
	Message string
}
type VerbatimStringValue struct {
	Format string
	Value  string
}
type MapValue map[string]ORSPValue
type SetValue []ORSPValue
type PushValue struct {
	Kind string
	Data []ORSPValue
}

func (v NullValue) Marshal() string {
	return "_\r\n"
}

func (v SimpleStringValue) Marshal() string {
	return fmt.Sprintf("+%s\r\n", string(v))
}

func (v BooleanValue) Marshal() string {
	if bool(v) {
		return "#t\r\n"
	}
	return "#f\r\n"
}

func (v DoubleValue) Marshal() string {
	return fmt.Sprintf(",%g\r\n", float64(v))
}

func (v *BigNumberValue) Marshal() string {
	if v == nil {
		return "(nil\r\n"
	}
	return fmt.Sprintf("(%s\r\n", v.String())
}

func (v BulkStringValue) Marshal() string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(v), v)
}

func (v IntegerValue) Marshal() string {
	return fmt.Sprintf(":%d\r\n", int64(v))
}

func (v ErrorValue) Marshal() string {
	return fmt.Sprintf("-%s\r\n", string(v))
}

func (v BulkErrorValue) Marshal() string {
	return fmt.Sprintf("!%d\r\n%s\r\n%s\r\n", len(v.Code), v.Code, v.Message)
}

func (v VerbatimStringValue) Marshal() string {
	return fmt.Sprintf("=%d\r\n%s:%s\r\n", len(v.Format)+1+len(v.Value), v.Format, v.Value)
}

func (v MapValue) Marshal() string {
	s := fmt.Sprintf("%%%d\r\n", len(v))
	for key, value := range v {
		s += BulkStringValue(key).Marshal() + value.Marshal()
	}
	return s
}

func (v SetValue) Marshal() string {
	s := fmt.Sprintf("~%d\r\n", len(v))
	for _, item := range v {
		s += item.Marshal()
	}
	return s
}

func (v PushValue) Marshal() string {
	s := fmt.Sprintf(">%d\r\n%s\r\n", len(v.Data)+1, v.Kind)
	for _, item := range v.Data {
		s += item.Marshal()
	}
	return s
}

func (v ArrayValue) Marshal() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*%d\r\n", len(v)))
	for _, item := range v {
		sb.WriteString(item.Marshal())
	}
	return sb.String()
}

func Unmarshal(reader *bufio.Reader) (ORSPValue, error) {
	typeChar, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch typeChar {
	case SimpleString:
		line, err := readLine(reader)
		if err != nil {
			return nil, err
		}
		return SimpleStringValue(line), nil
	case Error:
		line, err := readLine(reader)
		if err != nil {
			return nil, err
		}
		return ErrorValue(line), nil
	case Integer:
		line, err := readLine(reader)
		if err != nil {
			return nil, err
		}
		n, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return nil, err
		}
		return IntegerValue(n), nil
	case BulkString:
		length, err := parseLength(reader)
		if err != nil {
			return nil, err
		}
		if length < 0 {
			return nil, nil // Null bulk string
		}
		data := make([]byte, length)
		_, err = io.ReadFull(reader, data)
		if err != nil {
			return nil, err
		}
		// Read the trailing \r\n
		_, err = reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		return BulkStringValue(data), nil
	case Array:
		return readArray(reader)
	case Null:
		return NullValue{}, nil
	case Boolean:
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		_, err = reader.ReadString('\n')
		return BooleanValue(b == 't'), err
	case Double:
		line, err := readLine(reader)
		if err != nil {
			return nil, err
		}
		f, err := strconv.ParseFloat(line, 64)
		return DoubleValue(f), err
	case BigNumber:
		line, err := readLine(reader)
		if err != nil {
			return nil, err
		}
		n := new(big.Int)
		n, ok := n.SetString(line, 10)
		if !ok {
			return nil, fmt.Errorf("invalid big number: %s", line)
		}
		return &BigNumberValue{*n}, nil
	case BulkError:
		return readBulkError(reader)
	case VerbatimString:
		return readVerbatimString(reader)
	case Map:
		return readMap(reader)
	case Set:
		return readSet(reader)
	case Push:
		return readPush(reader)
	default:
		return nil, fmt.Errorf("unknown type: %c", typeChar)
	}
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
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
	data := make([]byte, length)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return "", err
	}
	_, err = reader.ReadString('\n')
	return BulkStringValue(data), err
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
	arr := make(ArrayValue, length)
	for i := 0; i < length; i++ {
		arr[i], err = Unmarshal(reader)
		if err != nil {
			return nil, err
		}
	}
	return arr, nil
}

func readBulkError(reader *bufio.Reader) (BulkErrorValue, error) {
	lenStr, err := readLine(reader)
	if err != nil {
		return BulkErrorValue{}, err
	}
	length, err := strconv.Atoi(lenStr)
	if err != nil {
		return BulkErrorValue{}, err
	}

	// Read the code and message
	code, err := readLine(reader)
	if err != nil {
		return BulkErrorValue{}, err
	}
	message, err := readLine(reader)
	if err != nil {
		return BulkErrorValue{}, err
	}

	// Check if the length matches the expected length
	expectedLength := len(code) + len(message) + 2 // +2 for the \r\n
	if length != expectedLength {
		return BulkErrorValue{}, fmt.Errorf("length mismatch: expected %d, got %d", length, expectedLength)
	}

	return BulkErrorValue{Code: code, Message: message}, nil
}

func readVerbatimString(reader *bufio.Reader) (VerbatimStringValue, error) {
	lenStr, err := readLine(reader)
	if err != nil {
		return VerbatimStringValue{}, err
	}
	length, err := strconv.Atoi(lenStr)
	if err != nil {
		return VerbatimStringValue{}, err
	}
	data := make([]byte, length)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return VerbatimStringValue{}, err
	}
	parts := strings.SplitN(string(data[:length]), ":", 2)
	if len(parts) != 2 {
		return VerbatimStringValue{}, fmt.Errorf("invalid verbatim string format")
	}
	return VerbatimStringValue{Format: parts[0], Value: parts[1]}, nil
}

func readMap(reader *bufio.Reader) (MapValue, error) {
	lenStr, err := readLine(reader)
	if err != nil {
		return nil, err
	}
	length, err := strconv.Atoi(lenStr)
	if err != nil {
		return nil, err
	}
	m := make(MapValue, length)
	for i := 0; i < length; i++ {
		key, err := Unmarshal(reader)
		if err != nil {
			return nil, err
		}
		keyStr, ok := key.(SimpleStringValue)
		if !ok {
			return nil, fmt.Errorf("map key must be a simple string")
		}
		value, err := Unmarshal(reader)
		if err != nil {
			return nil, err
		}
		m[string(keyStr)] = value
	}
	return m, nil
}

func readSet(reader *bufio.Reader) (SetValue, error) {
	lenStr, err := readLine(reader)
	if err != nil {
		return nil, err
	}
	length, err := strconv.Atoi(lenStr)
	if err != nil {
		return nil, err
	}
	set := make(SetValue, length)
	for i := 0; i < length; i++ {
		set[i], err = Unmarshal(reader)
		if err != nil {
			return nil, err
		}
	}
	return set, nil
}

func readPush(reader *bufio.Reader) (PushValue, error) {
	lenStr, err := readLine(reader)
	if err != nil {
		return PushValue{}, err
	}
	length, err := strconv.Atoi(lenStr)
	if err != nil {
		return PushValue{}, err
	}
	kind, err := readLine(reader)
	if err != nil {
		return PushValue{}, err
	}
	data := make([]ORSPValue, length-1)
	for i := 0; i < length-1; i++ {
		data[i], err = Unmarshal(reader)
		if err != nil {
			return PushValue{}, err
		}
	}
	return PushValue{Kind: kind, Data: data}, nil
}

func parseLength(reader *bufio.Reader) (int, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	line = strings.TrimSpace(line)
	return strconv.Atoi(line)
}

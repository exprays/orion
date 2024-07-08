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
		if err == io.EOF {
			return nil, err // Return EOF directly, it's handled in LoadAOF
		}
		return nil, fmt.Errorf("error reading type character: %w", err)
	}

	switch typeChar {
	case SimpleString:
		return unmarshalSimpleString(reader)
	case Error:
		return unmarshalError(reader)
	case Integer:
		return unmarshalInteger(reader)
	case BulkString:
		return unmarshalBulkString(reader)
	case Array:
		return unmarshalArray(reader)
	case Null:
		return unmarshalNull(reader)
	case Boolean:
		return unmarshalBoolean(reader)
	case Double:
		return unmarshalDouble(reader)
	case BigNumber:
		return unmarshalBigNumber(reader)
	case BulkError:
		return unmarshalBulkError(reader)
	case VerbatimString:
		return unmarshalVerbatimString(reader)
	case Map:
		return unmarshalMap(reader)
	case Set:
		return unmarshalSet(reader)
	case Push:
		return unmarshalPush(reader)
	default:
		return nil, fmt.Errorf("unknown type: %c", typeChar)
	}
}

func unmarshalSimpleString(reader *bufio.Reader) (SimpleStringValue, error) {
	line, err := readLine(reader)
	if err != nil {
		return "", fmt.Errorf("error reading simple string: %w", err)
	}
	return SimpleStringValue(line), nil
}

func unmarshalError(reader *bufio.Reader) (ErrorValue, error) {
	line, err := readLine(reader)
	if err != nil {
		return "", fmt.Errorf("error reading error value: %w", err)
	}
	return ErrorValue(line), nil
}

func unmarshalInteger(reader *bufio.Reader) (IntegerValue, error) {
	line, err := readLine(reader)
	if err != nil {
		return 0, fmt.Errorf("error reading integer: %w", err)
	}
	n, err := strconv.ParseInt(line, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer: %w", err)
	}
	return IntegerValue(n), nil
}

func unmarshalBulkString(reader *bufio.Reader) (BulkStringValue, error) {
	length, err := parseLength(reader)
	if err != nil {
		return "", fmt.Errorf("error parsing bulk string length: %w", err)
	}
	if length < 0 {
		return "", nil // Null bulk string
	}
	data := make([]byte, length)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return "", fmt.Errorf("error reading bulk string data: %w", err)
	}
	// Read the trailing \r\n
	_, err = reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading bulk string terminator: %w", err)
	}
	return BulkStringValue(data), nil
}

func unmarshalArray(reader *bufio.Reader) (ArrayValue, error) {
	length, err := parseLength(reader)
	if err != nil {
		return nil, fmt.Errorf("error parsing array length: %w", err)
	}
	arr := make(ArrayValue, length)
	for i := 0; i < length; i++ {
		arr[i], err = Unmarshal(reader)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling array element %d: %w", i, err)
		}
	}
	return arr, nil
}

func unmarshalNull(reader *bufio.Reader) (NullValue, error) {
	_, err := reader.ReadString('\n')
	if err != nil {
		return NullValue{}, fmt.Errorf("error reading null terminator: %w", err)
	}
	return NullValue{}, nil
}

func unmarshalBoolean(reader *bufio.Reader) (BooleanValue, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return false, fmt.Errorf("error reading boolean value: %w", err)
	}
	_, err = reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("error reading boolean terminator: %w", err)
	}
	return BooleanValue(b == 't'), nil
}

func unmarshalDouble(reader *bufio.Reader) (DoubleValue, error) {
	line, err := readLine(reader)
	if err != nil {
		return 0, fmt.Errorf("error reading double value: %w", err)
	}
	f, err := strconv.ParseFloat(line, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid double value: %w", err)
	}
	return DoubleValue(f), nil
}

func unmarshalBigNumber(reader *bufio.Reader) (*BigNumberValue, error) {
	line, err := readLine(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading big number: %w", err)
	}
	n := new(big.Int)
	n, ok := n.SetString(line, 10)
	if !ok {
		return nil, fmt.Errorf("invalid big number: %s", line)
	}
	return &BigNumberValue{*n}, nil
}

func unmarshalBulkError(reader *bufio.Reader) (BulkErrorValue, error) {
	length, err := parseLength(reader)
	if err != nil {
		return BulkErrorValue{}, fmt.Errorf("error parsing bulk error length: %w", err)
	}

	code, err := readLine(reader)
	if err != nil {
		return BulkErrorValue{}, fmt.Errorf("error reading bulk error code: %w", err)
	}
	message, err := readLine(reader)
	if err != nil {
		return BulkErrorValue{}, fmt.Errorf("error reading bulk error message: %w", err)
	}

	expectedLength := len(code) + len(message) + 2 // +2 for the \r\n
	if length != expectedLength {
		return BulkErrorValue{}, fmt.Errorf("bulk error length mismatch: expected %d, got %d", length, expectedLength)
	}

	return BulkErrorValue{Code: code, Message: message}, nil
}

func unmarshalVerbatimString(reader *bufio.Reader) (VerbatimStringValue, error) {
	length, err := parseLength(reader)
	if err != nil {
		return VerbatimStringValue{}, fmt.Errorf("error parsing verbatim string length: %w", err)
	}
	data := make([]byte, length)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return VerbatimStringValue{}, fmt.Errorf("error reading verbatim string data: %w", err)
	}
	parts := strings.SplitN(string(data[:length]), ":", 2)
	if len(parts) != 2 {
		return VerbatimStringValue{}, fmt.Errorf("invalid verbatim string format")
	}
	return VerbatimStringValue{Format: parts[0], Value: parts[1]}, nil
}

func unmarshalMap(reader *bufio.Reader) (MapValue, error) {
	length, err := parseLength(reader)
	if err != nil {
		return nil, fmt.Errorf("error parsing map length: %w", err)
	}
	m := make(MapValue, length)
	for i := 0; i < length; i++ {
		key, err := Unmarshal(reader)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling map key %d: %w", i, err)
		}
		keyStr, ok := key.(SimpleStringValue)
		if !ok {
			return nil, fmt.Errorf("map key must be a simple string, got %T", key)
		}
		value, err := Unmarshal(reader)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling map value %d: %w", i, err)
		}
		m[string(keyStr)] = value
	}
	return m, nil
}

func unmarshalSet(reader *bufio.Reader) (SetValue, error) {
	length, err := parseLength(reader)
	if err != nil {
		return nil, fmt.Errorf("error parsing set length: %w", err)
	}
	set := make(SetValue, length)
	for i := 0; i < length; i++ {
		set[i], err = Unmarshal(reader)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling set element %d: %w", i, err)
		}
	}
	return set, nil
}

func unmarshalPush(reader *bufio.Reader) (PushValue, error) {
	length, err := parseLength(reader)
	if err != nil {
		return PushValue{}, fmt.Errorf("error parsing push length: %w", err)
	}
	kind, err := readLine(reader)
	if err != nil {
		return PushValue{}, fmt.Errorf("error reading push kind: %w", err)
	}
	data := make([]ORSPValue, length-1)
	for i := 0; i < length-1; i++ {
		data[i], err = Unmarshal(reader)
		if err != nil {
			return PushValue{}, fmt.Errorf("error unmarshaling push data element %d: %w", i, err)
		}
	}
	return PushValue{Kind: kind, Data: data}, nil
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func parseLength(reader *bufio.Reader) (int, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	line = strings.TrimSpace(line)
	return strconv.Atoi(line)
}

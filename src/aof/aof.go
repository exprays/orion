package aof

import (
	"bufio"
	"fmt"
	"orion/src/protocol"
	"os"
	"sync"
)

var (
	aofFile           *os.File
	appendedCommands  map[string]struct{}
	appendedCommandsM sync.Mutex
)

// InitAOF initializes the AOF system
func InitAOF() error {
	var err error
	aofFile, err = os.OpenFile("appendonly.orion", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("error opening AOF file: %w", err)
	}
	appendedCommands = make(map[string]struct{})
	return nil
}

// AppendCommand appends a command to the AOF file (if not already appended)
func AppendCommand(command protocol.ArrayValue) error {
	appendedCommandsM.Lock()
	defer appendedCommandsM.Unlock()

	commandStr := command.Marshal()
	if _, ok := appendedCommands[commandStr]; ok {
		return nil // Command already exists
	}

	if aofFile == nil {
		return fmt.Errorf("AOF file not initialized")
	}

	_, err := aofFile.WriteString(commandStr + "\n")
	if err != nil {
		return fmt.Errorf("error writing to AOF file: %w", err)
	}

	appendedCommands[commandStr] = struct{}{}
	return nil
}

// LoadAOF replays commands in the AOF file to restore server state
func LoadAOF(handleCommand func(command protocol.ArrayValue) error) error {
	file, err := os.Open("appendonly.orion")
	if err != nil {
		return fmt.Errorf("error opening AOF file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	commandCount := 0

	for {
		command, err := protocol.Unmarshal(reader)
		if err != nil {
			if err.Error() == "EOF" {
				break // End of file reached
			}
			return fmt.Errorf("error unmarshaling command: %w", err)
		}

		arrayCommand, ok := command.(protocol.ArrayValue)
		if !ok {
			return fmt.Errorf("invalid command format in AOF file: expected ArrayValue")
		}

		err = handleCommand(arrayCommand)
		if err != nil {
			return fmt.Errorf("error handling command: %w", err)
		}

		commandCount++
	}

	fmt.Printf("Total commands replayed: %d\n", commandCount)
	return nil
}

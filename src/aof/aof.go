package aof

import (
	"bufio"
	"fmt"
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
func AppendCommand(command string) error {
	appendedCommandsM.Lock()
	defer appendedCommandsM.Unlock()

	if _, ok := appendedCommands[command]; ok {
		return nil // Command already exists
	}

	if aofFile == nil {
		return fmt.Errorf("AOF file not initialized")
	}

	_, err := aofFile.WriteString(command + "\n")
	if err != nil {
		return fmt.Errorf("error writing to AOF file: %w", err)
	}

	appendedCommands[command] = struct{}{}
	return nil
}

// LoadAOF replays commands in the AOF file to restore server state
func LoadAOF(handleCommand func(command string) error) error {
	file, err := os.Open("appendonly.orion")
	if err != nil {
		return fmt.Errorf("error opening AOF file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	commandCount := 0

	for scanner.Scan() {
		command := scanner.Text()
		if command == "" {
			continue // Skip empty lines
		}

		err := handleCommand(command)
		if err != nil {
			return fmt.Errorf("error handling command: %w", err)
		}

		commandCount++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading AOF file: %w", err)
	}

	fmt.Printf("Total commands replayed: %d\n", commandCount)
	return nil
}

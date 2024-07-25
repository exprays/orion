package aof

import (
	"bufio"
	"fmt"
	"io"
	"orion/src/protocol"
	"os"
	"sync"
)

var (
	aofFile           *os.File
	appendedCommands  map[string]struct{}
	appendedCommandsM sync.Mutex
)

// InitAOF initializes the AOF system should be called on server start
func InitAOF() error {
	var err error
	aofFile, err = os.OpenFile("appendonly.orion", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("error opening AOF file: %w", err)
	}
	appendedCommands = make(map[string]struct{})
	return nil
}

func AppendCommand(command protocol.ArrayValue) error {
	appendedCommandsM.Lock()
	defer appendedCommandsM.Unlock()

	if aofFile == nil {
		return fmt.Errorf("AOF file not initialized")
	}

	commandStr := command.Marshal()

	// Check if command has already been appended
	if _, exists := appendedCommands[commandStr]; exists {
		fmt.Printf("Command already appended, skipping: %s\n", commandStr)
		return nil
	}

	// _, err := aofFile.WriteString(commandStr)
	// if err != nil {
	// 	return fmt.Errorf("error writing to AOF file: %w", err)
	// }
	// error fixed: code-aof-rewrite-after-skipping

	// Mark command as appended
	appendedCommands[commandStr] = struct{}{}

	return aofFile.Sync()
}

func LoadAOF(handleCommand func(command protocol.ArrayValue) error) error {
	file, err := os.Open("appendonly.orion")
	if err != nil {
		if os.IsNotExist(err) {
			// AOF file doesn't exist, which is fine for a new instance
			return nil
		}
		return fmt.Errorf("error opening AOF file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	commandCount := 0

	for {
		// Skip any leading whitespace
		for {
			b, err := reader.ReadByte()
			if err != nil {
				if err == io.EOF {
					// End of file reached, we're done
					fmt.Printf("Total commands replayed: %d\n", commandCount)
					return nil
				}
				return fmt.Errorf("error reading AOF file: %w", err)
			}
			if !isWhitespace(b) {
				reader.UnreadByte()
				break
			}
		}

		command, err := protocol.Unmarshal(reader)
		if err != nil {
			if err == io.EOF {
				// End of file reached while trying to unmarshal, we're done
				fmt.Printf("Total commands replayed: %d\n", commandCount)
				return nil
			}
			// Print the content of the file at the point of error
			currentPosition, _ := file.Seek(0, io.SeekCurrent)
			errorContext := make([]byte, 100)
			_, readErr := file.ReadAt(errorContext, currentPosition-50)
			if readErr != nil && readErr != io.EOF {
				fmt.Printf("Error reading error context: %v\n", readErr)
			}
			fmt.Printf("Error context: %s\n", string(errorContext))

			fmt.Printf("Error unmarshaling command at position %d: %v\n", currentPosition, err)

			// Try to skip to the next command
			for {
				b, err := reader.ReadByte()
				if err != nil {
					if err == io.EOF {
						return nil // End of file reached while skipping
					}
					return fmt.Errorf("error skipping corrupted data: %w", err)
				}
				if b == '*' {
					reader.UnreadByte() // Put back the '*' character
					break
				}
			}
			continue
		}

		arrayCommand, ok := command.(protocol.ArrayValue)
		if !ok {
			fmt.Printf("Warning: Invalid command format in AOF file: expected ArrayValue, got %T\n", command)
			continue // Skip this command and continue with the next one
		}

		// fmt.Printf("Executing command from AOF: %v\n", arrayCommand)
		_ = handleCommand(arrayCommand)

		// if err := handleCommand(arrayCommand); err != nil {
		// 	// fmt.Printf("Warning: Error handling command %v: %v\n", arrayCommand, err)
		// 	fmt.Printf("error handling command")
		// 	// Continue loading instead of returning an error
		// }

		// Execute the command without printing
		if err := handleCommand(arrayCommand); err != nil {
			// Optionally, you can log errors here if needed
			fmt.Printf("Error handling command: %v", err)
		}

		commandCount++
	}
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

// CloseAOF closes the AOF file
func CloseAOF() error {
	if aofFile != nil {
		return aofFile.Close()
	}
	return nil
}

// RewriteAOF rewrites the AOF file to optimize storage
func RewriteAOF(getCurrentState func() ([]protocol.ArrayValue, error)) error {
	// Get the current state of the database
	commands, err := getCurrentState()
	if err != nil {
		return fmt.Errorf("error getting current state: %w", err)
	}

	// Create a temporary file for the new AOF
	tempFile, err := os.CreateTemp("", "appendonly.orion.temp")
	if err != nil {
		return fmt.Errorf("error creating temp file: %w", err)
	}
	defer tempFile.Close()

	// Write the current state to the temporary file
	for _, cmd := range commands {
		_, err := tempFile.WriteString(cmd.Marshal())
		if err != nil {
			return fmt.Errorf("error writing to temp file: %w", err)
		}
	}

	// Close the current AOF file
	if err := CloseAOF(); err != nil {
		return fmt.Errorf("error closing current AOF file: %w", err)
	}

	// Rename the temporary file to the AOF file
	if err := os.Rename(tempFile.Name(), "appendonly.orion"); err != nil {
		return fmt.Errorf("error renaming temp file: %w", err)
	}

	// Reinitialize the AOF
	return InitAOF()
}

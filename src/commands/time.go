// commands/time.go - TIME command handler

package commands

import (
	"orion/src/data"
)

// HandleTime retrieves the current server time
func HandleTime(args []string) string {
	if len(args) != 0 {
		return "ERROR: Usage: TIME"
	}
	response, err := data.Store.Time()
	if err != nil {
		return "ERROR: " + err.Error()
	}
	return response
}

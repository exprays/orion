// commands/dbsize.go - DBSIZE command handler

package commands

import (
	"fmt"
	"orion/src/data"
)

// HandleDBSize returns the number of keys in the data store
func HandleDBSize(args []string) string {
	if len(args) != 0 {
		return "ERROR: Usage: DBSIZE"
	}

	// Get the number of keys in the data store
	dbSize := data.Store.DBSize()
	return fmt.Sprintf("DBSIZE: %d", dbSize)
}

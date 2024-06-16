// commnds/getdel.go - GETDEL command handler

package commands

import "orion/src/data"

// HandleGetDel gets the value of a key and deletes it
func HandleGetDel(args []string) string {
	if len(args) != 1 {
		return "ERROR: Usage: GETDEL key"
	}
	key := args[0]
	value, ok := data.Store.Get(key)
	if !ok {
		return "nil"
	}
	data.Store.GetDel(key)
	return value
}

// commands/sadd.go

package commands

import (
	"orion/src/data"
	"strconv"
)

func HandleSAdd(args []string) string {
	if len(args) < 2 {
		return "ERROR: Wrong number of arguments for 'sadd' command"
	}

	key := args[0]
	members := args[1:]

	added := data.Store.SAdd(key, members...)

	return strconv.Itoa(added)
}

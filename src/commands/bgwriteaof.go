package commands

import (
	"fmt"
	"orion/src/aof"
	"orion/src/data"
	"orion/src/protocol"
)

// HandleBGRewriteAOF handles the BGREWRITEAOF command
func HandleBGRewriteAOF(args []protocol.ORSPValue) protocol.ORSPValue {
	go func() {
		err := aof.RewriteAOF(data.Store.GetAllCommands)
		if err != nil {
			fmt.Printf("Error in BGREWRITEAOF: %v\n", err)
		} else {
			fmt.Println("Background AOF rewrite completed")
		}
	}()
	return protocol.SimpleStringValue("Background AOF rewrite started")
}

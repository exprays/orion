package commands

import (
	"fmt"
	"orion/src/persistence"
	"orion/src/protocol"
	"sync"
	"time"
)

var (
	bgSaveMutex      sync.Mutex
	bgSaveInProgress bool
)

// HandleBGSave handles the BGSAVE command
func HandleBGSave(args []protocol.ORSPValue) protocol.ORSPValue {
	bgSaveMutex.Lock()
	if bgSaveInProgress {
		bgSaveMutex.Unlock()
		return protocol.ErrorValue("BGSAVE already in progress")
	}
	bgSaveInProgress = true
	bgSaveMutex.Unlock()

	go func() {
		defer func() {
			bgSaveMutex.Lock()
			bgSaveInProgress = false
			bgSaveMutex.Unlock()
		}()

		filename := fmt.Sprintf("dump_%d.orion", time.Now().Unix())
		err := persistence.SaveToFile(filename)
		if err != nil {
			fmt.Printf("Error in BGSAVE: %v\n", err)
		} else {
			fmt.Printf("Background save completed: %s\n", filename)
		}
	}()

	return protocol.SimpleStringValue("Background saving started")
}

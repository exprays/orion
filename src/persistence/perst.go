package persistence

import (
	"encoding/gob"
	"orion/src/data"
	"os"
)

func SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(data.Store.GetAllData())
}

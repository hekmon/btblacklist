package updater

import (
	"encoding/json"
	"fmt"
	"os"
)

type state struct {
	Ripe []string `json:"ripe"`
}

func loadStateFromDisk(path string, data interface{}) (err error) {
	file, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("can't open '%s' for reading: %v", path, err)
		return
	}
	if err = json.NewDecoder(file).Decode(data); err != nil {
		err = fmt.Errorf("JSON decoder failed: %v", err)
	}
	return file.Close()
}

func saveStateToDisk(path string, data interface{}, indent bool) (err error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		err = fmt.Errorf("can't open '%s' for writing: %v", path, err)
		return
	}
	encoder := json.NewEncoder(file)
	if indent {
		encoder.SetIndent("", "    ")
	}
	if err = encoder.Encode(data); err != nil {
		err = fmt.Errorf("JSON encoder failed: %v", err)
	}
	return file.Close()
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type config struct {
	Bind            string `json:"bind_address"`
	Port            uint16 `json:"bind_port"`
	UpdateFrequency time.Duration
	RipeSearch      []string `json:"ripe_search"`
}

func (c *config) UnmarshalJSON(data []byte) (err error) {
	type shadow config
	tmp := struct {
		UpdateFrequency uint `json:"update_frequency_hours"`
		*shadow
	}{
		shadow: (*shadow)(c),
	}
	if err = json.Unmarshal(data, &tmp); err != nil {
		return
	}
	c.UpdateFrequency = time.Duration(tmp.UpdateFrequency) * time.Hour
	return
}

func getConfig(filepath string) (c config, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				err = fmt.Errorf("%s | %s", err, closeErr)
			}
		}
	}()
	err = json.NewDecoder(file).Decode(&c)
	return
}

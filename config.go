package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"
)

type config struct {
	Bind            string `json:"bind_address"`
	Port            uint16 `json:"bind_port"`
	UpdateFrequency time.Duration
	RipeSearch      []string            `json:"ripe_search"`
	Blocklists      map[string]*url.URL `json:"external_blocklists"`
}

func (c *config) UnmarshalJSON(data []byte) (err error) {
	type shadow config
	tmp := struct {
		UpdateFrequency uint              `json:"update_frequency_hours"`
		Blocklists      map[string]string `json:"external_blocklists"`
		*shadow
	}{
		shadow: (*shadow)(c),
	}
	if err = json.Unmarshal(data, &tmp); err != nil {
		return
	}
	c.UpdateFrequency = time.Duration(tmp.UpdateFrequency) * time.Hour
	c.Blocklists = make(map[string]*url.URL, len(tmp.Blocklists))
	var tmpURL *url.URL
	for name, strURL := range tmp.Blocklists {
		if tmpURL, err = url.Parse(strURL); err != nil {
			return fmt.Errorf("can't parse '%s' external blocklist as URL: %v", name, err)
		}
		c.Blocklists[name] = tmpURL
	}
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
	if err = json.NewDecoder(file).Decode(&c); err != nil {
		return
	}
	if c.UpdateFrequency < time.Hour {
		err = fmt.Errorf("update frequency can't be lower than 1 hour")
		return
	}
	if c.Port == 0 {
		err = fmt.Errorf("binding port can't be 0")
		return
	}
	return
}

package updater

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

type cache struct {
	Compressed []byte              `json:"compressed"`
	Ripe       []string            `json:"ripe"`
	External   map[string][]string `json:"external"`
}

func (s *cache) MarshalJSON() (data []byte, err error) {
	// Prepare the shadow data structure for custom marshalling
	type shadow cache
	tmp := struct {
		Compressed string `json:"compressed"`
		*shadow
	}{
		shadow: (*shadow)(s),
	}
	// Encode gzip data to b64 before JSON marshalling
	tmp.Compressed = base64.StdEncoding.EncodeToString(s.Compressed)
	// Serialize !
	return json.Marshal(tmp)
}

func (s *cache) UnmarshalJSON(data []byte) (err error) {
	// Prepare the shadow data structure for custom marshalling
	type shadow cache
	tmp := struct {
		Compressed string `json:"compressed"`
		*shadow
	}{
		shadow: (*shadow)(s),
	}
	// Unmarshal to the tmp struct
	if err = json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("can't unmarshal data to the shadow struct: %v", err)
	}
	// Deserialize the tmp data
	if s.Compressed, err = base64.StdEncoding.DecodeString(tmp.Compressed); err != nil {
		return fmt.Errorf("can't deserialized compressed data as base64: %v", err)
	}
	return
}

func loadCacheFromDisk(path string, data interface{}) (err error) {
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

func saveCacheToDisk(path string, data interface{}, indent bool) (err error) {
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

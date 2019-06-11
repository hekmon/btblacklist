package updater

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"os"
)

const (
	cacheFile = "cache.gob.gz"
)

type cache struct {
	Ripe     []string            `json:"ripe"`
	External map[string][]string `json:"external"`
}

func loadCacheFromDisk(data interface{}) (err error) {
	// File
	file, err := os.Open(cacheFile)
	if err != nil {
		err = fmt.Errorf("can't open '%s' for reading: %v", cacheFile, err)
		return
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				err = fmt.Errorf("%s | %s", err, closeErr)
			}
		}
	}()
	// Gzip decompress
	decompressor, err := gzip.NewReader(file)
	if err != nil {
		err = fmt.Errorf("can't create the gzip decompressor: %v", err)
		return
	}
	defer func() {
		if closeErr := decompressor.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				err = fmt.Errorf("%s | %s", err, closeErr)
			}
		}
	}()
	// Deserialize
	decoder := gob.NewDecoder(decompressor)
	return decoder.Decode(data)
}

func saveCacheToDisk(data interface{}, indent bool) (err error) {
	// File
	file, err := os.OpenFile(cacheFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		err = fmt.Errorf("can't open '%s' for writing: %v", cacheFile, err)
		return
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				err = fmt.Errorf("%s | %s", err, closeErr)
			}
		}
	}()
	// Compressor
	compressor, err := gzip.NewWriterLevel(file, gzip.BestCompression)
	if err != nil {
		err = fmt.Errorf("can't create the gzip decompressor: %v", err)
		return
	}
	defer func() {
		if closeErr := compressor.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				err = fmt.Errorf("%s | %s", err, closeErr)
			}
		}
	}()
	// Serialize
	encoder := gob.NewEncoder(compressor)
	return encoder.Encode(data)
}

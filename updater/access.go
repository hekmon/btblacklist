package updater

import "bytes"

// GetGzippedDataReader returns a reader yielding gzipped blacklist
func (c *Controller) GetGzippedDataReader() (gzipReader *bytes.Reader) {
	defer c.compressedDataAccess.RUnlock()
	c.compressedDataAccess.RLock()
	if c.compressedData != nil {
		gzipReader = bytes.NewReader(c.compressedData)
	}
	return
}

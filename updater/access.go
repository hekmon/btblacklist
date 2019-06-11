package updater

import "bytes"

// GetGzippedDataReader returns a reader yielding gzipped blacklist
func (c *Controller) GetGzippedDataReader() (gzipReader *bytes.Reader, length int) {
	defer c.compressedDataAccess.RUnlock()
	c.compressedDataAccess.RLock()
	if c.compressedData != nil {
		gzipReader = bytes.NewReader(c.compressedData)
		length = len(c.compressedData)
	}
	return
}

package updater

import "bytes"

// GetGzippedDataReader returns a reader yielding gzipped blacklist
func (c *Controller) GetGzippedDataReader() (gzipReader *bytes.Reader) {
	c.compressedDataAccess.RLock()
	gzipReader = bytes.NewReader(c.compressedData)
	c.compressedDataAccess.RUnlock()
	return
}

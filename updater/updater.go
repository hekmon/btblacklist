package updater

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"time"

	"github.com/hekmon/cunits"
)

func (c *Controller) updater() {
	ticker := time.NewTicker(c.frequency)
	// Fist batch
	c.updaterBatch()
	// Next ones
	for {
		select {
		case <-ticker.C:
			c.updaterBatch()
		case <-c.ctx.Done():
			c.logger.Debug("[Updater] worker received stop signal")
			ticker.Stop()
			return
		}
	}
}

func (c *Controller) updaterBatch() {
	c.logger.Debug("[Updater] worker: starting a new batch")
	// Probing
	ripeUpdate := c.updateRipe()
	// Global update
	if len(c.compressedData) != 0 && !ripeUpdate {
		c.logger.Debug("[Updater] No new data, keeping cache")
		return
	}
	// Prepare the compressor
	compressed := bytes.NewBuffer(nil)
	compressor, err := gzip.NewWriterLevel(compressed, gzip.BestCompression)
	if err != nil {
		c.logger.Errorf("[Updater] Can't create the gzip compressor: %v", err)
		return
	}
	// Add the ripe data
	ripeReader := bytes.NewBufferString(strings.Join(c.ripeState, "\n"))
	if _, err = io.Copy(compressor, ripeReader); err != nil {
		c.logger.Errorf("[Updater] Can't copy ripe results to the compressor: %v", err)
		return
	}
	// Finalize
	if _, err = compressor.Write([]byte("\n")); err != nil {
		c.logger.Errorf("[Updater] Can't add \\n before EOF: %v", err)
		return
	}
	if err = compressor.Close(); err != nil {
		c.logger.Errorf("[Updater] Can't flush remaining bytes from the gzip compressor: %v", err)
		return
	}
	// Update the current data
	c.compressedDataAccess.Lock()
	c.compressedData = compressed.Bytes()
	c.compressedDataAccess.Unlock()
	c.logger.Infof("[Updater] %d range(s) from RIPE search compressed into %s",
		len(c.ripeState), cunits.ImportInByte(float64(len(c.compressedData))))
}

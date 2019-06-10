package updater

import (
	"bytes"
	"compress/gzip"
	"io"
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
			c.logger.Debug("[Update] worker received stop signal")
			ticker.Stop()
			return
		}
	}
}

func (c *Controller) updaterBatch() {
	c.logger.Debug("[Update] worker: starting a new batch")
	// Probing
	ripeUpdate := c.updateRipe()
	// Global update
	if !ripeUpdate {
		c.logger.Debug("[Update] No new data, keeping cache")
		return
	}
	// Prepare the compressor
	compressed := bytes.NewBuffer(nil)
	compressor, err := gzip.NewWriterLevel(compressed, gzip.BestCompression)
	if err != nil {
		c.logger.Errorf("[Update] Can't create the gzip compressor: %v", err)
		return
	}
	// Add the ripe data
	ripeReader := bytes.NewBufferString(c.ripeState)
	written, err := io.Copy(compressor, ripeReader)
	if err != nil {
		c.logger.Errorf("[Update] Can't copy ripe results to the compressor: %v", err)
		return
	}
	c.logger.Debugf("[Update] Copied %s ripe data to the compressor", cunits.ImportInByte(float64(written)))
	// Update the current data
	c.compressedDataAccess.Lock()
	c.compressedData = compressed.Bytes()
	c.compressedDataAccess.Unlock()
}

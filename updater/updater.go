package updater

import "time"

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
	// Proobing
	ripeUpdate := c.updateRipe()
	// Global update
	if ripeUpdate {
		// doStuff
	}
}

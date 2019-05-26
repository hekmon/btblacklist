package updater

import (
	"fmt"
	"strings"

	"github.com/hekmon/btblacklist/ripe"
)

func (c *Controller) updateRipe() (updateGlobal bool) {
	// Perform each searches
	results := make([][]ripe.Range, len(c.ripeSearch))
	for index, search := range c.ripeSearch {
		ranges, err := c.ripec.Search(search)
		if err != nil {
			c.logger.Errorf("[Updater] ripe search: %s: %v", search, err)
			return
		}
		results[index] = ranges
		c.logger.Debugf("[Updater] ripe search: %s: %d result(s)", search, len(ranges))
	}
	// Merge results
	uniq := ripe.RemoveDuplicates(results)
	// Create the final string
	buff := make([]string, len(uniq))
	for index, ripeRange := range uniq {
		// Debug can be nice
		if c.logger.IsDebugShown() {
			if ripeRange.Route == "" {
				c.logger.Debugf("[Updater] ripe search: got inetnum: %s: %s", ripeRange.Name, ripeRange.Range)
			} else {
				c.logger.Debugf("[Updater] ripe search: got inetnum: %s: %s (from %s)", ripeRange.Name, ripeRange.Range, ripeRange.Route)
			}
		}
		// Write the line
		buff[index] = fmt.Sprintf("%s:%s", ripeRange.Name, ripeRange.Range)
	}
	finalString := strings.Join(buff, "\n")
	// Do we need to update global state ?
	if finalString != c.ripeState {
		c.logger.Infof("[Updater] ripe results changed (%d uniq results): global state will be updated", len(uniq))
		c.ripeState = finalString
		updateGlobal = true
	}
	return
}

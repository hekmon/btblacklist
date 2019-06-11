package updater

import (
	"fmt"
	"reflect"
	"sort"

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
	// Merge & sort results
	uniq := ripe.RemoveDuplicates(results)
	list := make([]string, len(uniq))
	for index, ripeRange := range uniq {
		// Debug can be nice
		if c.logger.IsDebugShown() {
			if ripeRange.Route == "" {
				c.logger.Debugf("[Updater] ripe search: got inetnum: %s: %s", ripeRange.Name, ripeRange.Range)
			} else {
				c.logger.Debugf("[Updater] ripe search: got route: %s: %s (from %s)", ripeRange.Name, ripeRange.Range, ripeRange.Route)
			}
		}
		// Build P2B lines
		list[index] = fmt.Sprintf("%s:%s", ripeRange.Name, ripeRange.Range)
	}
	sort.Strings(list)
	// Do we need to update global state ?
	if !reflect.DeepEqual(list, c.ripeState) {
		c.logger.Infof("[Updater] ripe results changed (%d uniq results): global state will be updated", len(uniq))
		c.ripeState = list
		updateGlobal = true
	} else {
		c.logger.Debug("[Updater] ripe results identical: keeping cache")
	}
	return
}

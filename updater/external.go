package updater

import (
	"compress/gzip"
	"io/ioutil"
	"reflect"
	"strings"
)

func (c *Controller) getExternalBlockList(name, url string) (modified bool) {
	// Get blocklist
	resp, err := c.http.Get(url)
	if err != nil {
		c.logger.Errorf("[Updater] external blocklist '%s': %v", name, err)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Errorf("[Updater] external blocklist '%s': can't close the HTTP request body: %v", name, err)
		}
	}()
	// Ungzip it
	decompressor, err := gzip.NewReader(resp.Body)
	if err != nil {
		c.logger.Errorf("[Updater] external blocklist '%s': can't create the gzip decompressor: %v", name, err)
		return
	}
	defer func() {
		if err := decompressor.Close(); err != nil {
			c.logger.Errorf("[Updater] external blocklist '%s': can't close the gzip decompressor: %v", name, err)
		}
	}()
	c.logger.Debugf("[Updater] external blocklist '%s': gzip header: %+v", name, decompressor.Header)
	data, err := ioutil.ReadAll(decompressor)
	if err != nil {
		c.logger.Errorf("[Updater] external blocklist '%s': can't read data as gzip compressed: %v", name, err)
		return
	}
	// Compare to stored data
	lines := strings.Split(string(data), "\n")
	c.logger.Debugf("[Updater] external blocklist '%s': got %d line(s)", name, len(lines))
	if !reflect.DeepEqual(lines, c.externalStates[name]) {
		c.logger.Infof("[Updater] external blocklist '%s': data has changed, will update", name)
		c.externalStates[name] = lines
		modified = true
	} else {
		c.logger.Debugf("[Updater] external blocklist '%s': no changes", name)
	}
	return
}

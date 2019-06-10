package updater

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hekmon/btblacklist/ripe"

	"github.com/hekmon/hllogger"
)

const (
	timeout   = time.Minute
	stateFile = "state.json"
)

// Config allows to customize a Controller creation with New()
type Config struct {
	UpdateFrequency time.Duration
	RipeSearch      []string
	Logger          *hllogger.HlLogger
}

// New returns an initialized and ready to use Controller.
// Cancel ctx when you want to stop its workers.
// Use WaitForFullStop() to be sure they are all stopped.
func New(ctx context.Context, conf Config) (c *Controller, err error) {
	// Checks
	minFreq := timeout * time.Duration(len(conf.RipeSearch))
	if conf.UpdateFrequency < minFreq {
		err = fmt.Errorf("update frequency can not be lower than %v", minFreq)
		return
	}
	if conf.Logger == nil {
		err = errors.New("logger can't be nil")
		return
	}
	// Init the controller
	c = &Controller{
		// Config
		ripeSearch: conf.RipeSearch,
		frequency:  conf.UpdateFrequency,
		// Sub controllers
		ripec:  ripe.New(timeout),
		logger: conf.Logger,
		// State
		ctx:     ctx,
		stopped: make(chan struct{}),
	}
	// Load state
	var s state
	if err = loadStateFromDisk(stateFile, &s); err != nil {
		if !strings.HasSuffix(err.Error(), "no such file or directory") {
			err = fmt.Errorf("can't load previous searchs data from disk: %v", err)
			return
		}
		c.logger.Warningf("[Updater] can't load previous state from disk: %v", err)
		err = nil
	} else {
		c.logger.Infof("[Updater] previous state loaded from '%s'", stateFile)
	}
	// Start the workers
	c.workers.Add(1)
	go func() {
		c.updater()
		c.workers.Done()
	}()
	// Start the stop watcher
	go c.stopWatcher()
	// All good
	return
}

// Controller holds all the state & logic. Instanciate with New().
type Controller struct {
	// Config
	ripeSearch []string
	frequency  time.Duration
	// Global state
	compressedData       []byte
	compressedDataAccess sync.RWMutex
	// Sub states
	ripeState string
	// Sub controllers
	ripec  *ripe.Client
	logger *hllogger.HlLogger
	// State
	ctx     context.Context
	workers sync.WaitGroup
	stopped chan struct{}
}

func (c *Controller) stopWatcher() {
	<-c.ctx.Done()
	c.logger.Debug("[Updater] Stop signal received")
	// Wait for workers to end
	c.workers.Wait()
	// Save some state
	if err := saveStateToDisk(stateFile, state{
		Compressed: c.compressedData,
		Ripe:       c.ripeState,
	}, c.logger.IsDebugShown()); err != nil {
		c.logger.Errorf("[Updater] can't save state to disk: %v", err)
	} else {
		c.logger.Infof("[Updater] State dumped to %s", stateFile)
	}
	// We have fully stopped release WaitForFullStop()
	c.logger.Debug("[Updater] Fully stopped")
	close(c.stopped)
}

// WaitForFullStop will block until all workers have properly stopped.
// To initiate stop, cancel the context used with New().
func (c *Controller) WaitForFullStop() {
	<-c.stopped
}
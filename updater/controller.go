package updater

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hekmon/btblacklist/ripe"
	"github.com/hekmon/hllogger"
)

const (
	timeout = time.Minute
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
	// TODO
	// We have fully stopped release WaitForFullStop()
	c.logger.Debug("[Updater] Fully stopped")
	close(c.stopped)
}

// WaitForFullStop will block until all workers have properly stopped.
// To initiate stop, cancel the context used with New().
func (c *Controller) WaitForFullStop() {
	<-c.stopped
}

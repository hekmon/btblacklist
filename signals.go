package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	systemd "github.com/iguanesolutions/go-systemd"
)

const (
	httpStopTimeout = 30 * time.Second
)

func handleSignals() {
	var err error
	// If we exit, allow main goroutine to do so
	defer mainStop.Unlock()
	// Register signals
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	// Waiting for signals to catch
	for {
		sig := <-signalChannel
		switch sig {
		case syscall.SIGTERM:
			fallthrough
		case syscall.SIGINT:
			logger.Infof("[Main] Signal '%v' caught: cleaning up before exiting", sig)
			// Notify stopping
			if err = systemd.NotifyStopping(); err != nil {
				logger.Errorf("[Main] can't send systemd stopping notification: %v", err)
			} else {
				logger.Debug("[Main] systemd stopping notification sent")
			}
			// First stop the HTTP server
			logger.Debugf("[Main] Stopping the HTTP server... (%v timeout)", httpStopTimeout)
			httpCtx, httpCtxCancel := context.WithTimeout(mainCtx, httpStopTimeout)
			defer httpCtxCancel()
			if err = httpServer.Shutdown(httpCtx); err != nil {
				logger.Errorf("[Main] HTTP server shutdown: %v", err)
			} else {
				logger.Debug("[Main] HTTP server stopped")
			}
			// Cancel main ctx & wait for watcher
			mainCtxCancel()
			updaterController.WaitForFullStop()
			logger.Debugf("[Main] Signal '%v' caught: all controller(s) stopped: unlocking main goroutine to exit", sig)
			return
		default:
			logger.Warningf("[Main] Signal '%v' caught but no process set to handle it: skipping", sig)
		}
	}
}

package main

import (
	"context"
	"flag"
	"os"
	"sync"

	"github.com/hekmon/btblacklist/updater"
	"github.com/hekmon/hllogger"
	systemd "github.com/iguanesolutions/go-systemd"
)

var (
	logger            *hllogger.HlLogger
	mainCtx           context.Context
	mainCtxCancel     func()
	updaterController *updater.Controller
	mainStop          sync.Mutex
)

func main() {
	// Parse flags
	logLevelFlag := flag.Int("loglevel", 1, "Set loglevel: Debug(0) Info(1) Warning(2) Error(3) Fatal(4). Default Info.")
	confFile := flag.String("conf", "config.json", "Relative or absolute path to the json configuration file")
	flag.Parse()

	// Init logger
	var ll hllogger.LogLevel
	switch *logLevelFlag {
	case 0:
		ll = hllogger.Debug
	case 1:
		ll = hllogger.Info
	case 2:
		ll = hllogger.Warning
	case 3:
		ll = hllogger.Error
	case 4:
		ll = hllogger.Fatal
	default:
		ll = hllogger.Info
	}
	logger = hllogger.New(os.Stderr, &hllogger.Config{
		LogLevel:              ll,
		SystemdJournaldCompat: systemd.IsNotifyEnabled(),
	})
	logger.Output(" ")
	logger.Output(" • BT Blacklist •")
	logger.Output("      (￣ヘ￣)")
	logger.Output(" ")

	// Conf
	logger.Info("[Main] Loading configuration")
	conf, err := getConfig(*confFile)
	if err != nil {
		logger.Fatalf(1, "can't load config: %v", err)
	}

	// Init main context
	mainCtx, mainCtxCancel = context.WithCancel(context.Background())
	defer mainCtxCancel() // make the linter happy

	// Start the updater
	updaterController, err = updater.New(mainCtx, updater.Config{
		UpdateFrequency: conf.UpdateFrequency,
		RipeSearch:      conf.RipeSearch,
		Logger:          logger,
	})
	if err != nil {
		logger.Fatalf(2, "[Main] can't init the updater: %v", err)
	}

	// Handles signals
	mainStop.Lock()
	go handleSignals()

	// Finally ready
	if err = systemd.NotifyReady(); err != nil {
		logger.Errorf("[Main] Can't send READY signal to systemd: %v", err)
	}

	// Lock the main goroutine and wait for signal handling to unlock it
	mainStop.Lock()
}

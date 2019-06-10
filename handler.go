package main

import (
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/hekmon/cunits"
)

var (
	reqIDref uint64
)

func handler(w *loggingResponseWriter, r *http.Request) {
	// Init
	start := time.Now()
	reqID := atomic.AddUint64(&reqIDref, 1)
	// Prepare logging
	var (
		err  error
		size cunits.Bits
	)
	defer func() {
		if err != nil {
			logger.Errorf("[ReadHandler] (%d) '%s %s' from '%s': answered '%d %s' in %v (%s) but an error occured: %v",
				reqID, r.Method, r.URL, r.RemoteAddr, w.statusCode, http.StatusText(w.statusCode), time.Since(start), size, err)
		} else {
			logger.Infof("[ReadHandler] (%d) '%s %s' from '%s': answered '%d %s' in %v (%s)",
				reqID, r.Method, r.URL, r.RemoteAddr, w.statusCode, http.StatusText(w.statusCode), time.Since(start), size)
		}
	}()
	// Stream data
	w.Header().Add("Content-Type", "application/x-gzip")
	written, err := io.Copy(w, updaterController.GetGzippedDataReader())
	size = cunits.ImportInByte(float64(written))
}

package main

import (
	"io"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/hekmon/cunits/v2"
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
	// Do we have any data to stream ?
	reader, length := updaterController.GetGzippedDataReader()
	if length == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// We do !
	w.Header().Set("Content-Disposition", "attachment; filename=btblocklist.txt.gz;")
	w.Header().Set("Content-Length", strconv.Itoa(length))
	w.Header().Set("Content-Type", "application/x-gzip")
	written, err := io.Copy(w, reader)
	size = cunits.ImportInByte(float64(written))
}

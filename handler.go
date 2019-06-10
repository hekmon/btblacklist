package main

import (
	"io"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/x-gzip")
	io.Copy(w, updaterController.GetGzippedDataReader())
}

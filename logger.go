package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var logWriter io.Writer

func logger(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	crw := newCustomResponseWriter(rw)
	next(crw, r)
	end := time.Now()
	addr := r.RemoteAddr
	headersToCheck := []string{"X-Real-Ip", "X-Forwarded-For"}
	for _, headerKey := range headersToCheck {
		if val := r.Header.Get(headerKey); len(val) > 0 {
			addr = val
			break
		}
	}
	fmt.Fprintf(logWriter, "%v | %3d | %13v | %15s | %d | %s %s\n",
		end.Format("2006/01/02 - 15:04:05"),
		crw.status, end.Sub(start), addr, crw.size, r.Method, r.RequestURI)
}

type customResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (c *customResponseWriter) WriteHeader(status int) {
	c.status = status
	c.ResponseWriter.WriteHeader(status)
}

func (c *customResponseWriter) Write(b []byte) (int, error) {
	size, err := c.ResponseWriter.Write(b)
	c.size += size
	return size, err
}

func newCustomResponseWriter(w http.ResponseWriter) *customResponseWriter {
	// When WriteHeader is not called, it's safe to assume the status will be 200.
	return &customResponseWriter{
		ResponseWriter: w,
		status:         200,
	}
}

package middleware

/*
 * Inspired by Golang's TimeoutHandler: https://golang.org/src/net/http/server.go?s=101514:101582#L3212
 * and gin-timeout: https://github.com/vearne/gin-timeout
 */

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maxeth/go-account-api/model"
)

// Timeout wraps the request context with a timeout, and returns the passed http error errTimeout in case of a timeout
func Timeout(timeout time.Duration, errTimeout *model.Error) gin.HandlerFunc {
	return func(c *gin.Context) {
		// set Gin's writer as our custom writer
		tw := &timeoutWriter{ResponseWriter: c.Writer, h: make(http.Header)}
		c.Writer = tw

		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// update gin request context with out context with timeout
		c.Request = c.Request.WithContext(ctx)

		finished := make(chan struct{})        // to indicate handler finished
		panicChan := make(chan interface{}, 1) // used to handle panics if we can't recover

		go func() {
			// in case we panic, just defer and try to recover
			defer func() {
				if p := recover(); p != nil {

					panicChan <- p
				}
			}()

			c.Next() // calls subsequent middleware(s) and handler

			finished <- struct{}{} // indicate that c.Next() finished successfully
		}()

		// A considerable problem with this select statement/middleware is that one HAS to be matched if we don't want to perma-block, so we can't directly write to
		// the actual gin.responseWriter inside our own Write method because then the thread would wait for either panicChan or a timeout forever since we cant
		// send a message to the "success" channel from inside our application handlers
		// this means that we have to write any response body/header twice. once to our intermediary timeoutWriter, and once to the actual gin.ResponseWriter that is included in our "tw" struct
		select {
		case <-panicChan:
			// if we cannot recover from panic,
			// send internal server error
			e := model.NewInternal()
			tw.ResponseWriter.WriteHeader(e.Status())

			eResp, _ := json.Marshal(gin.H{
				"error": e,
			})
			tw.ResponseWriter.Write(eResp)

		case <-finished:
			// if c.next() finished successfully, set headers and write resp
			tw.mu.Lock()
			defer tw.mu.Unlock()

			// map Headers w.Header() from our intermediary timeout writer (written to by gin inside some handler)
			// to the "actual" gin responsewriter (tw.ResponseWriter) for a http response
			dst := tw.ResponseWriter.Header()
			for k, vv := range tw.Header() {
				dst[k] = vv
			}

			tw.ResponseWriter.WriteHeader(tw.code)
			// tw.wbuf will have been written to already when gin writes to tw.Write()
			tw.ResponseWriter.Write(tw.wbuf.Bytes())

		case <-ctx.Done(): // context's Done channel is closed when the deadline expires, when the returned cancel function is called, or when the parent context's Done channel is closed
			// timeout has occurred, send errTimeout and write headers
			tw.mu.Lock()
			defer tw.mu.Unlock()
			// ResponseWriter from gin
			tw.ResponseWriter.Header().Set("Content-Type", "application/json")
			tw.ResponseWriter.WriteHeader(errTimeout.Status())
			eResp, _ := json.Marshal(gin.H{
				"error": errTimeout,
			})
			tw.ResponseWriter.Write(eResp)
			c.Abort()
			tw.SetTimedOut()
		}
	}
}

// implements http.ResponseWriter, and additionally keeps state of whether the request timed out, and in that case blocks any writes.
// we return a universal error message in case of a timeout in the middleware, so this way we prevent the response writer inside a handler function to overwrite our error.

// also locks access to this writer to prevent race conditions
// holds the gin.ResponseWriter which we'll manually call Write()
// on in the middleware function to send response
type timeoutWriter struct {
	gin.ResponseWriter // we dont actually implement the writing logic, but rather just include a gin.responseWriter struct inside our timeoutWriter struct that we call the write methods on when needed

	h    http.Header
	wbuf bytes.Buffer // The zero value for Buffer is an empty buffer ready to use.

	mu          sync.Mutex
	timedOut    bool
	wroteHeader bool
	code        int
}

// Writes the response, but first makes sure there
// hasn't already been a timeout
// In http.ResponseWriter interface
func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return 0, nil
	}

	return tw.wbuf.Write(b)
}

// In http.ResponseWriter interface
func (tw *timeoutWriter) WriteHeader(code int) {
	checkWriteHeaderCode(code)
	tw.mu.Lock()
	defer tw.mu.Unlock()
	// We do not write the header if we've timed out or written the header
	if tw.timedOut || tw.wroteHeader {
		return
	}
	tw.writeHeader(code)
}

// set that the header has been written
func (tw *timeoutWriter) writeHeader(code int) {
	tw.wroteHeader = true
	tw.code = code
}

// Header "relays" the header, h, set in struct
// In http.ResponseWriter interface
func (tw *timeoutWriter) Header() http.Header {
	return tw.h
}

// SetTimeOut sets timedOut field to true
func (tw *timeoutWriter) SetTimedOut() {
	tw.timedOut = true
}

func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}

package main

// KEEP THIS: one log line per request (status + duration) is what the workspace Monitoring tab charts.

import (
	"encoding/json"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// requestLog replaces gin.Logger(): it logs every request once, in the
// contract shape, after the rest of the chain (Recovery included) has run.
func requestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logRequest(c.Request.Method, c.Request.URL.Path, c.Request.UserAgent(), c.Writer.Status(), time.Since(start))
	}
}

// requestLine is the fleet's request-log contract: one JSON object on stdout per
// finished HTTP request (health checks included), with the status and the
// duration in milliseconds. Never add bodies, cookies, auth headers or the query
// string.
type requestLine struct {
	Status     int     `json:"status"`
	DurationMs float64 `json:"durationMs"`
	Method     string  `json:"method"`
	Path       string  `json:"path"`
	UA         string  `json:"ua,omitempty"`
}

// logRequest writes one request line to stdout. json.Encoder emits the object
// and its newline in a single Write, so concurrent requests do not interleave.
func logRequest(method, path, ua string, status int, took time.Duration) {
	if status == 0 {
		status = http.StatusOK // nothing written explicitly: net/http sends 200
	}
	_ = json.NewEncoder(os.Stdout).Encode(requestLine{
		Status:     status,
		DurationMs: math.Round(float64(took)/float64(time.Microsecond)) / 1000,
		Method:     method,
		Path:       path,
		UA:         ua,
	})
}

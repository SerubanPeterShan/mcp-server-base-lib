package mcp

import (
	"net/http"
	"testing"
	"time"
)

const serverTimeoutMsg = "Server did not become ready in time"

func waitForServerReady(port string, t *testing.T) {
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatalf(serverTimeoutMsg)
		case <-ticker.C:
			resp, err := http.Get("http://localhost:" + port + "/health")
			if err == nil && resp.StatusCode == http.StatusOK {
				return
			}
		}
	}
}

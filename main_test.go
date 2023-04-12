package main

import (
	"net/http"
	"testing"
	"time"
)

const (
	addr        = "localhost:8080"
	testTimeout = 2 * time.Second
)

func TestServerIsRunning(t *testing.T) {
	// Start the server in a separate goroutine
	go func() {
		startServer(addr)
	}()

	time.Sleep(500 * time.Millisecond) // Give the server some time to start

	client := http.Client{
		Timeout: testTimeout,
	}

	// Send a request to the server
	resp, err := client.Get("http://" + addr)

	// If there's an error, fail the test
	if err != nil {
		t.Fatalf("Server is not running: %v", err)
	}

	// If the status code is not 200, fail the test
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Server is not running, expected status 200, got %d", resp.StatusCode)
	}
}


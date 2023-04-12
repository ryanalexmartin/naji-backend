package server

import (
	"net/http"
	"testing"
	"time"
)

func TestServerIsRunning(t *testing.T) {
	addr := "localhost:8082"
	go StartServer(addr)
	time.Sleep(500 * time.Millisecond)

	resp, err := http.Get("http://" + addr)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

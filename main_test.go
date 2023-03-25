package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

// func setupMockClientsMap(numClients int) func() {
// 	clients.Lock()
// 	defer clients.Unlock()
// 	for i := 0; i < numClients; i++ {
// 		conn := &websocket.Conn{}
// 		clients.m[conn] = true
// 	}
// 	return func() {
// 		clients.Lock()
// 		defer clients.Unlock()
// 		for conn := range clients.m {
// 			delete(clients.m, conn)
// 		}
// 	}
// }

func TestLoadTopics(t *testing.T) {
	filename := "topics.csv"
	topics, err := loadTopics(filename)

	if err != nil {
		t.Errorf("loadTopics failed with error: %v", err)
	}

	if len(topics) == 0 {
		t.Error("loadTopics returned an empty slice")
	}

	for _, topic := range topics {
		if len(topic) == 0 {
			t.Error("loadTopics returned a slice containing empty strings")
		}
	}
}

func TestHandleWebsocket(t *testing.T) {
	// Set up a mock WebSocket server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleWebsocket(w, r, []string{"topic1", "topic2"})
	}))
	defer server.Close()

	// Set up a mock WebSocket client
	u := url.URL{Scheme: "ws", Host: strings.TrimPrefix(server.URL, "http://")}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	// Check if the client was added to the waitingClients queue
	waitingClients.Lock()
	defer waitingClients.Unlock()
	if len(waitingClients.q) != 1 {
		t.Errorf("Client not added to waitingClients queue")
	}

	// Check if the client was removed from the waitingClients queue
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message from WebSocket connection: %v", err)
	}
	if len(waitingClients.q) != 0 {
		t.Errorf("Client not removed from waitingClients queue")
	}

	// Check if the client was added to the clients map
	clients.Lock()
	defer clients.Unlock()
	if _, ok := clients.m[conn]; !ok {
		t.Errorf("Client not added to clients map")
	}

	// Check if the client was removed from the clients map
	err = conn.Close()
	if err != nil {
		t.Fatalf("Failed to close WebSocket connection: %v", err)
	}
	if _, ok := clients.m[conn]; ok {
		t.Errorf("Client not removed from clients map")
	}

	// Check if the client received the correct message
	var receivedMessage message
	if err := json.Unmarshal(msg, &receivedMessage); err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}
	if receivedMessage.Type != "status" {
		t.Errorf("Expected message type %s, got %s", "status", receivedMessage.Type)
	}
	if receivedMessage.Text != "You are connected. You will be matched with another user shortly." {
		t.Errorf("Expected message text %s, got %s", "You are connected. You will be matched with another user shortly.", receivedMessage.Text)
	}

}

// func TestCleanupQueue(t *testing.T) {
// 	// Set up a mock WebSocket server
// 	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		handleConnections(w, r, []string{"topic1", "topic2"})
// 	}))
// 	defer server.Close()

// 	// Set up a mock WebSocket client and add it to the waitingClients queue
// 	u := url.URL{Scheme: "ws", Host: strings.TrimPrefix(server.URL, "http://")}
// 	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
// 	if err != nil {
// 		t.Fatalf("Failed to connect to WebSocket server: %v", err)
// 	}
// 	waitingClients.Lock()
// 	waitingClients.q = append(waitingClients.q, conn)
// 	waitingClients.Unlock()

// 	// Disconnect the client and check if it was removed from the waitingClients queue
// 	err = conn.Close()
// 	if err != nil {
// 		t.Fatalf("Failed to close WebSocket connection: %v", err)
// 	}
// 	waitingClients.Lock()
// 	defer waitingClients.Unlock()
// 	if len(waitingClients.q) != 0 {
// 		t.Errorf("Client not removed from waitingClients queue")
// 	}
// }

// func TestRemoveClient(t *testing.T) {
// 	// Set up a mock WebSocket server
// 	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		handleConnections(w, r, []string{"topic1", "topic2"})
// 	}))
// 	defer server.Close()

// 	// Set up a mock WebSocket client and add it to the clients map
// 	u := url.URL{Scheme: "ws", Host: strings.TrimPrefix(server.URL, "http://")}
// 	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
// 	if err != nil {
// 		t.Fatalf("Failed to connect to WebSocket server: %v", err)
// 	}
// 	clients.Lock()
// 	clients.m[conn] = true
// 	clients.Unlock()

// 	// Remove the client and check if it was removed from the clients map
// 	removeClient(conn)
// 	clients.RLock()
// 	defer clients.RUnlock()
// 	if _, ok := clients.m[conn]; ok {
// 		t.Errorf("Client not removed from clients map")
// 	}
// }

// func TestRelayMessages(t *testing.T) {
// 	// Set up a mock WebSocket server
// 	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// do nothing
// 	}))
// 	defer server.Close()

// 	// Set up two mock WebSocket clients
// 	u1 := url.URL{Scheme: "ws", Host: strings.TrimPrefix(server.URL, "http://")}
// 	u2 := url.URL{Scheme: "ws", Host: strings.TrimPrefix(server.URL, "http://")}
// 	conn1, _, err := websocket.DefaultDialer.Dial(u1.String(), nil)
// 	if err != nil {
// 		t.Fatalf("Failed to connect to WebSocket server: %v", err)
// 	}
// 	conn2, _, err := websocket.DefaultDialer.Dial(u2.String(), nil)
// 	if err != nil {
// 		t.Fatalf("Failed to connect to WebSocket server: %v", err)
// 	}

// 	// Send a message from client 1 to client 2 and check if it is received
// 	msg := message{Type: "chat", Text: "Hello, world!"}
// 	jsonMsg, _ := json.Marshal(msg)
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		relayMessages(conn1, conn2, []string{})
// 	}()
// 	err = conn1.WriteMessage(websocket.TextMessage, jsonMsg)
// 	if err != nil {
// 		t.Fatalf("Failed to write message: %v", err)
// 	}
// 	_, actualMsg, err := conn2.ReadMessage()
// 	if err != nil {
// 		t.Fatalf("Failed to read message: %v", err)
// 	}
// 	var actualMsgObj message
// 	err = json.Unmarshal(actualMsg, &actualMsgObj)
// 	if err != nil {
// 		t.Fatalf("Failed to unmarshal message: %v", err)
// 	}
// 	if actualMsgObj.Text != msg.Text {
// 		t.Errorf("Received message does not match sent message: got %v, want %v", actualMsgObj.Text, msg.Text)
// 	}

// 	// Send a disconnect message from client 1 and check if both clients are removed
// 	disconnectMsg := message{Type: "disconnect", Text: ""}
// 	jsonDisconnectMsg, _ := json.Marshal(disconnectMsg)
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		relayMessages(conn1, conn2, []string{})
// 	}()
// 	err = conn1.WriteMessage(websocket.TextMessage, jsonDisconnectMsg)
// 	if err != nil {
// 		t.Fatalf("Failed to write message: %v", err)
// 	}
// 	wg.Wait()
// 	clients.RLock()
// 	defer clients.RUnlock()
// 	// Check if client 1 was removed from the clients map
// 	if _, ok := clients.m[conn1]; ok {
// 		t.Errorf("Client 1 not removed from clients map")
// 	}

// 	// Check if client 2 was removed from the clients map
// 	if _, ok := clients.m[conn2]; ok {
// 		t.Errorf("Client 2 not removed from clients map")
// 	}
// }

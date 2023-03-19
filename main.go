package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

var clients = make(map[*websocket.Conn]bool)
var clientsLock = sync.RWMutex{}
var waitingClients = []*websocket.Conn{}
var waitingClientsLock = sync.Mutex{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	//	defer conn.Close()

	clientsLock.Lock()
	clients[conn] = true
	clientsLock.Unlock()

	log.Printf("New user connected: %v", conn.RemoteAddr())

	matchmaking(conn)
}

func requeueClient(conn *websocket.Conn) {
	waitingClientsLock.Lock()
	defer waitingClientsLock.Unlock()

	waitingClients = append(waitingClients, conn)
	log.Printf("User %v added back to the waiting queue", conn.RemoteAddr())
	matchmaking(conn)
}

func matchmaking(conn *websocket.Conn) {
	waitingClientsLock.Lock()
	defer waitingClientsLock.Unlock()

	if len(waitingClients) > 0 {
		conn2 := waitingClients[0]
		waitingClients = waitingClients[1:]

		connectedMsg := message{Type: "status", Text: "You are now connected with another user."}
		jsonMsg, _ := json.Marshal(connectedMsg)

		conn.WriteMessage(websocket.TextMessage, jsonMsg)
		conn2.WriteMessage(websocket.TextMessage, jsonMsg)

		log.Printf("User %v connected with user %v", conn.RemoteAddr(), conn2.RemoteAddr())

		go chatHandler(conn, conn2)
	} else {
		waitingClients = append(waitingClients, conn)
		log.Printf("User %v added to the waiting queue", conn.RemoteAddr())
	}
}

func removeClient(conn *websocket.Conn) {
	clientsLock.Lock()
	delete(clients, conn)
	clientsLock.Unlock()

	waitingClientsLock.Lock()
	index := -1
	for i, waitingClient := range waitingClients {
		if waitingClient == conn {
			index = i
			break
		}
	}
	if index != -1 {
		waitingClients = append(waitingClients[:index], waitingClients[index+1:]...)
	}
	waitingClientsLock.Unlock()

	log.Printf("Client %v removed", conn.RemoteAddr())
}

func relayMessages(src *websocket.Conn, dest *websocket.Conn) {
	for {
		_, msg, err := src.ReadMessage()
		if err != nil {
			disconnectMsg := message{Type: "status", Text: "The other user has disconnected."}
			jsonMsg, _ := json.Marshal(disconnectMsg)
			dest.WriteMessage(websocket.TextMessage, jsonMsg)

			removeClient(src)
			removeClient(dest)
			log.Printf("User %v disconnected", src.RemoteAddr())
			break
		}

		err = dest.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("Write message error: %v", err)
			break
		}
	}
}

func chatHandler(conn1 *websocket.Conn, conn2 *websocket.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)

	cleanup := func() {
		wg.Done()
		wg.Wait()
		conn1.Close()
		conn2.Close()
	}

	go func() {
		defer cleanup()
		relayMessages(conn1, conn2)
	}()
	go func() {
		defer cleanup()
		relayMessages(conn2, conn1)
	}()
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

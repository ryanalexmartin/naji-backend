package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

var clients = struct {
	sync.RWMutex
	m map[*websocket.Conn]bool
}{m: make(map[*websocket.Conn]bool)}

var waitingClients = struct {
	sync.Mutex
	q []*websocket.Conn
}{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func handleQueuedClients(ctx context.Context) {
	// ping clients in the waiting queue every 5 seconds, to keep the connection alive
	// Send a message to every client that is currently connected,
	// and if the client is not responding, remove it from waitingClients.q
	pingTicker := time.NewTicker(5 * time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("handleQueuedClients stopped due to context cancellation")
			return
		case <-pingTicker.C:
			log.Printf("Waiting queue length: %v",
				len(waitingClients.q))
			waitingClients.Lock()
			activeClients := make([]*websocket.Conn, 0, len(waitingClients.q))

			for _, waitingClient := range waitingClients.q {
				pingMsg, _ := json.Marshal(map[string]string{"type": "ping"})
				err := waitingClient.WriteMessage(websocket.TextMessage, pingMsg)
				if err != nil {
					waitingClient.Close()
					log.Printf("Client %v removed from the waiting queue due to no response", waitingClient.RemoteAddr())
				} else {
					activeClients = append(activeClients, waitingClient)
				}
			}

			waitingClients.q = activeClients
			waitingClients.Unlock()
		}
	}
}

func main() {
	topics, err := loadTopics("topics.csv")
	if err != nil {
		log.Fatalf("Failed to load topics: %v", err)
	}
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(w, r, topics)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Println("Health check")
	})

	http.HandleFunc("/online-users", onlineUsersHandler)

	// run goroutine to handle queued clients, this goroutine will run forever
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handleQueuedClients(ctx)

	if _, err := os.Stat("/etc/letsencrypt/live/ws.naji.live/fullchain.pem"); os.IsNotExist(err) {
		log.Println("Certificate and key files not found, starting server on :8080")
		http.ListenAndServe(":8080", nil)
	} else {
		log.Println("Certificate and key files found, starting server on :443")
		http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/ws.naji.live/fullchain.pem", "/etc/letsencrypt/live/ws.naji.live/privkey.pem", nil)
	}

}

func onlineUsersHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	clients.RLock()
	onlineUsers := len(clients.m)
	clients.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"onlineUsers": onlineUsers})
}

func handleConnections(w http.ResponseWriter, r *http.Request, topics []string) {
	enableCors(&w)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	clients.Lock()
	clients.m[conn] = true
	clients.Unlock()

	log.Printf("New user connected: %v", conn.RemoteAddr())

	matchmaking(conn, topics)
}

func requeueClient(conn *websocket.Conn, topics []string) {
	waitingClients.Lock()
	defer waitingClients.Unlock()

	waitingClients.q = append(waitingClients.q, conn)
	log.Printf("User %v added back to the waiting queue", conn.RemoteAddr())
	matchmaking(conn, topics)
}

func matchmaking(conn *websocket.Conn, topics []string) {
	waitingClients.Lock()
	defer waitingClients.Unlock()

	if len(waitingClients.q) > 0 {

		conn2 := waitingClients.q[0]
		waitingClients.q = waitingClients.q[1:]

		randomTopic := topics[rand.Intn(len(topics))]
		connectedMsg := message{Type: "status", Text: fmt.Sprintf("Now connected! Let's talk about %s", randomTopic)}
		jsonMsg, _ := json.Marshal(connectedMsg)

		conn.WriteMessage(websocket.TextMessage, jsonMsg)
		conn2.WriteMessage(websocket.TextMessage, jsonMsg)

		log.Printf("User %v connected with user %v", conn.RemoteAddr(), conn2.RemoteAddr())

		go chatHandler(conn, conn2, topics)
	} else {
		waitingClients.q = append(waitingClients.q, conn)
		log.Printf("User %v added to the waiting queue", conn.RemoteAddr())
		log.Printf("Waiting queue length: %v", len(waitingClients.q))
	}
}

func removeClient(conn *websocket.Conn) {
	clients.Lock()
	delete(clients.m, conn)
	clients.Unlock()

	waitingClients.Lock()
	defer waitingClients.Unlock()
	for i, waitingClient := range waitingClients.q {
		if waitingClient == conn {
			waitingClients.q = append(waitingClients.q[:i], waitingClients.q[i+1:]...)
			break
		}
	}

	// log.Printf("Client %v removed", conn.RemoteAddr())
}

func relayMessages(src *websocket.Conn, dest *websocket.Conn, topics []string) {
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

		var receivedMessage message
		if err := json.Unmarshal(msg, &receivedMessage); err == nil {
			if receivedMessage.Type == "disconnect" {
				disconnectMsg := message{Type: "status", Text: "The other user has disconnected."}
				jsonMsg, _ := json.Marshal(disconnectMsg)
				dest.WriteMessage(websocket.TextMessage, jsonMsg)

				removeClient(src)
				removeClient(dest)
				log.Printf("User %v disconnected", src.RemoteAddr())
				break
			}
		}

		log.Printf("User %v sent message: %s", src.RemoteAddr(), string(msg))

		err = dest.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("Write message error: %v", err)
			break
		}
	}
}

func chatHandler(conn1 *websocket.Conn, conn2 *websocket.Conn, topics []string) {
	var wg sync.WaitGroup
	wg.Add(2)

	cleanup := func() {
		wg.Done()
		wg.Wait()
		removeClient(conn1)
		removeClient(conn2)
	}

	go func() {
		defer cleanup()
		relayMessages(conn1, conn2, topics)
	}()
	go func() {
		defer cleanup()
		relayMessages(conn2, conn1, topics)
	}()
}

func loadTopics(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var topics []string
	reader := csv.NewReader(bufio.NewReader(file))
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		topics = append(topics, record[0])
	}
	return topics, nil
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
